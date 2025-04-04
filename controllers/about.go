package controllers

import (
	"example/go-htmx/request"
	"example/go-htmx/store"
	"example/go-htmx/views"

	"github.com/gin-gonic/gin"
)

type AboutRouterParams struct {
	UserStore store.UserStore
}

type aboutRouter struct {
	userStore store.UserStore
}

func NewAboutRouter(params AboutRouterParams) *aboutRouter {
	if params.UserStore == nil {
		panic("UserStore is required")
	}
	return &aboutRouter{
		userStore: params.UserStore,
	}
}

func (r *aboutRouter) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/", r.getAbout)
}

func (r *aboutRouter) getAbout(c *gin.Context) {
	user, _ := request.GetUser(c, r.userStore)
	views.AboutPage(user).Render(c, c.Writer)
}
