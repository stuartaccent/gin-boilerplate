package routing

import (
	"gin.go.dev/public"
	"gin.go.dev/ui/styles"
	"github.com/gin-gonic/gin"
	"io/fs"
	"log"
	"net/http"
)

// StaticRouter is a router for static files.
type StaticRouter struct {
	stylesheet *styles.StyleSheet
}

// NewStaticRouter create a new static router.
func NewStaticRouter(e *gin.Engine) {
	r := &StaticRouter{stylesheet: styles.NewStyleSheet()}
	e.StaticFS("/static", staticFS())
	e.GET("/ui.css", r.uiCss)
}

// staticFS returns the static file system.
func staticFS() http.FileSystem {
	s, err := fs.Sub(public.Static, "static")
	if err != nil {
		log.Fatalf("Unable to load static files: %v", err)
	}
	return http.FS(s)
}

// uiCss handles the request for the UI stylesheet.
func (r *StaticRouter) uiCss(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/css")
	if err := r.stylesheet.CSS(c.Writer); err != nil {
		log.Printf("error writing style: %v", err)
	}
}
