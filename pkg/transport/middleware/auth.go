package middleware

import (
	"gin.go.dev/pkg/storage/db/dbx"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	sloggin "github.com/samber/slog-gin"
	"log/slog"
	"net/http"
	"sync"
)

// Authenticated middleware func to ensure logged in, redirects to log-in if not.
func Authenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		var once sync.Once

		once.Do(func() {
			setCurrentUser(c)
		})

		if _, exists := c.Get("user"); !exists {
			hx := c.MustGet("htmx").(*HTMX)
			if hx.IsHTMXRequest() {
				hx.SetRedirect("/auth/login")
				c.Status(http.StatusNoContent)
			} else {
				c.Redirect(http.StatusFound, "/auth/login")
			}
			c.Abort()
		}
	}
}

// setCurrentUser set the current active user.
func setCurrentUser(c *gin.Context) {
	ctx := c.Request.Context()
	session := c.MustGet("session").(sessions.Session)
	queries := c.MustGet("queries").(*dbx.Queries)

	userID, ok := session.Get("user_id").([16]byte)
	if !ok {
		return
	}

	user, err := queries.GetUserByID(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil || !user.IsActive {
		return
	}

	sloggin.AddCustomAttributes(c, slog.String("user", user.Email))

	c.Set("user", user)
}
