package request

import (
	"example/go-htmx/store"

	"github.com/gin-gonic/gin"
)

func GetUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get("user_id")
	if exists && username != nil {
		return username.(string), true
	}
	return "", false
}

func GetUser(c *gin.Context, userStore store.UserStore) (*store.User, bool) {
	user, exists := c.Get("user")
	if exists && user != nil {
		return user.(*store.User), true
	}

	userId, exists := c.Get("user_id")

	if exists && userId != nil {
		id := userId.(string)

		user, err := userStore.GetByUsername(id)
		if err != nil {
			return nil, false
		}
		c.Set("user", user)
		return &user, true
	}

	return nil, false
}
