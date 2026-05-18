package middleware

import (
	"net/http"
	"strings"

	"github.com/Tedra-ez/AdvancedProgramming_Final/auth-service/internal/services"
	"github.com/gin-gonic/gin"
)

const cookieName = "auth_token"

func Auth(authSvc *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if len(tokenStr) > 7 && strings.HasPrefix(tokenStr, "Bearer ") {
			tokenStr = tokenStr[7:]
		}
		if tokenStr == "" {
			if cookie, _ := c.Cookie(cookieName); cookie != "" {
				tokenStr = cookie
			}
		}
		if tokenStr == "" {
			c.Next()
			return
		}
		user, err := authSvc.ParseToken(c.Request.Context(), tokenStr)
		if err != nil {
			c.Next()
			return
		}
		c.Set("user_id", user["id"])
		c.Set("user_role", user["role"])
		c.Set("user_email", user["email"])
		c.Set("user_name", user["name"])
		c.Next()
	}
}

func RequireAuth(c *gin.Context) {
	if _, ok := c.Get("user_id"); !ok {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}
	c.Next()
}

func RequireAdmin(c *gin.Context) {
	role, _ := c.Get("user_role")
	if role != "admin" && role != "administrator" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	c.Next()
}
