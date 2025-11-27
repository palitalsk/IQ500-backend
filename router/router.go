package router

import (
	"main/handlers"
	"main/middleware"
	"main/utils"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
)

// Router เป็นตัวห่อ gin.Engine ให้สามารถแนบ Serve() ใช้งานจาก main ได้ง่าย
type Router struct {
	*gin.Engine
}

// NewRouter ประกอบทุก handler และ middleware เข้าด้วยกัน
func NewRouter(
	resource *mongo.Database,
	chatHandler *handlers.RoomChatHandler,
	roomHandler *handlers.RoomHandler,
	authHandler *handlers.AuthHandler,
	uploadHandler *handlers.UploadHandler,
	bizChatHandler *handlers.ChatHandler,
	businessHandler *handlers.BusinessHandler,
) (*Router, error) {
	r := gin.Default()
	r.Use(CORS)

	// ping
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"version": "1.0.1",
		})
	})

	api := r.Group("/api")
	auth := r.Group("/api")
	auth.Use(utils.Authentication(resource))

	// authentication (IQ500)
	api.POST("/login", authHandler.Login)
	api.POST("/register", authHandler.Register)
	api.POST("/reset-password", authHandler.ResetPassword)

	// room chat (IQ500)
	chat := auth.Group("/chat")
	{
		chat.GET("/:id", chatHandler.GetChatByRoomID)
		chat.POST("", chatHandler.Chat)
	}

	// room (IQ500)
	room := auth.Group("/room")
	{
		room.GET("", roomHandler.GetRoom)
		room.POST("", roomHandler.CreateRoom)
		room.DELETE("/:id", roomHandler.DeleteRoomByID)
	}

	// === Business & Document & Chat (demo-api style) ===
	// These routes use header-based auth (X-User-ID) and business_code in path/query
	bizAPI := r.Group("/api")
	bizAPI.Use(middleware.AuthMiddleware())

	// Business routes
	bizAPI.POST("/businesses", businessHandler.CreateBusiness)
	bizAPI.GET("/businesses", businessHandler.GetUserBusinesses)
	bizAPI.GET("/businesses/:business_code", businessHandler.GetBusiness)

	// Business-specific routes
	business := bizAPI.Group("/businesses/:business_code")
	business.Use(middleware.BusinessAccessMiddleware())
	{
		business.POST("/upload", uploadHandler.HandleUpload)
		business.POST("/chat", bizChatHandler.HandleChat)
	}

	return &Router{r}, nil
}

// CORS ตั้งค่า header CORS พื้นฐานให้ทุก request
func CORS(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Content-Type", "application/json")

	if c.Request.Method != "OPTIONS" {
		c.Next()
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

// Serve เริ่มต้น HTTP server ด้วย address ที่กำหนด
func (r *Router) Serve(listenAddr string) error {
	srv := &http.Server{
		Addr:    listenAddr,
		Handler: r,
	}
	return srv.ListenAndServe()
}


