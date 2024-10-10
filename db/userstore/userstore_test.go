package userstore

import (
	"context"
	"testing"

	"github.com/benjaminkitson/bk-user-api/internal/testhelpers"
	"github.com/benjaminkitson/bk-user-api/models"
	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NewStore(t *testing.T) UserStore {
	th := testhelpers.DBTester{}
	testTableName := "user"
	tableName := th.CreateLocalTable(t, testTableName)
	client := th.GetTestClient()
	t.Cleanup(func() { th.DeleteLocalTable(t, tableName) })
	return NewUserStore(client, testTableName)
}

// TODO: generally tidy these up a bit, make them a bit more robust
func TestGetUserByID(t *testing.T) {
	ctx := context.Background()
	store := NewStore(t)
	r, err := store.GetByID(ctx, "12345")
	require.NoError(t, err)

	if r.Email != "benk13@gmail.com" {
		t.Fatalf("Expected email to be benk13@gmail.com, got %v", r.Email)
	}
}

func TestGetUserByEmail(t *testing.T) {
	ctx := context.Background()
	store := NewStore(t)
	expectedEmail := "benk13@gmail.com"
	r, err := store.GetByEmail(ctx, expectedEmail)
	require.NoError(t, err)

	if r.Email != expectedEmail {
		t.Fatalf("Expected email to be benk13@gmail.com, got %v", r.Email)
	}

	_, err = store.GetByEmail(ctx, "nonexistent@gmail.com")

}

func TestPutUser(t *testing.T) {
	ctx := context.Background()
	store := NewStore(t)
	id := uuid.New().String()
	email := "someother@gmail.com"

	u, err := store.Put(ctx, models.User{Email: email}, id)
	require.NoError(t, err)

	r, err := store.GetByID(ctx, id)
	require.NoError(t, err)

	assert.Equal(t, email, u.Email)

	if r.Email != email {
		t.Fatalf("Expected username to be user2, got %v", r.Email)
	}
}

func TestDeleteUser(t *testing.T) {
	ctx := context.Background()
	store := NewStore(t)
	id := uuid.New().String()
	email := "someother@gmail.com"

	_, err := store.Put(ctx, models.User{Email: email}, id)
	require.NoError(t, err)

	s, err := store.Delete(ctx, id)
	require.NoError(t, err)

	assert.Equal(t, s, id)
}
