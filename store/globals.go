package store

import (
	"github.com/gin-gonic/gin"
)

type Globals struct {
	User *User
}

func GetGlobals(c *gin.Context) *Globals {
	userVal, exists := c.Get("user")
	var user *User
	if !exists || userVal == nil {
		user = nil
	} else {
		user = userVal.(*User)
	}

	return &Globals{
		User: user,
	}
}

func (g *Globals) GetUser() (user *User, exists bool) {
	if g.User == nil {
		return &User{}, false
	}
	return g.User, true
}
