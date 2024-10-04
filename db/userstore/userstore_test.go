package userstore

import (
	"context"
	"fmt"
	"testing"

	"github.com/benjaminkitson/bk-user-api/internal/testhelpers"
	"github.com/google/uuid"

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

func TestGetUser(t *testing.T) {
	ctx := context.Background()
	store := NewStore(t)
	r, err := store.Get(ctx, "12345")
	require.NoError(t, err)

	if r.Email != "benk13@gmail.com" {
		t.Fatalf("Expected username to be benk13@gmail.com, got %v", r.Email)
	}
}

func TestPutUser(t *testing.T) {
	ctx := context.Background()
	store := NewStore(t)
	id := uuid.New().String()

	fmt.Println(id)

	err := store.Put(ctx, User{Email: "someother@gmail.com"}, id)
	require.NoError(t, err)

	r, err := store.Get(ctx, id)
	require.NoError(t, err)

	if r.Email != "user2" {
		t.Fatalf("Expected username to be user2, got %v", r.Email)
	}
}
