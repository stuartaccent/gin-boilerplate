package routing

import (
	"encoding/gob"
	"log"
	"strings"

	"gin.go.dev/internal/crypt"
	"gin.go.dev/internal/db"
	"gin.go.dev/internal/htmx"
	"gin.go.dev/internal/middleware"
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
	session := sessions.Default(c)
	session.Clear()
	c.HTML(200, "", pages.Login(pages.LoginData{
		Error: "",
		Csrf:  csrf.GetToken(c),
	}))
}

// login the user from the login form then redirect to home
func login(c *gin.Context) {
	ctx := c.Request.Context()
	hx := c.MustGet("htmx").(*htmx.Helper)
	queries := c.MustGet("queries").(*db.Queries)
	session := sessions.Default(c)

	invalid := func() {
		c.HTML(422, "", pages.Login(pages.LoginData{
			Error: "Invalid email address or password",
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
	if !crypt.CheckPassword(user.HashedPassword, password) {
		invalid()
		return
	}

	session.Set("user_id", user.ID.Bytes)
	if err = session.Save(); err != nil {
		log.Printf("session save error: %v\n", err)
		invalid()
		return
	}

	hx.SetRedirect("/")
	c.Status(200)
}

// logout the user then redirect to login
func logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1})
	_ = session.Save()
	c.Redirect(302, "/auth/login")
}

// userMenu the user menu in the header.
func userMenu(c *gin.Context) {
	_, open := c.GetQuery("open")
	if open {
		c.HTML(200, "", components.UserMenuOpen())
	} else {
		c.HTML(200, "", components.UserMenuClosed())
	}
}
