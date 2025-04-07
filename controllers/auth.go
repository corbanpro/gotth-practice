package controllers

import (
	"example/go-htmx/request"
	"example/go-htmx/store"
	"example/go-htmx/views"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthRouterParams struct {
	UserStore    store.UserStore
	SessionStore store.SessionStore
}

type authRouter struct {
	userStore    store.UserStore
	sessionStore store.SessionStore
}

func NewAuthRouter(params AuthRouterParams) *authRouter {
	if params.UserStore == nil {
		panic("UserStore is required")
	} else if params.SessionStore == nil {
		panic("SessionStore is required")
	}
	return &authRouter{
		userStore:    params.UserStore,
		sessionStore: params.SessionStore,
	}
}

func (r *authRouter) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/login", r.getLogin)
	router.GET("/register", r.getRegister)
	router.POST("/login", r.postLogin)
	router.POST("/register", r.postRegister)
	router.GET("/logout", r.getLogout)
}

func (r *authRouter) getLogin(c *gin.Context) {
	user, _ := request.GetUser(c, r.userStore)
	views.LoginPage(user).Render(c, c.Writer)
}

func (r *authRouter) getRegister(c *gin.Context) {
	user, _ := request.GetUser(c, r.userStore)
	views.RegisterPage(user).Render(c, c.Writer)
}

func (r *authRouter) postLogin(c *gin.Context) {
	c.Request.ParseForm()
	formData := c.Request.PostForm
	username := formData.Get("username")
	password := formData.Get("password")

	if username == "" || password == "" {
		c.Status(http.StatusBadRequest)
		views.LoginError("Please fill out all fields").Render(c, c.Writer)
		return
	}

	loginValid, err := r.userStore.ValidateLogin(username, password)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		views.LoginError("Error. Please try again later.").Render(c, c.Writer)
		return
	}

	if !loginValid {
		c.Status(http.StatusUnauthorized)
		views.LoginError("Username or password incorrect").Render(c, c.Writer)
		return
	}

	session, err := r.sessionStore.Create(username)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		views.LoginError("Error. Please try again later.").Render(c, c.Writer)
		return
	}

	c.SetCookie("session_id", session.Id, 3600, "/", "", false, true)

	c.Header("HX-Redirect", "/home")
}

func (r *authRouter) postRegister(c *gin.Context) {
	c.Request.ParseForm()

	formData := c.Request.PostForm
	username := formData.Get("username")
	password := formData.Get("password")
	confirmPassword := formData.Get("confirm_password")
	firstName := formData.Get("first_name")
	lastName := formData.Get("last_name")

	if username == "" || firstName == "" || lastName == "" || password == "" || confirmPassword == "" {
		c.Status(http.StatusBadRequest)
		views.RegisterError("All fields are required.").Render(c, c.Writer)
		return
	}

	if len(password) < 8 {
		c.Status(http.StatusBadRequest)
		views.RegisterError("Password must be at least 8 characters.").Render(c, c.Writer)
		return
	}

	if len(username) < 4 {
		c.Status(http.StatusBadRequest)
		views.RegisterError("Username must be at least 4 characters.").Render(c, c.Writer)
		return
	}

	if password != confirmPassword {
		c.Status(http.StatusBadRequest)
		views.RegisterError("Passwords do not match.").Render(c, c.Writer)
		return
	}

	existingUser, err := r.userStore.GetByUsername(username)

	if err == nil || existingUser.Username != "" {
		c.Status(http.StatusConflict)
		views.RegisterError("User already exists. Please log in").Render(c, c.Writer)
		return
	}

	user, err := r.userStore.Create(username, password, firstName, lastName)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		views.RegisterError("Error creating user. Please try again later.").Render(c, c.Writer)
		return
	}

	session, err := r.sessionStore.Create(user.Username)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		views.RegisterError("Error creating session. Please try again later.").Render(c, c.Writer)
		return
	}

	c.SetCookie("session_id", session.Id, 3600, "/", "", false, true)

	c.Header("HX-Redirect", "/home")
}

func (r *authRouter) getLogout(c *gin.Context) {
	sessionId, err := c.Cookie("session_id")

	if err == nil {
		r.sessionStore.Delete(sessionId)
		c.SetCookie("session_id", "", -1, "/", "", false, true)
	}

	c.Redirect(http.StatusSeeOther, "/auth/login")
}
