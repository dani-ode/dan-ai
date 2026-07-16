// apps/api/bootstrap/app.go
package bootstrap

import (
	"fmt"
	"net"

	aimodelgrpc "dan-ai/internal/aimodel/grpc"
	aimodelrepo "dan-ai/internal/aimodel/repository"
	aimodelsvc "dan-ai/internal/aimodel/service"
	authgrpc "dan-ai/internal/auth/grpc"
	"dan-ai/internal/auth/jwt"
	certgrpc "dan-ai/internal/certificate/grpc"
	certrepo "dan-ai/internal/certificate/repository"
	certsvc "dan-ai/internal/certificate/service"
	chatgrpc "dan-ai/internal/chat/grpc"
	chatrepo "dan-ai/internal/chat/repository"
	chatsvc "dan-ai/internal/chat/service"
	experiencegrpc "dan-ai/internal/experience/grpc"
	experiencerepo "dan-ai/internal/experience/repository"
	experiencesvc "dan-ai/internal/experience/service"
	knowledgegrpc "dan-ai/internal/knowledge/grpc"
	knowledgerepo "dan-ai/internal/knowledge/repository"
	knowledgesvc "dan-ai/internal/knowledge/service"
	outboxrepo "dan-ai/internal/outbox/repository"
	profilegrpc "dan-ai/internal/profile/grpc"
	profilerepo "dan-ai/internal/profile/repository"
	profilesvc "dan-ai/internal/profile/service"
	projectgrpc "dan-ai/internal/project/grpc"
	projectrepo "dan-ai/internal/project/repository"
	projectsvc "dan-ai/internal/project/service"
	promptgrpc "dan-ai/internal/prompt/grpc"
	promptrepo "dan-ai/internal/prompt/repository"
	promptsvc "dan-ai/internal/prompt/service"
	skillgrpc "dan-ai/internal/skill/grpc"
	skillrepo "dan-ai/internal/skill/repository"
	skillsvc "dan-ai/internal/skill/service"
	techgrpc "dan-ai/internal/technology/grpc"
	techrepo "dan-ai/internal/technology/repository"
	techsvc "dan-ai/internal/technology/service"
	visitorgrpc "dan-ai/internal/visitor/grpc"
	visitorrepo "dan-ai/internal/visitor/repository"
	visitorsvc "dan-ai/internal/visitor/service"
	"dan-ai/pkg/config"
	"dan-ai/pkg/logger"
	aimodelpb "dan-ai/proto/aimodel"
	pb "dan-ai/proto/auth"
	certpb "dan-ai/proto/certificate"
	chatpb "dan-ai/proto/chat"
	experiencepb "dan-ai/proto/experience"
	knowledgepb "dan-ai/proto/knowledge"
	profilepb "dan-ai/proto/profile"
	projectpb "dan-ai/proto/project"
	promptpb "dan-ai/proto/prompt"
	skillpb "dan-ai/proto/skill"
	techpb "dan-ai/proto/technology"
	visitorpb "dan-ai/proto/visitor"

	"google.golang.org/grpc"
	gormdb "gorm.io/gorm"
)

type App struct {
	Config     *config.Config
	DB         *gormdb.DB
	GRPCServer *grpc.Server
}

// NewApp initializes configuration, logger, database, and gRPC server.
func NewApp() (*App, error) {
	// 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// 2. Initialize structured logger
	logger.Init(cfg.App.Env)
	logger.Info("Starting application", "app_name", cfg.App.Name, "env", cfg.App.Env)

	// 3. Connect to database
	db, err := NewDatabase(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}
	logger.Info("Database connection established", "host", cfg.DB.Host, "database", cfg.DB.Name)

	// 4. Initialize JWT manager
	jwtManager := jwt.NewManager(cfg.Auth.JWTSecret, cfg.Auth.JWTExpiry)
	logger.Info("JWT manager initialized", "expiry", cfg.Auth.JWTExpiry.String())

	// 5. Create gRPC server (with auth interceptor)
	grpcServer := NewGRPCServer(jwtManager)
	logger.Info("gRPC server initialized")

	// 6. Register Auth service (login without DB)
	authHandler := authgrpc.NewHandler(jwtManager, cfg.Auth.AdminUsername, cfg.Auth.AdminPassword, cfg.Auth.JWTExpiry)
	pb.RegisterAuthServiceServer(grpcServer, authHandler)
	logger.Info("Auth service registered")

	// 6.5 Initialize Knowledge module (core of Phase 3)
	knowledgeRepo := knowledgerepo.NewPostgresKnowledgeRepository(db)
	knowledgeService := knowledgesvc.NewService(knowledgeRepo)
	knowledgeHandler := knowledgegrpc.NewKnowledgeHandler(knowledgeService)
	knowledgepb.RegisterKnowledgeServiceServer(grpcServer, knowledgeHandler)
	logger.Info("Knowledge service registered")

	// 7. Initialize Profile module
	profileRepo := profilerepo.NewPostgresRepository(db)
	profileService := profilesvc.NewService(profileRepo, knowledgeService)
	profileHandler := profilegrpc.NewHandler(profileService)
	profilepb.RegisterProfileServiceServer(grpcServer, profileHandler)
	logger.Info("Profile service registered")

	// 8. Initialize Prompt module
	promptRepo := promptrepo.NewPostgresRepository(db)
	promptService := promptsvc.NewService(promptRepo)
	promptHandler := promptgrpc.NewHandler(promptService)
	promptpb.RegisterPromptServiceServer(grpcServer, promptHandler)
	logger.Info("Prompt service registered")

	// 9. Initialize Visitor module
	visitorRepo := visitorrepo.NewPostgresRepository(db)
	visitorService := visitorsvc.NewService(visitorRepo)
	visitorHandler := visitorgrpc.NewHandler(visitorService)
	visitorpb.RegisterVisitorServiceServer(grpcServer, visitorHandler)
	logger.Info("Visitor service registered")

	// 10. Initialize Outbox module for chat events
	outboxRepo := outboxrepo.NewPostgresRepository(db)

	// 11. Initialize Chat module
	chatRepo := chatrepo.NewPostgresRepository(db)
	chatService := chatsvc.NewService(chatRepo, outboxRepo)
	chatHandler := chatgrpc.NewHandler(chatService, visitorService)
	chatpb.RegisterChatServiceServer(grpcServer, chatHandler)
	logger.Info("Chat service registered")

	// 11. Initialize AIModel module
	aimodelRepo := aimodelrepo.NewPostgresRepository(db)
	aimodelService := aimodelsvc.NewService(aimodelRepo)
	aimodelHandler := aimodelgrpc.NewHandler(aimodelService)
	aimodelpb.RegisterAIModelServiceServer(grpcServer, aimodelHandler)
	logger.Info("AIModel service registered")

	// 12. Initialize Project module
	projectRepo := projectrepo.NewPostgresRepository(db)
	projectService := projectsvc.NewService(projectRepo, knowledgeService)
	projectHandler := projectgrpc.NewHandler(projectService)
	projectpb.RegisterProjectServiceServer(grpcServer, projectHandler)
	logger.Info("Project service registered")

	// 13. Initialize Experience module
	experienceRepo := experiencerepo.NewPostgresRepository(db)
	experienceService := experiencesvc.NewService(experienceRepo, knowledgeService)
	experienceHandler := experiencegrpc.NewHandler(experienceService)
	experiencepb.RegisterExperienceServiceServer(grpcServer, experienceHandler)
	logger.Info("Experience service registered")

	// 14. Initialize Technology module
	techRepo := techrepo.NewPostgresRepository(db)
	techService := techsvc.NewService(techRepo)
	techHandler := techgrpc.NewHandler(techService)
	techpb.RegisterTechnologyServiceServer(grpcServer, techHandler)
	logger.Info("Technology service registered")

	// 15. Initialize Certificate module
	certRepo := certrepo.NewPostgresRepository(db)
	certService := certsvc.NewService(certRepo, knowledgeService)
	certHandler := certgrpc.NewHandler(certService)
	certpb.RegisterCertificateServiceServer(grpcServer, certHandler)
	logger.Info("Certificate service registered")

	// 16. Initialize Skill module
	skillRepo := skillrepo.NewPostgresRepository(db)
	skillService := skillsvc.NewService(skillRepo, knowledgeService)
	skillHandler := skillgrpc.NewHandler(skillService)
	skillpb.RegisterSkillServiceServer(grpcServer, skillHandler)
	logger.Info("Skill service registered")

	return &App{
		Config:     cfg,
		DB:         db,
		GRPCServer: grpcServer,
	}, nil
}

// Run starts the gRPC server on the configured port.
func (a *App) Run() error {
	addr := ":" + a.Config.App.Port
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	logger.Info("gRPC server listening", "address", addr)
	return a.GRPCServer.Serve(lis)
}
