package database

import (
	"github.com/asdine/storm"
	"github.com/google/uuid"

	"github.com/chrootlogin/go-docstore/internal/common"
)

type UsersDB struct {
	r storm.Node
}

func (u *UsersDB) Get(username string) (*common.User, error) {
	var user common.User
	err := u.r.One("Username", username, &user)
	if err != nil {
		if err == storm.ErrNotFound {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (u *UsersDB) Add(user common.User) (*uuid.UUID, error) {
	// create UUID for user
	uid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	user.ID = uid

	// add user to database
	err = u.r.Save(&user)
	if err != nil {
		return nil, err
	}

	return &uid, nil
}
