package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/chrootlogin/go-docstore/internal/api/v1/user"
	"github.com/chrootlogin/go-docstore/internal/common"
	"github.com/chrootlogin/go-docstore/internal/helper"
	"github.com/chrootlogin/go-docstore/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
)

func TestGetAuthMiddleware(t *testing.T) {
	assert := assert.New(t)

	os.Setenv("SESSION_KEY", "not-a-secret-key")
	am := GetAuthMiddleware()

	assert.NotNil(am)
}

// non-existing user
func TestAuthMiddleware_LoginHandler(t *testing.T) {
	assert := assert.New(t)

	apiReq := ApiLogin{
		Username: "non-existing",
		Password: "non-existing",
	}

	data, err := json.Marshal(apiReq)
	if assert.NoError(err) {
		w := httptest.NewRecorder()

		os.Setenv("SESSION_KEY", "not-a-secret-key")
		am := GetAuthMiddleware()

		r := gin.Default()
		r.POST("/user/login", am.LoginHandler)

		req, _ := http.NewRequest("POST", "/user/login", bytes.NewReader(data))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Content-Length", string(len(data)))
		r.ServeHTTP(w, req)

		assert.Equal(http.StatusUnauthorized, w.Code)
	}
}

// existing user
func TestAuthMiddleware_LoginHandler2(t *testing.T) {
	assert := assert.New(t)

	const USERNAME = "admin"
	const PASSWORD = "test1234"

	hash, err := helper.HashPassword(PASSWORD)
	if assert.NoError(err) {
		store.Users().Add(common.User{
			Username:     "admin",
			PasswordHash: hash,
			Email:        "admin@example.org",
		})

		apiReq := ApiLogin{
			Username: USERNAME,
			Password: PASSWORD,
		}

		data, err := json.Marshal(apiReq)
		if assert.NoError(err) {
			w := httptest.NewRecorder()

			os.Setenv("SESSION_KEY", "not-a-secret-key")
			am := GetAuthMiddleware()

			r := gin.Default()
			r.POST("/user/login", am.LoginHandler)

			req, _ := http.NewRequest("POST", "/user/login", bytes.NewReader(data))
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Content-Length", string(len(data)))
			r.ServeHTTP(w, req)

			assert.Equal(http.StatusOK, w.Code)
		}
	}
}

// wrong password
func TestAuthMiddleware_LoginHandler3(t *testing.T) {
	assert := assert.New(t)

	const (
		USERNAME = "admin2"
	    PASSWORD = "test1234"
	)

	hash, err := helper.HashPassword(PASSWORD)
	if assert.NoError(err) {
		store.Users().Add(common.User{
			Username:     USERNAME,
			PasswordHash: hash,
			Email:        "admin2@example.org",
		})

		apiReq := ApiLogin{
			Username: USERNAME,
			Password: "wrong-password",
		}

		data, err := json.Marshal(apiReq)
		if assert.NoError(err) {
			w := httptest.NewRecorder()

			os.Setenv("SESSION_KEY", "not-a-secret-key")
			am := GetAuthMiddleware()

			r := gin.Default()
			r.POST("/user/login", am.LoginHandler)

			req, _ := http.NewRequest("POST", "/user/login", bytes.NewReader(data))
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Content-Length", string(len(data)))
			r.ServeHTTP(w, req)

			assert.Equal(http.StatusUnauthorized, w.Code)
		}
	}
}

func TestAuthMiddleware_MiddlewareFunc(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()

	os.Setenv("SESSION_KEY", "not-a-secret-key")
	am := GetAuthMiddleware()

	r := gin.Default()
	api := r.Group("/api/")
	api.Use(am.MiddlewareFunc())
	{
		api.GET("/user/*username", user.GetUserHandler)
	}

	req, _ := http.NewRequest("GET", "/api/user/non-existing", nil)
	r.ServeHTTP(w, req)

	assert.Equal(http.StatusUnauthorized, w.Code)
}

// create user, login, and get user info
func TestAuthMiddleware_MiddlewareFunc2(t *testing.T) {
	assert := assert.New(t)

	const (
		USERNAME = "testuser"
		PASSWORD = "test1234"
		EMAIL = "testuser@example.org"
	)

	hash, err := helper.HashPassword(PASSWORD)
	if assert.NoError(err) {
		err := store.Users().Add(common.User{
			Username:     USERNAME,
			PasswordHash: hash,
			Email:        EMAIL,
		})

		if assert.NoError(err) {
			apiReq := ApiLogin{
				Username: USERNAME,
				Password: PASSWORD,
			}

			data, err := json.Marshal(apiReq)
			if assert.NoError(err) {
				w := httptest.NewRecorder()

				os.Setenv("SESSION_KEY", "not-a-secret-key")
				am := GetAuthMiddleware()

				r := gin.Default()
				r.POST("/user/login", am.LoginHandler)

				req, _ := http.NewRequest("POST", "/user/login", bytes.NewReader(data))
				req.Header.Add("Content-Type", "application/json")
				req.Header.Add("Content-Length", string(len(data)))
				r.ServeHTTP(w, req)

				if assert.Equal(http.StatusOK, w.Code) {
					body, err := ioutil.ReadAll(w.Body)

					if assert.NoError(err) {
						var resp map[string]string

						err = json.Unmarshal(body, &resp)
						if assert.NoError(err) {
							w = httptest.NewRecorder()
							r = gin.Default()

							api := r.Group("/api/")
							api.Use(am.MiddlewareFunc())
							{
								api.GET("/user/*username", user.GetUserHandler)
							}

							req, _ := http.NewRequest("GET", "/api/user/" + USERNAME, nil)
							req.Header.Add("Authorization", "Bearer "+resp["token"])
							r.ServeHTTP(w, req)

							assert.Equal(http.StatusOK, w.Code)
						}
					}
				}
			}
		}
	}
}


func TestAuthMiddleware_MiddlewareFunc3(t *testing.T) {
	assert := assert.New(t)

	os.Setenv("SESSION_KEY", "not-a-secret-key")
	am := GetAuthMiddleware()

	w := httptest.NewRecorder()
	r := gin.Default()

	api := r.Group("/api/")
	api.Use(am.MiddlewareFunc())
	{
		api.GET("/user/*username", user.GetUserHandler)
	}

	req, _ := http.NewRequest("GET", "/api/user/admin", nil)
	req.Header.Add("Authorization", "Bearer vkldklWEfweklergkl3KSDFJJWHEGKLSDFKSKDFJAHWEH")
	r.ServeHTTP(w, req)

	assert.Equal(http.StatusUnauthorized, w.Code)
}
