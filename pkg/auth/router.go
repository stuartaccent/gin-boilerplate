package auth

import (
	"encoding/gob"
	"gin.go.dev/pkg/storage/db/dbx"
	"gin.go.dev/pkg/transport/middleware"
	"gin.go.dev/pkg/ui/components"
	"gin.go.dev/pkg/ui/pages"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	csrf "github.com/stuartaccent/gin-csrf"
	"golang.org/x/time/rate"
	"net/http"
	"strings"
)

func init() {
	gob.Register([16]byte{})
}

// LoginCredentials used in the login validation
type LoginCredentials struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required,min=6"`
}

// Router create a new Router.
func Router(e *gin.Engine, csrf gin.HandlerFunc) {
	limiter := middleware.RateLimiter(rate.Limit(2), 5)
	allowForm := middleware.AllowContentType("application/x-www-form-urlencoded")
	auth := middleware.Authenticated()
	g := e.Group("/auth")
	{
		g.GET("/login", csrf, loginForm)
		g.POST("/login", limiter, allowForm, csrf, login)
		g.GET("/logout", logout)
		g.GET("/user-menu", auth, userMenu)
	}
}

// loginForm get the login form
func loginForm(c *gin.Context) {
	session := c.MustGet("session").(sessions.Session)
	session.Clear()
	c.HTML(http.StatusOK, "", pages.Login(pages.LoginData{
		Csrf: csrf.GetToken(c),
	}))
}

// login the user from the login form then redirect to home
func login(c *gin.Context) {
	ctx := c.Request.Context()
	hx := c.MustGet("htmx").(*middleware.HTMX)
	queries := c.MustGet("queries").(*dbx.Queries)
	session := c.MustGet("session").(sessions.Session)

	invalid := func() {
		c.HTML(http.StatusUnprocessableEntity, "", pages.Login(pages.LoginData{
			Error: "invalid email address or password",
			Csrf:  csrf.GetToken(c),
		}))
	}

	var credentials LoginCredentials
	if err := c.ShouldBind(&credentials); err != nil {
		invalid()
		return
	}

	email := strings.ToLower(credentials.Email)
	user, err := queries.GetUserByEmail(ctx, email)
	if err != nil || !user.IsActive {
		invalid()
		return
	}

	password := []byte(credentials.Password)
	if !CheckPassword(user.HashedPassword, password) {
		invalid()
		return
	}

	session.Set("user_id", user.ID.Bytes)
	if err = session.Save(); err != nil {
		_ = c.Error(err)
		invalid()
		return
	}

	hx.SetRedirect("/")
	c.Status(http.StatusOK)
}

// logout the user then redirect to login
func logout(c *gin.Context) {
	session := c.MustGet("session").(sessions.Session)
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1})
	if err := session.Save(); err != nil {
		_ = c.Error(err)
	}
	c.Redirect(http.StatusFound, "/auth/login")
}

// userMenu the user menu in the header.
func userMenu(c *gin.Context) {
	_, open := c.GetQuery("open")
	c.HTML(http.StatusOK, "", components.UserMenu(open))
}
