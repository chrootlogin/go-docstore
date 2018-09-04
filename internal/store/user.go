package store

import (
	"errors"
	"time"

	"github.com/asdine/storm"
	"github.com/patrickmn/go-cache"

	"github.com/chrootlogin/go-docstore/internal/common"
	"github.com/chrootlogin/go-docstore/internal/database"
	"fmt"
)

var (
	userCache       *cache.Cache
	ErrUserNotExist = errors.New("user does not exist")
)

func init() {
	userCache = cache.New(30*time.Minute, 10*time.Minute)
}

type userList struct{}

func (ul *userList) Get(username string) (*common.User, error) {
	u, found := userCache.Get(username)
	if found {
		fmt.Println(fmt.Sprintf("found %v", u))
		return u.(*common.User), nil
	}

	var user common.User
	err := database.DB().One("Username", username, &user)
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
	err := database.DB().Save(&user)
	if err != nil {
		return err
	}

	// add user to cache
	userCache.Set(user.Username, &user, cache.DefaultExpiration)
	return nil
}

func Users() *userList {
	return &userList{}
}
