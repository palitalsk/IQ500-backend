package http

import (
	"fmt"
	"log"
	"main/internal/adapter/config"
	httpServer "main/internal/adapter/handler/http"
	"main/internal/adapter/storage/mongo"
	"main/internal/adapter/storage/mongo/repository"
	"main/internal/core/service"
)

func HttpMain(
	config *config.Container,
	resource *mongo.Resource,
) {
	// === Dependency Injection ===
	// Repositories
	userRepo := repository.NewUserRepository(resource.DB)
	roomRepo := repository.NewRoomRepository(resource.DB)
	chatRepo := repository.NewChatRepository(resource.DB)

	// AI Service
	aiService := service.NewAIService(config.AI.BaseURL, config.AI.Timeout)

	// Services (Application)
	userService := service.NewAuthenticationService(userRepo)
	roomService := service.NewRoomService(roomRepo)
	chatService := service.NewChatService(chatRepo, roomRepo, aiService)

	// HTTP Handlers (Adapters)
	authenticationHandler := httpServer.NewAuthenticationHandler(userService)
	roomHandler := httpServer.NewRoomHandler(roomService)
	chatHandler := httpServer.NewChatHandler(chatService)

	// router
	fmt.Println("initializing router")
	router, err := httpServer.NewRouter(
		resource.DB,
		chatHandler,
		roomHandler,
		authenticationHandler,
	)
	if err != nil {
		log.Fatalf("Error initializing router", err)
	}

	// server
	fmt.Printf("Server listening on port %s\n", config.HTTP.Port)
	listenAddr := fmt.Sprintf(":%s", config.HTTP.Port)

	err = router.Serve(listenAddr)
	if err != nil {
		log.Fatalf("listen: %s\n", err)
	}
}
