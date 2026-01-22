package http

import (
	"fmt"
	"log"

	"main/config"
	"main/handlers"
	"main/internal/adapter/storage/mongo"
	"main/repo"
	"main/router"
	"main/services"
)

func HttpMain(
	cfg *config.Container,
	resource *mongo.Resource,
) {
	// === Dependency Injection (Hexagonal part) ===
	// Repositories
	userRepo := repo.NewUserRepository(resource.DB)
	roomRepo := repo.NewRoomRepository(resource.DB)
	chatRepo := repo.NewChatRepository(resource.DB)

	// Services (Application)
	userService := services.NewAuthenticationService(userRepo)
	roomService := services.NewRoomService(roomRepo)
	chatService := services.NewRoomChatService(chatRepo, roomRepo)

	// HTTP Handlers (Adapters)
	authenticationHandler := handlers.NewAuthHandler(userService)
	roomHandler := handlers.NewRoomHandler(roomService)
	chatHandler := handlers.NewRoomChatHandler(chatService, cfg.HTTP.SlipAPIURL)

	bcfg := config.LoadConfig()

	embedRepo := repo.NewEmbedRepository(bcfg.EmbedServiceURL, bcfg.GeminiApiKey)
	pineconeRepo := repo.NewPineconeRepository(bcfg.PineconeIndexURL, bcfg.PineconeApiKey)

	documentService := services.NewDocumentService(embedRepo, pineconeRepo)
	bizChatService := services.NewChatService(embedRepo, pineconeRepo)

	uploadHandler := handlers.NewUploadHandler(documentService)
	bizChatHandler := handlers.NewChatHandler(bizChatService)

	// router
	fmt.Println("initializing router")
	r, err := router.NewRouter(
		resource.DB,
		chatHandler,
		roomHandler,
		authenticationHandler,
		uploadHandler,
		bizChatHandler,
	)
	if err != nil {
		log.Fatalf("Error initializing router: %v", err)
	}

	// server
	fmt.Printf("Server listening on port %s\n", cfg.HTTP.Port)
	listenAddr := fmt.Sprintf(":%s", cfg.HTTP.Port)

	err = r.Serve(listenAddr)
	if err != nil {
		log.Fatalf("listen: %s\n", err)
	}
}
