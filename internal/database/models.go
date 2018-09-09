package database

import (
	"crypto/sha256"

	"github.com/google/uuid"

	"github.com/chrootlogin/go-docstore/pkg/docstore"
	"github.com/chrootlogin/go-docstore/internal/common"
)

type IDocumentsDB interface {
	Create(string, []byte) (uuid.UUID, error)
	Read(string) (*docstore.Document, error)
	ReadFile([sha256.Size]byte) ([]byte, error)
}

type IUsersDB interface {
	Add(common.User) (*uuid.UUID, error)
	Get(string) (*common.User, error)
}
