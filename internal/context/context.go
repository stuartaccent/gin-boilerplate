package context

import (
	"gin.go.dev/internal/db"
	"gin.go.dev/internal/htmx"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// GinContext middleware func to set the context for the request.
func GinContext(dbPool *pgxpool.Pool) gin.HandlerFunc {
	queries := db.New(dbPool)
	return func(c *gin.Context) {
		hx := htmx.New(c.Request, c.Writer)
		session := sessions.Default(c)
		c.Set("htmx", hx)
		c.Set("postgres", dbPool)
		c.Set("queries", queries)
		c.Set("session", session)
		c.Next()
	}
}
