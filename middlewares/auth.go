package middlewares

import (
	"github.com/gin-gonic/gin"
)

func RequireLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := c.Cookie("user_id")
		if err != nil {
			c.Redirect(302, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}
