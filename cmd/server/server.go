package main

import (
	"context"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"gin.go.dev/internal/routing"
	"gin.go.dev/internal/webx"
	"github.com/gin-contrib/secure"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	csrf "github.com/utrack/gin-csrf"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	helpFlag           = flag.Bool("help", false, "Display help information")
	port               = flag.Int("port", 80, "The server port")
	dbDns              = flag.String("db-dns", getEnv("DB_DNS", "default_db_dns"), "The database DNS")
	sessionAuthDefault = "bc5d8c0382c39241a6b7d7394ff8af49243ed64e154f0f1500771d879b70d689"
	sessionEncDefault  = "8fbb8e195bd4a8dce59f7eabc2283196"
	csrfSecretDefault  = "some-csrf-secret"
	sessionAuthKey     = flag.String("session-auth-key", sessionAuthDefault, "The hex representation of the session authentication key")
	sessionEncKey      = flag.String("session-enc-key", sessionEncDefault, "The hex representation of the session encryption key")
	csrfSecret         = flag.String("csrf-secret", csrfSecretDefault, "The csrf secret")
)

func getEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func main() {
	flag.Parse()
	if *helpFlag {
		flag.PrintDefaults()
		os.Exit(0)
	}

	dbPool := setupPool()
	defer dbPool.Close()

	e := setupGin(dbPool)
	startServer(e)
}

func setupPool() *pgxpool.Pool {
	dbPool, err := pgxpool.New(context.Background(), *dbDns)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	return dbPool
}

func setupGin(dbPool *pgxpool.Pool) *gin.Engine {
	g := gin.Default()

	csrfMiddleware := csrf.Middleware(csrf.Options{
		Secret: *csrfSecret,
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	})
	sessionKey := decodeHex(sessionAuthKey)
	sessionEnc := decodeHex(sessionEncKey)
	sessionStore := cookie.NewStore(sessionKey, sessionEnc)
	sessionStore.Options(sessions.Options{
		Path:     "/",
		Domain:   "",
		MaxAge:   86400 * 30,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	g.Use(secure.New(secure.Config{
		AllowedHosts:          []string{},
		STSSeconds:            86400 * 365,
		STSIncludeSubdomains:  true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
	}))
	g.Use(sessions.Sessions("session", sessionStore))
	g.Use(webx.NewCustomContext(dbPool))

	g.Static("/static", "./static")

	g.LoadHTMLFiles(getTemplates()...)

	routing.NewMainRouter(g)
	routing.NewAuthRouter(g, csrfMiddleware)

	return g
}

func decodeHex(hexStr *string) []byte {
	decoded, err := hex.DecodeString(*hexStr)
	if err != nil {
		log.Fatal(err)
	}
	return decoded
}

func getTemplates() []string {
	var templates []string
	err := filepath.Walk("templates", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if filepath.Ext(path) == ".gohtml" {
			templates = append(templates, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return templates
}

func startServer(e *gin.Engine) {
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: e,
	}

	log.Printf("Starting server on port %d", *port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
