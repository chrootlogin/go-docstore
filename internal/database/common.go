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
	"github.com/chrootlogin/go-docstore/pkg/docstore"
	"github.com/chrootlogin/go-docstore/internal/helper"
)

const (
	DEFAULT_USERNAME = "admin"
	DEFAULT_PASSWORD = "admin"
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
		newDB := false
		if _, err := os.Stat("dbPath"); os.IsNotExist(err) {
			newDB = true
		}

		d, err := storm.Open(dbPath, storm.Codec(json.Codec))
		if err != nil {
			log.Fatal(fmt.Sprintf("Can't open database: %s", err.Error()))
		}

		// Init objects
		err = d.Init(&common.User{})
		if err != nil {
			log.Fatal(fmt.Sprintf("Can't init database: %s", err.Error()))
		}

		err = d.Init(&docstore.Document{})
		if err != nil {
			log.Fatal(fmt.Sprintf("Can't init database: %s", err.Error()))
		}

		instance = db{
			db: d,
		}

		if newDB {
			instance.initDatabase()
		}
	})

	return &instance
}

func (d *db) Users() storm.Node {
	return d.db.From("users")
}

func (d *db) User() IUsersDB {
	return &UsersDB{
		d.db.From("users"),
	}
}

func (d *db) Documents() IDocumentsDB {
	return &DocumentsDB{
		d.db.From("documents"),
	}
}

func (d *db) initDatabase() {
	pwHash, err := helper.HashPassword(DEFAULT_PASSWORD)
	if err != nil {
		log.Fatal(fmt.Sprintf("Can't hash password: %s", err.Error()))
	}

	// creating default user
	u := common.User{
		Username: DEFAULT_USERNAME,
		PasswordHash: pwHash,
		Email: "admin@example.com",
		IsEnabled: true,
		Permissions: []string{
			"ADMIN",
		},
	}

	_, err = d.User().Add(u)
	if err != nil {
		log.Fatal(fmt.Sprintf("Couldn't create default user: %s", err.Error()))
	}

	log.Println(fmt.Sprintf("Created default user '%s' with default password '%s'. Please change this password immediately", DEFAULT_USERNAME, DEFAULT_PASSWORD))
}