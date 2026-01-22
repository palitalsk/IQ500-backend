package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"main/domain"
	"main/port"
	"main/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

type RoomChatHandler struct {
	svc        port.ChatService
	slipAPIURL string
}

func NewRoomChatHandler(svc port.ChatService, slipAPIURL string) *RoomChatHandler {
	return &RoomChatHandler{
		svc:        svc,
		slipAPIURL: slipAPIURL,
	}
}

func (h *RoomChatHandler) GetChatByRoomID(c *gin.Context) {
	fmt.Println("GetChatByRoomID")

	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	chat, err := h.svc.GetChatByRoomID(c, id)
	if err != nil {
		return
	}

	utils.Response(c, http.StatusOK, 200, "success", "ok", chat)
}

func (h *RoomChatHandler) Chat(c *gin.Context) {
	fmt.Println("Chat")

	var payload domain.PayloadChat
	if err := c.ShouldBind(&payload); err != nil {
		fmt.Println("error bind", err)
		utils.Response(c, http.StatusBadRequest, 1, "ข้อมูลไม่ถูกต้อง กรุณาลองใหม่อีกครั้ง", err.Error(), nil)
		return
	}

	if fileHeader, err := c.FormFile("image"); err == nil {
		file, err := fileHeader.Open()
		if err != nil {
			fmt.Println("error open image file", err)
			utils.Response(c, http.StatusBadRequest, 1, "เกิดข้อผิดพลาดในการอ่านไฟล์สลิป", err.Error(), nil)
			return
		}
		defer file.Close()

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			fmt.Println("error read image file", err)
			utils.Response(c, http.StatusBadRequest, 1, "เกิดข้อผิดพลาดในการอ่านไฟล์สลิป", err.Error(), nil)
			return
		}

		encoded := base64.StdEncoding.EncodeToString(fileBytes)

		reqBody := map[string]string{
			"image": encoded,
		}

		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			fmt.Println("error marshal slip request", err)
			utils.Response(c, http.StatusInternalServerError, 500, "เกิดข้อผิดพลาดภายในระบบ", err.Error(), nil)
			return
		}

		if h.slipAPIURL == "" {
			fmt.Println("error: SLIP_API_URL is not configured")
			utils.Response(c, http.StatusInternalServerError, 500, "บริการตรวจสอบสลิปยังไม่ได้ตั้งค่า", "SLIP_API_URL is not configured", nil)
			return
		}

		resp, err := http.Post(h.slipAPIURL, "application/json", bytes.NewReader(bodyBytes))
		if err != nil {
			fmt.Println("error call slip api", err)
			utils.Response(c, http.StatusInternalServerError, 500, "ไม่สามารถเชื่อมต่อบริการตรวจสอบสลิปได้", err.Error(), nil)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(resp.Body)
			fmt.Println("slip api returned non-200", resp.StatusCode, string(respBody))
			utils.Response(c, http.StatusInternalServerError, 500, "บริการตรวจสอบสลิปตอบกลับไม่สำเร็จ", string(respBody), nil)
			return
		}

		var slipResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&slipResp); err != nil {
			fmt.Println("error decode slip api response", err)
			utils.Response(c, http.StatusInternalServerError, 500, "ไม่สามารถอ่านผลลัพธ์จากบริการตรวจสอบสลิปได้", err.Error(), nil)
			return
		}

		var payloadData map[string]interface{}

		if p, ok := slipResp["payload"].(map[string]interface{}); ok {
			payloadData = p
			fmt.Println("Found data in payload key")
		} else {
			payloadData = slipResp
			fmt.Println("Using slipResp directly")
		}

		aiResult := make(map[string]interface{})

		if isSlip, ok := payloadData["is_slip"].(bool); ok {
			aiResult["is_slip"] = isSlip
			fmt.Printf("Extracted is_slip: %v\n", isSlip)
		}

		if ocrResult, ok := payloadData["ocr_result"].(map[string]interface{}); ok {
			aiResult["ocr_result"] = ocrResult
			fmt.Printf("Extracted ocr_result with %d fields\n", len(ocrResult))
		} else {
			fmt.Println("Warning: ocr_result is not a map[string]interface{}")
		}

		if len(aiResult) > 0 {
			payload.AiResult = aiResult
			fmt.Printf("AiResult prepared for MongoDB: is_slip=%v, ocr_result_fields=%d\n",
				aiResult["is_slip"],
				func() int {
					if ocr, ok := aiResult["ocr_result"].(map[string]interface{}); ok {
						return len(ocr)
					}
					return 0
				}())
		} else {
			fmt.Println("Warning: aiResult is empty, not setting in payload")
		}

		payload.Img = fileHeader.Filename
	} else {
		fmt.Println("no image file in request:", err)
	}

	chat, err := h.svc.Chat(c, payload)
	if err != nil {
		return
	}

	utils.Response(c, http.StatusOK, 200, "success", "ok", chat)
}
