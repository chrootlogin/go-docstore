package doc

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/chrootlogin/go-docstore/internal/common"
	"github.com/chrootlogin/go-docstore/internal/database"
)

type ApiDocument struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func CreateDocumentHandler(c *gin.Context) {
	path := c.Param("path")
	if len(path) <= 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.ApiResponse{Message: common.WrongAPIUsageError})
		return
	}

	var doc ApiDocument
	if c.BindJSON(&doc) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.ApiResponse{Message: common.WrongAPIUsageError})
		return
	}

	content, err := base64.StdEncoding.DecodeString(doc.Content)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, common.ApiResponse{Message: "error decoding content"})
		return
	}

	docUuid, err := database.DB().Documents().Create(path, content)
	if err != nil {
		log.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, common.ApiResponse{Message: "error saving document"})
		return
	}

	c.JSON(http.StatusCreated, common.ApiResponse{Message: fmt.Sprintf("created document: %s", docUuid)})
}

func ReadDocumentHandler(c *gin.Context) {
	path := c.Param("path")
	if len(path) <= 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.ApiResponse{Message: common.WrongAPIUsageError})
		return
	}

	doc, err := database.DB().Documents().Read(path)
	if err != nil {
		if err == database.ErrNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, common.ApiResponse{Message: "not found"})
			return
		}

		log.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, common.ApiResponse{Message: "error saving document"})
		return
	}

	// load last revision
	content, err := database.DB().Documents().ReadFile(doc.Revisions[0].FileHash)
	if err != nil {
		log.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, common.ApiResponse{Message: "error saving document"})
		return
	}

	data := base64.StdEncoding.EncodeToString(content)
	c.JSON(http.StatusOK, ApiDocument{
		Name:    doc.Name,
		Content: data,
	})
}
