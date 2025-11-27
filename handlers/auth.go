package handlers

import (
	"fmt"
	"main/domain"
	"main/port"
	"main/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthHandler จัดการ endpoint ด้าน authentication (login / register / reset-password)
type AuthHandler struct {
	svc port.AuthenticationService
}

func NewAuthHandler(svc port.AuthenticationService) *AuthHandler {
	return &AuthHandler{
		svc: svc,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	fmt.Println("Login")

	var payload domain.PayloadUser

	if err := c.ShouldBind(&payload); err != nil {
		fmt.Println("error bind", err)
		utils.Response(c, http.StatusBadRequest, 400, "ข้อมูลไม่ถูกต้อง กรุณาลองใหม่อีกครั้ง", err.Error(), nil)
		return
	}

	token, err := h.svc.Login(c, payload)
	if err != nil {
		return
	}

	utils.Response(c, http.StatusOK, 200, "success", "ok", token)
}

func (h *AuthHandler) Register(c *gin.Context) {
	fmt.Println("Register")

	var payload domain.PayloadUser
	if err := c.ShouldBind(&payload); err != nil {
		fmt.Println("error bind", err)
		utils.Response(c, http.StatusBadRequest, 400, "ข้อมูลไม่ถูกต้อง กรุณาลองใหม่อีกครั้ง", err.Error(), nil)
		return
	}

	if err := h.svc.Register(c, payload); err != nil {
		return
	}

	utils.Response(c, http.StatusOK, 200, "success", "ok", nil)
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	fmt.Println("ResetPassword")

	var payload domain.PayloadResetPassword
	if err := c.ShouldBind(&payload); err != nil {
		fmt.Println("error bind", err)
		utils.Response(c, http.StatusBadRequest, 400, "ข้อมูลไม่ถูกต้อง กรุณาลองใหม่อีกครั้ง", err.Error(), nil)
		return
	}

	if err := h.svc.ResetPassword(c, payload); err != nil {
		return
	}

	utils.Response(c, http.StatusOK, 200, "success", "ok", nil)
}


