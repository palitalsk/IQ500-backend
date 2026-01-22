package router

import (
	"main/handlers"
	"main/utils"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
)

type Router struct {
	*gin.Engine
}

func NewRouter(
	resource *mongo.Database,
	chatHandler *handlers.RoomChatHandler,
	roomHandler *handlers.RoomHandler,
	authHandler *handlers.AuthHandler,
	uploadHandler *handlers.UploadHandler,
	bizChatHandler *handlers.ChatHandler,
) (*Router, error) {
	r := gin.Default()
	r.Use(CORS)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"version": "1.0.1",
		})
	})

	api := r.Group("/api")
	auth := r.Group("/api")
	auth.Use(utils.Authentication(resource))

	api.POST("/login", authHandler.Login)
	api.POST("/register", authHandler.Register)
	api.POST("/reset-password", authHandler.ResetPassword)

	chat := auth.Group("/chat")
	{
		chat.GET("/:id", chatHandler.GetChatByRoomID)
		chat.POST("", chatHandler.Chat)
	}

	room := auth.Group("/room")
	{
		room.GET("", roomHandler.GetRoom)
		room.POST("", roomHandler.CreateRoom)
		room.DELETE("/:id", roomHandler.DeleteRoomByID)
	}

	chatbotAPI := r.Group("/api")
	chatbotAPI.Use(utils.Authentication(resource))

	chatbotAPI.POST("/upload", uploadHandler.HandleUpload)
	chatbotAPI.POST("/chatbot", bizChatHandler.HandleChat)

	return &Router{r}, nil
}

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

func (r *Router) Serve(listenAddr string) error {
	srv := &http.Server{
		Addr:    listenAddr,
		Handler: r,
	}
	return srv.ListenAndServe()
}
