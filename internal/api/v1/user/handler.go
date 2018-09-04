package user

import (
	"net/http"
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/chrootlogin/go-docstore/internal/common"
	"github.com/chrootlogin/go-docstore/internal/store"
)

type apiResponse struct {
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	Permissions []string `json:"permissions"`
}

func GetUserHandler(c *gin.Context) {
	userName := c.Param("username")
	if len(userName) <= 3 {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.ApiResponse{Message: common.WrongAPIUsageError})
		return
	}

	// remove first character because it's always /
	userName = trimLeftChar(userName)

	user, err := store.Users().Get(userName)
	if err != nil {
		if err == store.ErrUserNotExist {
			c.AbortWithStatusJSON(http.StatusNotFound, common.ApiResponse{Message: fmt.Sprintf("User not found: %s", userName)})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, common.ApiResponse{Message: fmt.Sprintf("Can't get user list: %s", err.Error())})
		return
	}

	resp := apiResponse{
		Username:    user.Username,
		Email:       user.Email,
		Permissions: user.Permissions,
	}

	c.JSON(http.StatusOK, resp)
}

// https://stackoverflow.com/a/48798875
func trimLeftChar(s string) string {
	for i := range s {
		if i > 0 {
			return s[i:]
		}
	}
	return s[:0]
}