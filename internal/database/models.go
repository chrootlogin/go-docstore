package database

import (
	"github.com/chrootlogin/go-docstore/pkg/docstore"
	"github.com/google/uuid"
)

type IDocumentsDB interface {
	Create(string, []byte) (uuid.UUID, error)
	Read(string) (*docstore.Document, error)
}
