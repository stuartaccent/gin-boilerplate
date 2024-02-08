package webx

import (
	"gin.go.dev/internal/db"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CustomContext hold various functionality for gin usage in gin.HandlerFunc
// such as Queries and Session.
type CustomContext struct {
	*gin.Context
	Postgres  *pgxpool.Pool
	Queries   *db.Queries
	Session   sessions.Session
	Validator *CustomValidator
}

// NewCustomContext gin handler to create a new CustomContext.
func NewCustomContext(dbPool *pgxpool.Pool) gin.HandlerFunc {
	validator := NewCustomValidator()
	queries := db.New(dbPool)
	return func(c *gin.Context) {
		session := sessions.Default(c)
		cc := &CustomContext{
			Context:   c,
			Postgres:  dbPool,
			Queries:   queries,
			Session:   session,
			Validator: validator,
		}
		c.Set("custom", cc)
		c.Next()
	}
}
