package doc

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/json"
	"github.com/stretchr/testify/assert"

	"github.com/chrootlogin/go-docstore/internal/database"
)

func TestCreateDocumentHandler(t *testing.T) {
	assert := assert.New(t)

	type TestCases struct {
		Name    string
		Path    string
		Content []byte
		Status  int
	}

	testCases := []TestCases{
		{
			Name:    "index.md",
			Path:    "/index.md",
			Content: []byte("# Hello\nThis is an index file"),
			Status:  http.StatusCreated,
		},
		{
			Name:    "mammut.md",
			Path:    "/docs/mammut.md",
			Content: []byte("### This is a mammut"),
			Status:  http.StatusCreated,
		},
		{
			Path:    "/docs/",
			Content: []byte("Test 1234"),
			Status:  http.StatusInternalServerError,
		},
		{
			Path:    "/",
			Content: []byte("Test 1234"),
			Status:  http.StatusBadRequest,
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.Path)

		w := httptest.NewRecorder()

		r := gin.Default()
		r.POST("/doc/*path", CreateDocumentHandler)

		doc := ApiDocument{
			Content: base64.StdEncoding.EncodeToString(testCase.Content),
		}

		data, err := json.Marshal(doc)
		if assert.NoError(err) {
			req, _ := http.NewRequest("POST", "/doc"+testCase.Path, bytes.NewReader(data))
			r.ServeHTTP(w, req)
			if assert.Equal(testCase.Status, w.Code) {
				if testCase.Status == http.StatusCreated {
					document, err := database.DB().Documents().Read(testCase.Path)
					if assert.NoError(err) {
						assert.Equal(testCase.Name, document.Name)
					}
				}
			}
		}
	}
}

func TestCreateDocumentHandler2(t *testing.T) {
	assert := assert.New(t)
	w := httptest.NewRecorder()

	r := gin.Default()
	r.POST("/doc/*path", CreateDocumentHandler)

	req, _ := http.NewRequest("POST", "/doc/empty.md", bytes.NewReader([]byte{}))
	r.ServeHTTP(w, req)

	assert.Equal(http.StatusBadRequest, w.Code)
}

func TestCreateDocumentHandler3(t *testing.T) {
	assert := assert.New(t)
	w := httptest.NewRecorder()

	r := gin.Default()
	r.POST("/doc/*path", CreateDocumentHandler)

	doc := ApiDocument{
		Content: "z,y",
	}

	data, err := json.Marshal(doc)
	if assert.NoError(err) {
		req, _ := http.NewRequest("POST", "/doc/wrong.md", bytes.NewReader(data))
		r.ServeHTTP(w, req)

		assert.Equal(http.StatusInternalServerError, w.Code)
	}
}
