package middleware

import (
	"example/go-htmx/store"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(sessionStore store.SessionStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			c.Next()
			return
		}

		session, err := sessionStore.GetById(sessionID)
		if err != nil {
			c.SetCookie("session_id", "", -1, "/", "", false, true)
			c.Next()
			return
		}

		c.Set("user_id", session.Username)

		c.Next()
	}
}

func RequireAuth(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		userVal, exists := c.Get("user_id")
		if !exists || userVal == nil {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		} else {
			handler(c)
		}
	}
}
