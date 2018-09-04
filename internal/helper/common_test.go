package helper

import (
	"testing"
)

func TestCheckPasswordHash(t *testing.T) {
	type TestMatrix struct {
		ClearTextPassword string
		Hash              string
		Correct           bool
	}

	testMatrix := []TestMatrix{
		{
			ClearTextPassword: ")sdDfjl=+BBb091!",
			Hash:              "$2a$14$F.jfRxIxMysxJA2nQv4zhuwq97hfdNBoKsRis0wy1edesof48o6sO",
			Correct:           true,
		},
		{
			ClearTextPassword: "oaskfRWE%çç)df093!",
			Hash:              "$2a$14$lELHxPz4dMPcMfFg7HVY8OyWjIuLBz/2dDeoQLL6CRiymcxpT/8um",
			Correct:           true,
		},
		{
			ClearTextPassword: "!ldsfSDGJwkelg$)*",
			Hash:              "$2a$14$0pxZWwT6Y2lvK4iMZqduEeWGDOngsRwVDVoXGKfliz6YTSmhY1bpi",
			Correct:           false,
		},
	}

	for _, test := range testMatrix {
		result := CheckPasswordHash(test.ClearTextPassword, test.Hash)

		if result != test.Correct {
			t.Errorf("Password check on %s was not expected result", test.ClearTextPassword)
		}
	}
}
