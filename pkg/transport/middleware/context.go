package middleware

import (
	"gin.go.dev/pkg/storage/db/dbx"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Context middleware func to set the app context.
func Context(postgres *pgxpool.Pool) gin.HandlerFunc {
	queries := dbx.New(postgres)

	return func(c *gin.Context) {
		htmx := &HTMX{Request: c.Request, Response: c.Writer}

		c.Set("htmx", htmx)
		c.Set("postgres", postgres)
		c.Set("queries", queries)
		c.Set("session", sessions.Default(c))

		c.Next()
	}
}
