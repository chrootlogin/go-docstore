package user

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/chrootlogin/go-docstore/internal/common"
	"github.com/chrootlogin/go-docstore/internal/store"
)

// try existing user
func TestGetUserHandler(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()

	u := common.User{
		Username: "test-user",
		Email:    "test@example.com",
	}

	err := store.Users().Add(u)
	if assert.NoError(err) {
		r := gin.Default()
		r.GET("/user/*username", GetUserHandler)

		req, _ := http.NewRequest("GET", "/user/"+u.Username, nil)
		r.ServeHTTP(w, req)

		if assert.Equal(w.Code, http.StatusOK) {
			data, err := ioutil.ReadAll(w.Body)
			if assert.NoError(err) {
				var resp apiResponse
				err = json.Unmarshal(data, &resp)
				if assert.NoError(err) {
					assert.Equal(u.Username, resp.Username)
					assert.Equal(u.Email, resp.Email)
				}
			}
		}
	}
}

// try non existing user
func TestGetUserHandler2(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()

	r := gin.Default()
	r.GET("/user/*username", GetUserHandler)

	req, _ := http.NewRequest("GET", "/user/not-existing", nil)
	r.ServeHTTP(w, req)

	assert.Equal(http.StatusNotFound, w.Code)
}

// try with to short username
func TestGetUserHandler3(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()

	r := gin.Default()
	r.GET("/user/*username", GetUserHandler)

	req, _ := http.NewRequest("GET", "/user/ab", nil)
	r.ServeHTTP(w, req)

	assert.Equal(http.StatusBadRequest, w.Code)
}
