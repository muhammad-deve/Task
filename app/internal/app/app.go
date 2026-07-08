package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gitlab.yurtal.tech/company/blitz/back/internal/api/docs"
	"gitlab.yurtal.tech/company/blitz/back/internal/config"
	"gitlab.yurtal.tech/company/blitz/back/internal/handler"
	"gitlab.yurtal.tech/company/blitz/back/internal/middleware"
	"gitlab.yurtal.tech/company/blitz/back/internal/migrate"
	"gitlab.yurtal.tech/company/blitz/back/internal/repository"
	"gitlab.yurtal.tech/company/blitz/back/internal/service"
	"gitlab.yurtal.tech/company/blitz/back/pkg/logger"
	"gitlab.yurtal.tech/company/blitz/back/pkg/minio"
	pg "gitlab.yurtal.tech/company/blitz/back/pkg/postgres"
)

// @title Driver Registry Service API
// @version 1.0
// @description Taxi driver management service with CRUD operations and status statistics
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @BasePath /
// @schemes https http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @tag.name auth
// @tag.description Authentication endpoints

// @tag.name drivers
// @tag.description Driver management operations

// @tag.name system
// @tag.description System health and monitoring
func Run(cfg *config.Config) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	l := logger.New(cfg.Logger.Level)
	e := echo.New()

	middleware.SetupMiddleware(e, cfg)

	pgClient, err := pg.New(pg.Username(cfg.Postgres.User), pg.Password(cfg.Postgres.Password),
		pg.Host(cfg.Postgres.Host), pg.Port(cfg.Postgres.Port),
		pg.Database(cfg.Postgres.Db), pg.MaxPoolSize(cfg.Postgres.MaxPoolSize))
	if err != nil {
		l.Fatalf("app - Run - pg.New: %v", err)
	}
	defer pgClient.Close()

	err = migrate.RunMigrations(ctx, pgClient.Pool)
	if err != nil {
		l.Fatalf("app - Run - RunMigrations: %v", err)
	}

	minioClient, err := minio.New(minio.Endpoint(cfg.Minio.Endpoint), minio.AccessKeyID(cfg.Minio.AccessKey), minio.SecretAccessKey(cfg.Minio.SecretKey))
	if err != nil {
		l.Fatalf("app - Run - minio.New: %v", err)
	}

	repos := repository.New(pgClient, minioClient)

	service := service.New(cfg, repos)

	handler := handler.New(l, cfg, service, pgClient.Pool)
	handler.Register(e)

	// Dynamic swagger spec — patches host+scheme from each incoming request so
	// the same binary works on localhost:8080 (http) and task.goport.uz (https)
	// without any manual changes or re-building.
	e.GET("/swagger/doc.json", func(c echo.Context) error {
		var spec map[string]interface{}
		if err := json.Unmarshal([]byte(docs.SwaggerInfo.ReadDoc()), &spec); err != nil {
			return err
		}

		// By removing host and schemes completely, Swagger UI will automatically
		// use the exact host and protocol (http/https) from the browser's URL bar.
		// This makes it work flawlessly on localhost, goport, ngrok, or any domain.
		delete(spec, "host")
		delete(spec, "schemes")

		return c.JSON(http.StatusOK, spec)
	})

	e.GET("/swagger/*", echoSwagger.EchoWrapHandler(
		echoSwagger.URL("/swagger/doc.json"),
		echoSwagger.DocExpansion("list"),
	))

	errc := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		l.Info("starting server on %s", fmt.Sprintf(":%d", cfg.Server.Http.Port))
		if err := e.Start(fmt.Sprintf(":%d", cfg.Server.Http.Port)); err != nil && err != http.ErrServerClosed {
			l.Fatalf("app - Run - e.Start: %v", err)
		}
	}()

	err = <-errc
	l.Infof("shutdown initiated: %v", err)

	cancel()

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := e.Shutdown(ctxShutdown); err != nil {
		l.Errorf("server shutdown: %v", err)
	}

	l.Info("exited cleanly")
}
