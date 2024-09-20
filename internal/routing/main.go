package routing

import (
	"gin.go.dev/internal/middleware"
	"gin.go.dev/ui/pages"
	"github.com/gin-gonic/gin"
	"net/http"
)

// NewMainRouter create a new MainRouter.
func NewMainRouter(e *gin.Engine) {
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
