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
	chatHandler := handlers.NewRoomChatHandler(chatService)

	// === Business / Document / Chat (demo-api style) ===
	// Load business-related config (Pinecone, embed-service, MongoDB)
	bcfg := config.LoadConfig()

	// Separate MongoDB repository for business features (uses same URI / DB via env)
	mongoRepo, err := repo.NewMongoDBRepository(bcfg.MongoDBURI, bcfg.DatabaseName)
	if err != nil {
		log.Fatalf("Failed to connect MongoDB for business features: %v", err)
	}
	defer mongoRepo.Close()

	// External repositories
	embedRepo := repo.NewEmbedRepository(bcfg.EmbedServiceURL)
	pineconeRepo := repo.NewPineconeRepository(bcfg.PineconeIndexURL, bcfg.PineconeApiKey)
	businessRepo := repo.NewBusinessRepository(mongoRepo)
	businessUserRepo := repo.NewBusinessUserRepository(mongoRepo)

	// Services
	documentService := services.NewDocumentService(embedRepo, pineconeRepo)
	bizChatService := services.NewChatService(embedRepo, pineconeRepo)
	businessService := services.NewBusinessService(embedRepo, pineconeRepo, businessRepo, businessUserRepo)

	// Handlers
	uploadHandler := handlers.NewUploadHandler(documentService)
	bizChatHandler := handlers.NewChatHandler(bizChatService)
	businessHandler := handlers.NewBusinessHandler(businessService)

	// router
	fmt.Println("initializing router")
	r, err := router.NewRouter(
		resource.DB,
		chatHandler,
		roomHandler,
		authenticationHandler,
		uploadHandler,
		bizChatHandler,
		businessHandler,
	)
	if err != nil {
		log.Fatalf("Error initializing router", err)
	}

	// server
	fmt.Printf("Server listening on port %s\n", cfg.HTTP.Port)
	listenAddr := fmt.Sprintf(":%s", cfg.HTTP.Port)

	err = r.Serve(listenAddr)
	if err != nil {
		log.Fatalf("listen: %s\n", err)
	}
}
