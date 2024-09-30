package home

import (
	"gin.go.dev/pkg/transport/middleware"
	"gin.go.dev/pkg/ui/pages"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Router create a new Router.
func Router(e *gin.Engine) {
	auth := middleware.Authenticated()
	g := e.Group("/", auth)
	{
		g.GET("/", index)
	}
}

// index root page endpoint.
func index(c *gin.Context) {
	c.HTML(http.StatusOK, "", pages.Home())
}
