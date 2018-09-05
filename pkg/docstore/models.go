package docstore

import (
	"crypto/sha256"
	"time"

	"github.com/google/uuid"
)

type Document struct {
	ID        uuid.UUID `storm:"id"`
	Name      string
	Revisions []FileRevision
}

type FileRevision struct {
	ID       uuid.UUID `storm:"id"`
	Time     time.Time
	FileHash [sha256.Size]byte
}

type File struct {
	Hash    [sha256.Size]byte `storm:"id"`
	Content []byte
}
