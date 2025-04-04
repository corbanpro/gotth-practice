package controllers

import (
	"example/go-htmx/request"
	"example/go-htmx/store"
	"example/go-htmx/views"

	"github.com/gin-gonic/gin"
)

type HomeRouterParams struct {
	UserStore store.UserStore
}

type homeRouter struct {
	userStore store.UserStore
}

func NewHomeRouter(params HomeRouterParams) *homeRouter {
	if params.UserStore == nil {
		panic("UserStore is required")
	}
	return &homeRouter{
		userStore: params.UserStore,
	}
}

func (r *homeRouter) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/", r.getHome)
	router.GET("/more", r.getMoreHome)
}

func (r *homeRouter) getHome(c *gin.Context) {
	user, _ := request.GetUser(c, r.userStore)
	views.HomePage(user).Render(c, c.Writer)
}

func (*homeRouter) getMoreHome(c *gin.Context) {
	views.MoreHome().Render(c, c.Writer)
}
