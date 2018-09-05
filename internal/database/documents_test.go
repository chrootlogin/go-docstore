package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDb_Documents(t *testing.T) {
	assert := assert.New(t)

	docs := DB().Documents()

	assert.NotNil(docs)
}

func TestDb_Users(t *testing.T) {
	assert := assert.New(t)

	docs := DB().Users()

	assert.NotNil(docs)
}

func TestDocumentsDB_Create(t *testing.T) {
	assert := assert.New(t)

	type TestCases struct {
		Path    string
		Content []byte
		Correct bool
	}

	testCases := []TestCases{
		{
			Path:    "/index.md",
			Content: []byte("# Hello\nThis is an index file"),
			Correct: true,
		},
		{
			Path:    "/docs/mammut.md",
			Content: []byte("### This is a mammut"),
			Correct: true,
		},
		{
			Path:    "/docs/",
			Content: []byte("Test 1234"),
			Correct: false,
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.Path)
		uuid, err := DB().Documents().Create(testCase.Path, testCase.Content)

		t.Log(uuid)
		if testCase.Correct {
			assert.NoError(err)
		} else {
			assert.Error(err)
		}
	}
}

func TestDocumentsDB_Read(t *testing.T) {
	assert := assert.New(t)

	type TestCases struct {
		Path    string
		Name    string
		Content []byte
		Correct bool
	}

	testCases := []TestCases{
		{
			Path:    "/index.md",
			Name:    "index.md",
			Content: []byte("# Hello\nThis is an index file"),
			Correct: true,
		},
		{
			Path:    "/docs/mammut.md",
			Name:    "mammut.md",
			Content: []byte("### This is a mammut"),
			Correct: true,
		},
		{
			Path:    "/docs/",
			Content: []byte("Test 1234"),
			Correct: false,
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.Path)

		doc, err := DB().Documents().Read(testCase.Path)
		if testCase.Correct {
			if assert.NoError(err) {
				assert.Equal(testCase.Name, doc.Name)
			}
		} else {
			if assert.Error(err) {
				assert.Equal(ErrNoFilename, err)
			}
		}

	}
}
