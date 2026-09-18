package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"genealogy-tree/internal/core/config"
	corelogger "genealogy-tree/internal/core/logger"
	"genealogy-tree/internal/core/repository/postgres"
	"genealogy-tree/internal/core/security"
	"genealogy-tree/internal/core/transport/http/middleware"
	httpserver "genealogy-tree/internal/core/transport/http/server"
	authservice "genealogy-tree/internal/features/auth/service"
	authhttp "genealogy-tree/internal/features/auth/transport/http"
	personrepo "genealogy-tree/internal/features/persons/repository"
	personminio "genealogy-tree/internal/features/persons/repository/minio"
	personservice "genealogy-tree/internal/features/persons/service"
	personhttp "genealogy-tree/internal/features/persons/transport/http"
	relrepo "genealogy-tree/internal/features/relationships/repository"
	relservice "genealogy-tree/internal/features/relationships/service"
	relhttp "genealogy-tree/internal/features/relationships/transport/http"
	treerepo "genealogy-tree/internal/features/trees/repository"
	treeservice "genealogy-tree/internal/features/trees/service"
	treehttp "genealogy-tree/internal/features/trees/transport/http"
	usersrepo "genealogy-tree/internal/features/users/repository"
	usersservice "genealogy-tree/internal/features/users/service"
	usershttp "genealogy-tree/internal/features/users/transport/http"

	"go.uber.org/zap"

	_ "genealogy-tree/docs"
)

// @title 		Genealogy Tree API
// @version 	1.0
// @description Genealogy Tree REST API
// @host 		127.0.0.1:8080
// @BasePath 	/api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	logger, err := corelogger.NewLogger(corelogger.NewConfigMust())
	if err != nil {
		log.Fatal("failed to init application logger:", err)
	}
	defer logger.Close()

	logger.Debug("initializing postgres connection pool")

	pool, err := postgres.NewPool(context.Background(), cfg.DatabaseURL())
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	tokenManager, err := security.NewTokenManager(cfg.JWTSecret, cfg.JWTIssuer, cfg.JWTAccessTTL)
	if err != nil {
		logger.Fatal("failed to init token manager", zap.Error(err))
	}

	usersRepository := usersrepo.NewUsersRepository(pool)
	usersService := usersservice.NewUsersService(usersRepository)
	authService := authservice.NewAuthService(usersRepository, tokenManager)
	authTransportHTTP := authhttp.NewAuthHTTPHandler(authService)
	usersTransportHTTP := usershttp.NewUsersHTTPHandler(usersService)
	logger.Debug("initializing minio S3-storage")

	fileStorage, err := personminio.NewMinIOFileStorage(
		cfg.MinIOEndpoint,
		cfg.MinIOAccessKey,
		cfg.MinIOSecretKey,
		cfg.MinIOBucket,
		cfg.MinIOUseSSL,
	)
	if err != nil {
		logger.Fatal("failed to init MinIO file storage", zap.Error(err))
	}

	treesRepository := treerepo.NewTreesRepository(pool)
	treesService := treeservice.NewTreesService(treesRepository, fileStorage)
	treesTransportHTTP := treehttp.NewTreesHTTPHandler(treesService)

	logger.Debug("initializing feature", zap.String("feature", "persons"))

	personsRepository := personrepo.NewPersonsRepository(pool)
	personsService := personservice.NewPersonsService(personsRepository, fileStorage, treesService)
	personsTransportHTTP := personhttp.NewPersonsHTTPHandlers(personsService)

	logger.Debug("initializing feature", zap.String("feature", "relationships"))

	relationshipsRepository := relrepo.NewRelationshipsRepository(pool)
	relationshipsService := relservice.NewRelationshipsService(relationshipsRepository, treesService)
	relationshipsTransportHTTP := relhttp.NewRelationshipsHTTPHandlers(relationshipsService)

	httpConfig := httpserver.NewConfigMust()
	server := httpserver.NewHTTPServer(
		httpConfig,
		logger,
		middleware.CORS(httpConfig.AllowedOrigins),
		middleware.RequestID(),
		middleware.Logger(logger),
		middleware.Trace(),
		middleware.Panic(),
	)

	apiVersionRouterV1 := httpserver.NewApiVersionRouter(httpserver.ApiVersion1)
	apiVersionRouterV1.RegisterRoutes(authTransportHTTP.Routes()...)
	authMiddleware := middleware.Auth(tokenManager)
	protectedRoutes := [][]httpserver.Route{
		usersTransportHTTP.Routes(),
		treesTransportHTTP.Routes(),
		personsTransportHTTP.Routes(),
		relationshipsTransportHTTP.Routes(),
	}
	for _, routes := range protectedRoutes {
		for _, route := range routes {
			route.Middleware = append(route.Middleware, authMiddleware)
			apiVersionRouterV1.RegisterRoutes(route)
		}
	}
	server.RegisterAPIRouters(apiVersionRouterV1)

	server.RegisterRoutes(httpserver.Route{Method: http.MethodGet, Path: "/health", Handler: func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}})
	server.RegisterSwagger()

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	logger.Debug("HTTP server configured", zap.String("addr", cfg.HTTPAddr))

	if err := server.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
