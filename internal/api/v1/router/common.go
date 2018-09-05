package router

import (
	"github.com/gin-gonic/gin"

	"github.com/chrootlogin/go-docstore/internal/api/v1/user"
)

func InitRouter(g *gin.RouterGroup) {
	g.GET("/user/*username", user.GetUserHandler)
}
