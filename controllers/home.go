package controllers

import (
	"example/go-htmx/views"

	"github.com/gin-gonic/gin"
)

type HomeRouterParams struct {
	Renderer views.Renderer
}

type homeRouter struct {
	renderer views.Renderer
}

func NewHomeRouter(params HomeRouterParams) *homeRouter {
	return &homeRouter{
		renderer: params.Renderer,
	}
}

func (r *homeRouter) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/", r.getHome)
	router.GET("/more", r.getMoreHome)
}

func (*homeRouter) getHome(c *gin.Context) {
	views.NewRenderer(c).Home().Render(c, c.Writer)
}

func (*homeRouter) getMoreHome(c *gin.Context) {
	views.NewRenderer(c).MoreHome().Render(c, c.Writer)
}
