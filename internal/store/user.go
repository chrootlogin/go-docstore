package store

import (
	"time"

	"github.com/asdine/storm"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"

	"github.com/chrootlogin/go-docstore/internal/common"
	"github.com/chrootlogin/go-docstore/internal/database"
)

var (
	userCache *cache.Cache
)

func init() {
	userCache = cache.New(30*time.Minute, 10*time.Minute)
}

type userList struct{}

func (ul *userList) Get(username string) (*common.User, error) {
	u, found := userCache.Get(username)
	if found {
		return u.(*common.User), nil
	}

	var user common.User
	err := database.DB().Users().One("Username", username, &user)
	if err != nil {
		if err == storm.ErrNotFound {
			return nil, ErrUserNotExist
		}

		return nil, err
	}

	userCache.Set(username, &user, cache.DefaultExpiration)

	return &user, nil
}

func (ul *userList) Add(user common.User) error {
	// create UUID for user
	uid, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	user.ID = uid

	// add user to database
	err = database.DB().Users().Save(&user)
	if err != nil {
		return err
	}

	// add user to cache
	userCache.Set(user.Username, &user, cache.DefaultExpiration)

	return nil
}

func (ul *userList) Delete(username string) error {
	user, err := ul.Get(username)
	if err != nil {
		return err
	}

	// remove from cache
	userCache.Delete(user.Username)

	// delete from database
	return database.DB().Users().DeleteStruct(user)
}

func Users() *userList {
	return &userList{}
}
