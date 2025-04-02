package views

import (
	"example/go-htmx/store"

	"github.com/gin-gonic/gin"
)

type Renderer struct {
	Globals *store.Globals
	c       *gin.Context
}

func NewRenderer(c *gin.Context) *Renderer {
	return &Renderer{
		Globals: store.GetGlobals(c),
		c:       c,
	}
}
