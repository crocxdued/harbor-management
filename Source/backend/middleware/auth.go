package middleware

import (
	"harbor/models"
	"harbor/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth(svc *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{Error: "требуется авторизация"})
			return
		}
		claims, err := svc.ParseToken(strings.TrimPrefix(h, "Bearer "))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{Error: "недействительный токен"})
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		for _, r := range roles {
			if role == r { c.Next(); return }
		}
		c.AbortWithStatusJSON(http.StatusForbidden, models.ErrorResponse{Error: "недостаточно прав"})
	}
}
