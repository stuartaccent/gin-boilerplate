package static

import (
	"embed"
	"gin.go.dev/pkg/ui/styles"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
)

var (
	//go:embed public/*
	Public embed.FS
)

// StaticRouter is a router for static files.
type router struct {
	stylesheet *styles.StyleSheet
}

// Router create a new static router.
func Router(e *gin.Engine) {
	r := &router{stylesheet: styles.NewStyleSheet()}
	e.StaticFS("/static", staticFS())
	e.GET("/ui.css", r.uiCss)
}

// staticFS returns the static file system.
func staticFS() http.FileSystem {
	s, err := fs.Sub(Public, "public")
	if err != nil {
		log.Fatalf("Unable to load static files: %v", err)
	}
	return http.FS(s)
}

// uiCss handles the request for the UI stylesheet.
func (r *router) uiCss(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/css")
	if err := r.stylesheet.CSS(c.Writer); err != nil {
		sloggin.AddCustomAttributes(c, slog.String("error", err.Error()))
		c.Status(http.StatusInternalServerError)
		return
	}
}
