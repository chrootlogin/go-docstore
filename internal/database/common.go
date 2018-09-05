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

type db struct {
	db *storm.DB
}

var (
	dbPath   = ""
	instance db
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

func DB() *db {
	once.Do(func() {
		d, err := storm.Open(dbPath, storm.Codec(json.Codec))
		if err != nil {
			log.Fatal(fmt.Sprintf("Can't open database: %s", err.Error()))
		}

		// Init objects
		err = d.Init(&common.User{})
		if err != nil {
			log.Fatal(fmt.Sprintf("Can't init database: %s", err.Error()))
		}

		instance = db{
			db: d,
		}
	})

	return &instance
}

func (d *db) Users() storm.Node {
	return d.db.From("users")
}

func (d *db) Documents() IDocumentsDB {
	return &DocumentsDB{
		d.db.From("documents"),
	}
}
