package routing

import (
	"encoding/gob"
	"gin.go.dev/internal/db"
	"gin.go.dev/internal/webx"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
	"log"
	"net/http"
	"strings"
)

func init() {
	gob.Register([16]byte{})
}

// LoginCredentials used in the login validation
type LoginCredentials struct {
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required,min=6"`
}

// AuthRouter route handler.
type AuthRouter struct {
}

// NewAuthRouter create a new AuthRouter.
func NewAuthRouter(e *gin.Engine, csrf gin.HandlerFunc) {
	r := AuthRouter{}
	g := e.Group("/auth")
	g.GET("/login", csrf, r.loginForm)
	g.POST("/login", csrf, webx.ContentTypes("application/x-www-form-urlencoded"), r.login)
	g.GET("/logout", r.logout)
	g.GET("/user-menu", webx.CurrentUser(), r.userMenu)
}

// loginForm get the login form
func (r *AuthRouter) loginForm(c *gin.Context) {
	cc := c.MustGet("custom").(*webx.CustomContext)
	cc.Session.Clear()
	c.HTML(http.StatusOK, "loginPage", gin.H{
		"Csrf": csrf.GetToken(c),
	})
}

// login the user from the login form then redirect to home
func (r *AuthRouter) login(c *gin.Context) {
	cc := c.MustGet("custom").(*webx.CustomContext)
	invalid := func() {
		c.HTML(http.StatusOK, "loginPage", gin.H{
			"Error": "Invalid email address or password",
			"Csrf":  csrf.GetToken(c),
		})
	}

	var credentials LoginCredentials
	if err := c.Bind(&credentials); err != nil {
		invalid()
		return
	}

	if err := cc.Validator.Validate(credentials); err != nil {
		invalid()
		return
	}

	email := strings.ToLower(credentials.Email)
	user, err := cc.Queries.GetUserByEmail(c.Request.Context(), email)
	if err != nil {
		invalid()
		return
	}

	password := []byte(credentials.Password)
	if !webx.CheckPassword(user.HashedPassword, password) {
		invalid()
		return
	}

	cc.Session.Set("user_id", user.ID.Bytes)
	if err := cc.Session.Save(); err != nil {
		log.Printf("session save error: %v\n", err)
		invalid()
		return
	}

	c.Redirect(http.StatusFound, "/")
}

// logout the user and redirect to login
func (r *AuthRouter) logout(c *gin.Context) {
	cc := c.MustGet("custom").(*webx.CustomContext)
	cc.Session.Clear()
	cc.Session.Options(sessions.Options{MaxAge: -1})
	_ = cc.Session.Save()
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
