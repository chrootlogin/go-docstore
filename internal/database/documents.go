package database

import (
	"crypto/sha256"
	"fmt"
	"path/filepath"
	"time"

	"github.com/asdine/storm"
	"github.com/google/uuid"

	"github.com/chrootlogin/go-docstore/pkg/docstore"
	"sort"
)

const (
	ROOT_NODE = "root"
	DATA_NODE = "data"
)

type DocumentsDB struct {
	r storm.Node
}

func (d *DocumentsDB) Create(path string, content []byte) (uuid.UUID, error) {
	dir, file := filepath.Split(path)
	if len(file) == 0 {
		return uuid.UUID{}, ErrNoFilename
	}

	filehash, err := d.saveFile(content)
	if err != nil {
		return uuid.UUID{}, err
	}

	docUuid, err := uuid.NewUUID()
	if err != nil {
		return uuid.UUID{}, err
	}

	revUuid, err := uuid.NewUUID()
	if err != nil {
		return uuid.UUID{}, err
	}

	doc := docstore.Document{
		ID:   docUuid,
		Name: file,
		Revisions: []docstore.FileRevision{
			{
				ID:       revUuid,
				Time:     time.Now(),
				FileHash: filehash,
			},
		},
	}

	// sort file revisions
	sort.Slice(doc.Revisions, func(i, j int) bool {
		return doc.Revisions[i].Time.Before(doc.Revisions[j].Time)
	})

	node := travelPath(dir, d.r.From(ROOT_NODE))

	var document docstore.Document
	fmt.Println("Doc")
	err = node.One("Name", file, &document)
	if err != nil {
		if err == storm.ErrNotFound {
			// save document
			node.Save(&doc)

			return docUuid, nil
		}

		return uuid.UUID{}, err
	}

	return uuid.UUID{}, ErrExists
}

func (d *DocumentsDB) Read(path string) (*docstore.Document, error) {
	dir, file := filepath.Split(path)
	if len(file) == 0 {
		return nil, ErrNoFilename
	}

	node := travelPath(dir, d.r.From(ROOT_NODE))

	var document docstore.Document
	err := node.One("Name", file, &document)
	if err != nil {
		if err == storm.ErrNotFound {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &document, nil
}

func (d *DocumentsDB) ReadFile(hash [sha256.Size]byte) ([]byte, error) {
	node := d.r.From(DATA_NODE)

	var file docstore.File
	err := node.One("Hash", hash, &file)
	if err != nil {
		if err == storm.ErrNotFound {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return file.Content, nil
}

func (d *DocumentsDB) saveFile(content []byte) ([sha256.Size]byte, error) {
	hash := sha256.Sum256(content)

	node := d.r.From(DATA_NODE)

	var file docstore.File

	fmt.Println("File")
	err := node.One("Hash", hash, &file)
	if err != nil {
		// save
		if err == storm.ErrNotFound {
			file = docstore.File{
				Hash:    hash,
				Content: content,
			}

			node.Save(&file)

			return hash, nil
		}

		return hash, err
	}

	return hash, nil
}

func travelPath(dir string, node storm.Node) storm.Node {
	// travel trough subdirectories
	dirs := filepath.SplitList(dir)
	for _, dirname := range dirs {
		node = node.From(dirname)
	}

	return node
}
