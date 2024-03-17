package routing

import (
	"gin.go.dev/components"
	"gin.go.dev/internal/middleware"
	"gin.go.dev/internal/renderer"
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
	h := renderer.New(ctx, 200, components.EmptyPage("Home"))
	c.Render(200, h)
}
