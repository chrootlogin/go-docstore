package store

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDocuments(t *testing.T) {
	assert := assert.New(t)

	docs := Documents()

	assert.NotNil(docs)
}

func TestDocuments_Create(t *testing.T) {
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
		err := Documents().Create(testCase.Path, testCase.Content)

		if testCase.Correct {
			assert.NoError(err)
		} else {
			assert.Error(err)
		}
	}
}
