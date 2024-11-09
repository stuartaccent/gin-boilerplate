package cmd

import (
	"context"
	"errors"
	"fmt"
	"gin.go.dev/pkg/auth"
	"gin.go.dev/pkg/home"
	"gin.go.dev/pkg/static"
	"gin.go.dev/pkg/transport/html"
	"gin.go.dev/pkg/transport/middleware"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/secure"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	slogformatter "github.com/samber/slog-formatter"
	sloggin "github.com/samber/slog-gin"
	"github.com/spf13/cobra"
	csrf "github.com/stuartaccent/gin-csrf"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

var cmdServer = &cobra.Command{
	Use:   "server",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		runServer()
	},
}

func initLogging(e *gin.Engine) {
	logger := slog.New(
		slogformatter.NewFormatterHandler(
			slogformatter.TimezoneConverter(time.UTC),
			slogformatter.TimeFormatter(time.RFC3339, nil),
		)(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}),
		),
	)
	logger = logger.With("gin_mode", cfg.Server.Mode.ToGinMode())

	config := sloggin.Config{
		WithUserAgent: true,
		WithRequestID: true,
	}

	e.Use(sloggin.NewWithConfig(logger, config))
}

func initPool() *pgxpool.Pool {
	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, cfg.Database.URL().String())
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}

	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}

	return dbPool
}

func runServer() {
	gin.SetMode(cfg.Server.Mode.ToGinMode())

	dbPool := initPool()
	defer dbPool.Close()

	csrfMiddleware := csrf.Middleware(csrf.Options{
		Secret: cfg.Security.CsrfSecret,
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	})

	secureMiddleware := secure.New(secure.Config{
		AllowedHosts:          cfg.Security.AllowedHosts,
		STSSeconds:            cfg.Security.StsSeconds,
		STSIncludeSubdomains:  cfg.Security.StsIncludeSubdomains,
		FrameDeny:             cfg.Security.FrameDeny,
		ContentTypeNosniff:    cfg.Security.ContentTypeNosniff,
		BrowserXssFilter:      cfg.Security.BrowserXSSFilter,
		ContentSecurityPolicy: cfg.Security.ContentSecurityPolicy,
	})

	sessionStore := cookie.NewStore(cfg.Session.KeyBytes(), cfg.Session.EncKeyBytes())
	sessionStore.Options(sessions.Options{
		Path:     cfg.Session.Path,
		Domain:   cfg.Session.Domain,
		MaxAge:   cfg.Session.MaxAge,
		Secure:   cfg.Session.Secure,
		HttpOnly: cfg.Session.HttpOnly,
		SameSite: cfg.Session.SameSite,
	})
	sessionMiddleware := sessions.Sessions("session", sessionStore)

	gzipMiddleware := gzip.Gzip(gzip.DefaultCompression)

	engine := gin.New()

	initLogging(engine)

	engine.Use(
		gin.Recovery(),
		secureMiddleware,
		sessionMiddleware,
		gzipMiddleware,
		middleware.Context(dbPool),
	)

	engine.HTMLRender = &html.Render{Fallback: engine.HTMLRender}

	static.Router(engine)
	home.Router(engine)
	auth.Router(engine, csrfMiddleware)

	server := http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Server.Port),
		ReadHeaderTimeout: 10 * time.Second,
		Handler:           engine,
	}

	log.Printf("Starting server on port %d", cfg.Server.Port)
	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}
}
