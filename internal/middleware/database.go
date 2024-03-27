package middleware

import (
	"gin.go.dev/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Database middleware func to set the database pool and queries.
func Database(dbPool *pgxpool.Pool) gin.HandlerFunc {
	queries := db.New(dbPool)
	return func(c *gin.Context) {
		c.Set("postgres", dbPool)
		c.Set("queries", queries)
		c.Next()
	}
}
