package middleware

import (
	"gin.go.dev/pkg/storage/db/dbx"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	sloggin "github.com/samber/slog-gin"
	"log/slog"
	"net/http"
)

// Authenticated middleware func to ensure logged in, redirects to log-in if not.
func Authenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, exists := c.Get("user"); !exists {
			currentUser()(c)
		}

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

// currentUser middleware func to set the current active user.
func currentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
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
}
