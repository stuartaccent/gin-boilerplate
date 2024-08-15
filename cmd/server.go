package cmd

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"gin.go.dev/embedded"
	"gin.go.dev/internal/middleware"
	"gin.go.dev/internal/renderer"
	"gin.go.dev/internal/routing"
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
)

var cmdServer = &cobra.Command{
	Use:   "server",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		// set up the db pool
		dbPool, err := pgxpool.New(context.Background(), fmt.Sprintf(
			"host=%s port=%v user=%s password=%s database=%s sslmode=%s",
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Db,
			cfg.Database.SslMode,
		))
		if err != nil {
			log.Fatalf("Unable to create connection pool: %v\n", err)
		}
		defer dbPool.Close()

		// set up gin mode
		gin.SetMode(cfg.Server.Mode.ToGinMode())

		// create new engine
		g := gin.New()
		g.Use(gin.Recovery())

		// Conditionally use Metrics middleware
		if cmd.Name() == "monitor" {
			g.Use(middleware.MetricsMiddleware())
		} else {
			g.Use(gin.Logger())
		}

		// csrf middleware
		csrfMiddleware := csrf.Middleware(csrf.Options{
			Secret: cfg.Security.CsrfSecret,
			ErrorFunc: func(c *gin.Context) {
				c.String(400, "CSRF token mismatch")
				c.Abort()
			},
		})

		// secure middleware
		g.Use(secure.New(secure.Config{
			AllowedHosts:          cfg.Security.AllowedHosts,
			STSSeconds:            cfg.Security.StsSeconds,
			STSIncludeSubdomains:  cfg.Security.StsIncludeSubdomains,
			FrameDeny:             cfg.Security.FrameDeny,
			ContentTypeNosniff:    cfg.Security.ContentTypeNosniff,
			BrowserXssFilter:      cfg.Security.BrowserXSSFilter,
			ContentSecurityPolicy: cfg.Security.ContentSecurityPolicy,
		}))

		// session middleware
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
		g.Use(sessions.Sessions("session", sessionStore))

		// custom middleware
		g.Use(middleware.Database(dbPool))
		g.Use(middleware.HTMX())

		// html renderer
		g.HTMLRender = &renderer.HTMLRenderer{Fallback: g.HTMLRender}

		// static
		staticFS, err := fs.Sub(embedded.Static, "static")
		if err != nil {
			log.Fatalf("Unable to load static files: %v", err)
		}
		g.StaticFS("/static", http.FS(staticFS))

		// ui css
		stylesheet := styles.NewStyleSheet()

		g.Handle("GET", "/ui.css", func(c *gin.Context) {
			c.Writer.Header().Set("Content-Type", "text/css")
			if err := stylesheet.CSS(c.Writer); err != nil {
				log.Printf("error writing style: %v", err)
			}
		})

		// routes
		routing.NewMainRouter(g)
		routing.NewAuthRouter(g, csrfMiddleware)

		// start the server
		server := http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
			Handler: g,
		}

		log.Printf("Starting server on port %d", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Fatal(err)
			}
		}
	},
}

func decodeHex(hexStr string) []byte {
	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		log.Fatal(err)
	}
	return decoded
}
