package main

import (
	"context"
	"embed"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"gin.go.dev/internal/middleware"

	"gin.go.dev/internal/config"
	"gin.go.dev/internal/renderer"
	"gin.go.dev/internal/routing"
	"github.com/gin-contrib/secure"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	csrf "github.com/utrack/gin-csrf"
)

var (
	helpFlag  = flag.Bool("help", false, "Display help information")
	appConfig = flag.String("app-config", "config.toml", "The path of the app config eg: config.toml")

	//go:embed static/*
	static embed.FS
)

func main() {
	flag.Parse()
	if *helpFlag {
		flag.PrintDefaults()
		os.Exit(0)
	}

	// parse the config file
	cfg, err := config.FromPath(*appConfig)
	if err != nil {
		log.Fatalf("Invalid config: %v", err)
	}

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
	g := gin.Default()

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
	staticFS, err := fs.Sub(static, "static")
	if err != nil {
		log.Fatalf("Unable to load static files: %v", err)
	}
	g.StaticFS("/static", http.FS(staticFS))

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
}

func decodeHex(hexStr string) []byte {
	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		log.Fatal(err)
	}
	return decoded
}
