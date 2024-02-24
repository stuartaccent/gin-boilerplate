package webx

import (
	"gin.go.dev/internal/db"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SetGinContext middleware func to set the context for the request.
func SetGinContext(dbPool *pgxpool.Pool) gin.HandlerFunc {
	queries := db.New(dbPool)
	return func(c *gin.Context) {
		session := sessions.Default(c)
		c.Set("postgres", dbPool)
		c.Set("queries", queries)
		c.Set("session", session)
		c.Next()
	}
}
