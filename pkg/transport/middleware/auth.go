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
// This calls CurrentUser first, so you don't need to chain both.
func Authenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, exists := c.Get("user"); !exists {
			CurrentUser()(c)
		}

		if _, exists := c.Get("user"); !exists {
			hx := c.MustGet("htmx").(*HTMXHelper)
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

// CurrentUser middleware func to set the current active user.
func CurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		queries := c.MustGet("queries").(*dbx.Queries)

		userID, ok := session.Get("user_id").([16]byte)
		if !ok {
			return
		}

		userUUID := pgtype.UUID{Bytes: userID, Valid: true}
		user, err := queries.GetUserByID(c.Request.Context(), userUUID)
		if err != nil || !user.IsActive {
			return
		}

		sloggin.AddCustomAttributes(c, slog.String("user", user.Email))

		c.Set("user", user)
	}
}
