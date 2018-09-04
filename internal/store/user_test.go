package store

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/chrootlogin/go-docstore/internal/common"
)

func TestUsers(t *testing.T) {
	assert := assert.New(t)

	ul := Users()

	assert.NotNil(ul)
}

func TestUserList_Get(t *testing.T) {
	assert := assert.New(t)

	_, err := Users().Get("no-exist")
	if assert.Error(err) {
		assert.Equal(ErrUserNotExist, err)
	}
}

func TestUserList_Add(t *testing.T) {
	assert := assert.New(t)

	newUser := common.User{
		Username: "test-user",
		Email:    "test@example.org",
	}

	err := Users().Add(newUser)
	if assert.NoError(err) {
		user, err := Users().Get(newUser.Username)
		if assert.NoError(err) {
			assert.Equal(newUser.Username, user.Username)
			assert.Equal(newUser.Email, user.Email)
		}
	}
}

func TestUserList_Delete(t *testing.T) {
	assert := assert.New(t)

	newUser := common.User{
		Username: "test-user2",
		Email:    "test2@example.org",
	}

	// add user
	err := Users().Add(newUser)
	if assert.NoError(err) {
		// delete user
		err := Users().Delete(newUser.Username)

		t.Log(userCache.Items())
		// check if exists
		if assert.NoError(err) {
			t.Log("GET")

			u, err := Users().Get(newUser.Username)
			t.Log(u)
			if assert.Error(err) {
				assert.Equal(ErrUserNotExist, err)
			}
		}
	}
}

// delete non existing user
func TestUserList_Delete2(t *testing.T) {
	assert := assert.New(t)

	err := Users().Delete("no-exist")

	if assert.Error(err) {
		assert.Equal(ErrUserNotExist, err)
	}
}
