package middleware

import (
	"goBlogApp/models"
	"net/http"
	"fmt"

	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
)

func Authorize() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, ok := ctx.Get("sub")
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		
		enforcer := casbin.NewEnforcer("config/acl_model.conf", "config/policy.csv")
		fmt.Println("- enforcer -", enforcer)
		ok = enforcer.Enforce(user.(*models.User), ctx.Request.URL.Path, ctx.Request.Method)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}
		ctx.Next()
	}
}
