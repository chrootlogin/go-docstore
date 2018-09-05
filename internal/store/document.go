package store

import (
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/chrootlogin/go-docstore/internal/database"
)

var (
	documentCache *cache.Cache
)

func init() {
	documentCache = cache.New(30*time.Minute, 10*time.Minute)
}

type documents struct{}

func (d *documents) Create(path string, content []byte) error {
	database.DB().Documents().Create(path, content)

	return nil
}

func Documents() *documents {
	return &documents{}
}
