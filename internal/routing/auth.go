package routing

import (
	"encoding/gob"
	"log"
	"net/http"
	"strings"

	"gin.go.dev/internal/db"
	"gin.go.dev/internal/webx"
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

// AuthRouter route handler.
type AuthRouter struct {
}

// NewAuthRouter create a new AuthRouter.
func NewAuthRouter(e *gin.Engine, csrf gin.HandlerFunc) {
	r := AuthRouter{}
	g := e.Group("/auth")
	g.GET("/login", csrf, r.loginForm)
	g.POST("/login", webx.RateLimiter(rate.Limit(2), 5), webx.ContentTypes("application/x-www-form-urlencoded"), csrf, r.login)
	g.GET("/logout", r.logout)
	g.GET("/user-menu", webx.Authenticated(), r.userMenu)
}

// loginForm get the login form
func (r *AuthRouter) loginForm(c *gin.Context) {
	session := c.MustGet("session").(sessions.Session)
	session.Clear()
	c.HTML(http.StatusOK, "loginPage", gin.H{
		"Csrf": csrf.GetToken(c),
	})
}

// login the user from the login form then redirect to home
func (r *AuthRouter) login(c *gin.Context) {
	session := c.MustGet("session").(sessions.Session)
	queries := c.MustGet("queries").(*db.Queries)
	invalid := func() {
		c.HTML(http.StatusOK, "loginPage", gin.H{
			"Error": "Invalid email address or password",
			"Csrf":  csrf.GetToken(c),
		})
	}

	var credentials LoginCredentials
	if err := c.ShouldBind(&credentials); err != nil {
		invalid()
		return
	}

	email := strings.ToLower(credentials.Email)
	user, err := queries.GetUserByEmail(c.Request.Context(), email)
	if err != nil || !user.IsActive {
		invalid()
		return
	}

	password := []byte(credentials.Password)
	if !webx.CheckPassword(user.HashedPassword, password) {
		invalid()
		return
	}

	session.Set("user_id", user.ID.Bytes)
	if err := session.Save(); err != nil {
		log.Printf("session save error: %v\n", err)
		invalid()
		return
	}

	c.Redirect(http.StatusFound, "/")
}

// logout the user and redirect to login
func (r *AuthRouter) logout(c *gin.Context) {
	session := c.MustGet("session").(sessions.Session)
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1})
	_ = session.Save()
	c.Redirect(http.StatusFound, "/auth/login")
}

// userMenu the user menu in the header.
func (r *AuthRouter) userMenu(c *gin.Context) {
	_, open := c.GetQuery("open")
	c.HTML(http.StatusOK, "userMenu", gin.H{
		"User": c.MustGet("user").(db.AuthUser),
		"Open": open,
	})
}
