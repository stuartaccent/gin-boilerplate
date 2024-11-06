package static

import (
	"embed"
	"github.com/gin-gonic/gin"
	"io/fs"
	"log"
	"net/http"
)

var (
	//go:embed public/*
	Public embed.FS
)

// Router create a new static router.
func Router(e *gin.Engine) {
	e.StaticFS("/static", staticFS())
}

// staticFS returns the static file system.
func staticFS() http.FileSystem {
	s, err := fs.Sub(Public, "public")
	if err != nil {
		log.Fatalf("Unable to load static files: %v", err)
	}
	return http.FS(s)
}
