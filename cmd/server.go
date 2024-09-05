package cmd

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"gin.go.dev/internal/middleware"
	"gin.go.dev/internal/renderer"
	"gin.go.dev/internal/routing"
	"gin.go.dev/public"
	"gin.go.dev/ui/styles"
	"github.com/gin-contrib/secure"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
	csrf "github.com/utrack/gin-csrf"
	"io/fs"
	"log"
	"net/http"
	"time"
)

var (
	monitorDuration time.Duration
	monitorLines    int
)

var cmdMonitor = &cobra.Command{
	Use:   "monitor",
	Short: "Start the server with monitoring",
	Run: func(cmd *cobra.Command, args []string) {
		gin.SetMode(cfg.Server.Mode.ToGinMode())

		engine := gin.New()
		engine.Use(gin.Recovery(), middleware.MetricsMiddleware())

		go func() {
			for {
				<-time.After(monitorDuration)
				middleware.MetricsResults.WriteMetrics(monitorLines)
			}
		}()

		runServer(engine)
	},
}

var cmdServer = &cobra.Command{
	Use:   "server",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		gin.SetMode(cfg.Server.Mode.ToGinMode())

		engine := gin.Default()
		runServer(engine)
	},
}

func init() {
	cmdMonitor.Flags().DurationVarP(&monitorDuration, "duration", "d", time.Second*5, "Duration between metrics collection")
	cmdMonitor.Flags().IntVarP(&monitorLines, "lines", "l", 50, "Number of lines to print")
	cmdServer.AddCommand(cmdMonitor)
}

func decodeHex(hexStr string) []byte {
	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		log.Fatal(err)
	}
	return decoded
}

func runServer(engine *gin.Engine) {
	dbPool, err := pgxpool.New(context.Background(), cfg.Database.URL().String())
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer dbPool.Close()

	csrfMiddleware := csrf.Middleware(csrf.Options{
		Secret: cfg.Security.CsrfSecret,
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	})

	engine.Use(secure.New(secure.Config{
		AllowedHosts:          cfg.Security.AllowedHosts,
		STSSeconds:            cfg.Security.StsSeconds,
		STSIncludeSubdomains:  cfg.Security.StsIncludeSubdomains,
		FrameDeny:             cfg.Security.FrameDeny,
		ContentTypeNosniff:    cfg.Security.ContentTypeNosniff,
		BrowserXssFilter:      cfg.Security.BrowserXSSFilter,
		ContentSecurityPolicy: cfg.Security.ContentSecurityPolicy,
	}))

	sessionKey := decodeHex(cfg.Session.Key)
	sessionEnc := decodeHex(cfg.Session.EncKey)
	sessionStore := cookie.NewStore(sessionKey, sessionEnc)
	sessionStore.Options(sessions.Options{
		Path:     cfg.Session.Path,
		Domain:   cfg.Session.Domain,
		MaxAge:   cfg.Session.MaxAge,
		Secure:   cfg.Session.Secure,
		HttpOnly: cfg.Session.HttpOnly,
		SameSite: cfg.Session.SameSite,
	})
	engine.Use(sessions.Sessions("session", sessionStore))

	engine.Use(middleware.Database(dbPool))
	engine.Use(middleware.HTMX())

	engine.HTMLRender = &renderer.HTMLRenderer{Fallback: engine.HTMLRender}

	staticFS, err := fs.Sub(public.Static, "static")
	if err != nil {
		log.Fatalf("Unable to load static files: %v", err)
	}
	engine.StaticFS("/static", http.FS(staticFS))

	stylesheet := styles.NewStyleSheet()
	engine.Handle("GET", "/ui.css", func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/css")
		if err := stylesheet.CSS(c.Writer); err != nil {
			log.Printf("error writing style: %v", err)
		}
	})

	routing.NewMainRouter(engine)
	routing.NewAuthRouter(engine, csrfMiddleware)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: engine,
	}

	log.Printf("Starting server on port %d", cfg.Server.Port)
	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}
}
