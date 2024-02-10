package routing

import (
	"gin.go.dev/internal/webx"
	"github.com/gin-gonic/gin"
	"net/http"
)

// MainRouter route handler.
type MainRouter struct {
}

// NewMainRouter create a new MainRouter.
func NewMainRouter(e *gin.Engine) {
	r := MainRouter{}
	g := e.Group("/", webx.Authenticated())
	g.GET("/", r.index)
}

// index root page endpoint.
func (r *MainRouter) index(c *gin.Context) {
	c.HTML(http.StatusOK, "indexPage", nil)
}
