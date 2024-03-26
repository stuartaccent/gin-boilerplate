package routing

import (
	"gin.go.dev/internal/middleware"
	"gin.go.dev/internal/renderer"
	"gin.go.dev/ui/pages"
	"github.com/gin-gonic/gin"
)

// NewMainRouter create a new MainRouter.
func NewMainRouter(e *gin.Engine) {
	g := e.Group("/", middleware.Authenticated())
	g.GET("/", index)
}

// index root page endpoint.
func index(c *gin.Context) {
	ctx := c.Request.Context()
	h := renderer.New(ctx, 200, pages.Home())
	c.Render(200, h)
}
