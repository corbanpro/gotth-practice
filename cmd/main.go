package main

import (
	"example/go-htmx/controllers"
	"example/go-htmx/store"

	"github.com/gin-gonic/gin"
)

// TODO: learn htmx and templ again.
// TODO: add middleware to folder

func main() {
	db := store.InitDb()
	redis := store.InitRedis(&gin.Context{})

	todoStore := store.NewTodoStore(db)
	userStore := store.NewUserStore(db)
	sessionStore := store.NewSessionStore(redis, &gin.Context{})

	router := gin.Default()
	router.Static("assets", "./assets")
	router.Any("/", func(c *gin.Context) {
		c.Redirect(303, "/home")
	})
	router.Use(func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
	})
	router.Use(func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			c.Set("user", nil) // No user
			c.Next()
			return
		}

		// Retrieve user from DB
		var user store.User

		session, err := sessionStore.GetById(sessionID)
		if err != nil {
			c.Set("user", nil) // No user
			c.Next()
			return
		}

		user, err = userStore.GetByUsername(session.Username)
		if err != nil {
			c.Set("user", nil) // No user
			c.Next()
			return
		}

		c.Set("user", &user)

		c.Next()
	})

	authRouterGroup := router.Group("/auth")
	authRouter := controllers.NewAuthRouter(controllers.AuthRouterParams{UserStore: userStore, SessionStore: sessionStore})
	authRouter.RegisterRoutes(authRouterGroup)

	homeRouterGroup := router.Group("/home")
	homeRouter := controllers.NewHomeRouter(controllers.HomeRouterParams{})
	homeRouter.RegisterRoutes(homeRouterGroup)

	aboutRouterGroup := router.Group("/about")
	aboutRouter := controllers.NewAboutRouter(controllers.AboutRouterParams{})
	aboutRouter.RegisterRoutes(aboutRouterGroup)

	todoRouterGroup := router.Group("/todo")
	todoRouter := controllers.NewTodoRouter(controllers.TodoRouterParams{TodoStore: todoStore})
	todoRouter.RegisterRoutes(todoRouterGroup)

	router.Run(":4000")
}
