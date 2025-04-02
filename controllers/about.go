package controllers

import (
	"example/go-htmx/views"

	"github.com/gin-gonic/gin"
)

type AboutRouterParams struct {
}

type aboutRouter struct {
}

func NewAboutRouter(params AboutRouterParams) *aboutRouter {
	return &aboutRouter{}
}

func (r *aboutRouter) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/", r.getAbout)
}

func (r *aboutRouter) getAbout(c *gin.Context) {
	renderer := views.NewRenderer(c)
	renderer.About().Render(c, c.Writer)
}
