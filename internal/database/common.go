package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/json"

	"github.com/chrootlogin/go-docstore/internal/common"
)

var (
	dbPath   = ""
	instance *storm.DB
	once     sync.Once
)

func init() {
	dbPath = os.Getenv("DB_PATH")

	if len(dbPath) == 0 {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		dbPath = filepath.Join(dir, "data.db")
		log.Println("Environment variable DB_PATH is empty!")
	}
}

func DB() *storm.DB {
	once.Do(func() {
		db, err := storm.Open(dbPath, storm.Codec(json.Codec))
		if err != nil {
			log.Fatal(fmt.Sprintf("Can't open database: %s", err.Error()))
		}

		// Init objects
		err = db.Init(&common.User{})
		if err != nil {
			log.Fatal(fmt.Sprintf("Can't init database: %s", err.Error()))
		}

		instance = db
	})

	return instance
}

func Users() storm.Node {
	db := DB()

	return db.From("users")
}
