package main

import (
	"example/go-htmx/controllers"
	"example/go-htmx/middleware"
	"example/go-htmx/store"

	"github.com/gin-gonic/gin"
)

func main() {
	// Database and Redis initialization
	db := store.InitDb()
	redis := store.InitRedis(&gin.Context{})

	// Store initialization
	todoStore := store.NewTodoStore(db)
	userStore := store.NewUserStore(db)
	sessionStore := store.NewSessionStore(redis, &gin.Context{})

	// Router initialization
	router := gin.Default()
	router.Static("assets", "./assets")
	router.Any("/", func(c *gin.Context) { c.Redirect(303, "/home") })

	// Middleware
	router.Use(func(c *gin.Context) { c.Header("Content-Type", "text/html; charset=utf-8") })
	router.Use(middleware.AuthMiddleware(sessionStore))

	// Route registration
	authRouterGroup := router.Group("/auth")
	authRouter := controllers.NewAuthRouter(controllers.AuthRouterParams{UserStore: userStore, SessionStore: sessionStore})
	authRouter.RegisterRoutes(authRouterGroup)

	homeRouterGroup := router.Group("/home")
	homeRouter := controllers.NewHomeRouter(controllers.HomeRouterParams{UserStore: userStore})
	homeRouter.RegisterRoutes(homeRouterGroup)

	aboutRouterGroup := router.Group("/about")
	aboutRouter := controllers.NewAboutRouter(controllers.AboutRouterParams{UserStore: userStore})
	aboutRouter.RegisterRoutes(aboutRouterGroup)

	todoRouterGroup := router.Group("/todo")
	todoRouter := controllers.NewTodoRouter(controllers.TodoRouterParams{TodoStore: todoStore, UserStore: userStore})
	todoRouter.RegisterRoutes(todoRouterGroup)

	// Start the server
	router.Run(":4000")
}
