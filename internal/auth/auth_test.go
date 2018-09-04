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

/*
func TestAuthMiddleware_MiddlewareFunc2(t *testing.T) {
	assert := assert.New(t)

	apiReq := ApiLogin{
		Username: "admin",
		Password: "admin1234",
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
			var resp map[string]string

			body, err := ioutil.ReadAll(w.Body)
			if assert.NoError(err) {
				err = json.Unmarshal(body, &resp)
				if assert.NoError(err) {
					w = httptest.NewRecorder()
					r = gin.Default()

					api := r.Group("/api/")
					api.Use(am.MiddlewareFunc())
					{
						//api.POST("/page/*path", page.PostPageHandler)
					}

					apiReq := map[string]string{
						"content": "# Test content",
					}

					data, err := json.Marshal(apiReq)
					if assert.NoError(err) {
						req, _ := http.NewRequest("POST", "/api/page/new-page.md", bytes.NewReader(data))
						req.Header.Add("Content-Type", "application/json")
						req.Header.Add("Content-Length", string(len(data)))
						req.Header.Add("Authorization", "Bearer "+resp["token"])
						r.ServeHTTP(w, req)

						assert.Equal(http.StatusCreated, w.Code)
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
		//api.GET("/page/*path", page.PostPageHandler)
	}

	req, _ := http.NewRequest("GET", "/api/page/index.md", nil)
	req.Header.Add("Authorization", "Bearer vkldklWEfweklergkl3KSDFJJWHEGKLSDFKSKDFJAHWEH")
	r.ServeHTTP(w, req)

	assert.Equal(http.StatusUnauthorized, w.Code)
}*/
