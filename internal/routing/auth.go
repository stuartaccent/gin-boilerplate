package routing

import (
	"encoding/gob"
	"fmt"
	"log"
	"strings"

	"gin.go.dev/internal/crypt"
	"gin.go.dev/internal/db"
	"gin.go.dev/internal/middleware"
	"gin.go.dev/internal/renderer"
	"gin.go.dev/ui/components"
	"gin.go.dev/ui/pages"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
	"golang.org/x/time/rate"
)

func init() {
	gob.Register([16]byte{})
}

// LoginCredentials used in the login validation
type LoginCredentials struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required,min=6"`
}

// NewAuthRouter create a new AuthRouter.
func NewAuthRouter(e *gin.Engine, csrf gin.HandlerFunc) {
	g := e.Group("/auth")
	g.GET("/login", csrf, loginForm)
	g.POST("/login", middleware.RateLimiter(rate.Limit(2), 5), middleware.ContentTypes("application/x-www-form-urlencoded"), csrf, login)
	g.GET("/logout", logout)
	g.GET("/user-menu", middleware.Authenticated(), userMenu)
}

// loginForm get the login form
func loginForm(c *gin.Context) {
	ctx := c.Request.Context()
	session := c.MustGet("session").(sessions.Session)
	session.Clear()
	h := renderer.New(ctx, 200, pages.Login(pages.LoginData{
		Error: "",
		Csrf:  csrf.GetToken(c),
	}))
	c.Render(200, h)
}

// login the user from the login form then redirect to home
func login(c *gin.Context) {
	ctx := c.Request.Context()
	queries := c.MustGet("queries").(*db.Queries)
	session := c.MustGet("session").(sessions.Session)

	invalid := func() {
		h := renderer.New(ctx, 200, pages.Login(pages.LoginData{
			Error: "Invalid email address or password",
			Csrf:  csrf.GetToken(c),
		}))
		c.Render(200, h)
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
	if !crypt.CheckPassword(user.HashedPassword, password) {
		invalid()
		return
	}

	session.Set("user_id", user.ID.Bytes)
	if err := session.Save(); err != nil {
		log.Printf("session save error: %v\n", err)
		invalid()
		return
	}

	c.Redirect(302, "/")
}

// logout the user then redirect to login
func logout(c *gin.Context) {
	session := c.MustGet("session").(sessions.Session)
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1})
	_ = session.Save()
	c.Redirect(302, "/auth/login")
}

// userMenu the user menu in the header.
func userMenu(c *gin.Context) {
	user := c.MustGet("user").(db.AuthUser)
	ctx := c.Request.Context()
	name := fmt.Sprint(user.FirstName, " ", user.LastName)
	_, open := c.GetQuery("open")
	if open {
		h := renderer.New(ctx, 200, components.UserMenuOpen(name))
		c.Render(200, h)
	} else {
		h := renderer.New(ctx, 200, components.UserMenuClosed(name))
		c.Render(200, h)
	}
}
