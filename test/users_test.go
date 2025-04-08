package test

import (
	"snippetbox/internal/assert"
	"snippetbox/internal/modules/user/repository"
	"testing"
)

func TestUserExist(t *testing.T) {
	db, err := newTestDB(t)
	if err != nil {
		t.Fatal(err)
	}
	users := &repository.UserRepository{DB: db}
	tests := []struct {
		name   string
		userID int
		want   bool
	}{
		{
			name:   "Valid ID",
			userID: 1,
			want:   true,
		},
		{
			name:   "Zero ID",
			userID: 0,
			want:   false,
		},
		{
			name:   "Non-existent ID",
			userID: 2,
			want:   false,
		},
	}
	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			ok, err := users.Exists(test.userID)

			assert.Equal(t, test.want, ok)
			assert.NilError(t, err)
		})
	}
}
