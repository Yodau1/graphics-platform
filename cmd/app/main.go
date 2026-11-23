package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	core_logger "github.com/Yodau1/graphics-platform/internal/core/logger"
	core_pgx_pool "github.com/Yodau1/graphics-platform/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/Yodau1/graphics-platform/internal/core/transport/http/middleware"
	core_http_server "github.com/Yodau1/graphics-platform/internal/core/transport/http/server"
	categories_postgres_repository "github.com/Yodau1/graphics-platform/internal/features/categories/repository/postgres"
	categories_service "github.com/Yodau1/graphics-platform/internal/features/categories/service"
	categories_transport_http "github.com/Yodau1/graphics-platform/internal/features/categories/transport/http"
	resources_filesystem_repository "github.com/Yodau1/graphics-platform/internal/features/resources/repository/filesystem"
	resources_postgres_repository "github.com/Yodau1/graphics-platform/internal/features/resources/repository/postgres"
	resources_service "github.com/Yodau1/graphics-platform/internal/features/resources/service"
	resources_transport_http "github.com/Yodau1/graphics-platform/internal/features/resources/transport/http"
	users_postgres_repository "github.com/Yodau1/graphics-platform/internal/features/users/repository/postgres"
	users_service "github.com/Yodau1/graphics-platform/internal/features/users/service"
	users_transport_http "github.com/Yodau1/graphics-platform/internal/features/users/transport/http"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to init application logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	postgresPool, err := core_pgx_pool.NewPool(
		ctx,
		core_pgx_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to init PostgreSQL pool", zap.Error(err))
	}
	defer postgresPool.Close()

	//Получаем путь, куда сохранять файлы
	fileStorageConfig, err := resources_filesystem_repository.NewConfig()
	if err != nil {
		logger.Fatal("failed to load resource file storage config", zap.Error(err))
	}

	resourceFileStorage := resources_filesystem_repository.NewResourceFileStorage(
		fileStorageConfig.Folder,
	)

	usersRepository := users_postgres_repository.NewUsersRepository(postgresPool)
	usersService := users_service.NewUsersService(usersRepository)
	requireAuthenticated := core_http_middleware.RequireAuthenticatedUser(usersService)
	usersHTTPHandler := users_transport_http.NewUsersHTTPHandler(usersService, requireAuthenticated)

	resourcesRepository := resources_postgres_repository.NewResourcesRepository(postgresPool)
	resourcesService := resources_service.NewResourcesService(resourcesRepository, resourceFileStorage)
	resourcesHTTPHandler := resources_transport_http.NewResourcesHTTPHandler(resourcesService, requireAuthenticated)

	categoriesRepository := categories_postgres_repository.NewCategoriesRepository(postgresPool)
	categoriesService := categories_service.NewCategoriesService(categoriesRepository)
	categoriesHTTPHandler := categories_transport_http.NewCategoriesHTTPHandler(categoriesService, requireAuthenticated)

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	apiV1 := core_http_server.NewAPIVersionRouter(core_http_server.APIVersion1)
	apiV1.AddRoutes(usersHTTPHandler.Routes()...)
	apiV1.AddRoutes(categoriesHTTPHandler.Routes()...)
	apiV1.AddRoutes(resourcesHTTPHandler.Routes()...)
	httpServer.RegisterAPIRouters(apiV1)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("run HTTP server", zap.Error(err))
		panic(err)
	}
}
