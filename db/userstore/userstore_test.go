package userstore

import (
	"context"
	"testing"

	"github.com/benjaminkitson/bk-user-api/internal/testhelpers"

	"github.com/stretchr/testify/require"
)

func NewStore(t *testing.T) UserStore {
	th := testhelpers.DBTester{}
	tableName := th.CreateLocalTable(t)
	client := th.GetTestClient()
	t.Cleanup(func() { th.DeleteLocalTable(t, tableName) })
	return NewUserStore(client)
}

func TestGetUser(t *testing.T) {
	ctx := context.Background()
	store := NewStore(t)
	r, err := store.Get(ctx, "12345")
	require.NoError(t, err)

	if r.Username != "benk13" {
		t.Fatalf("Expected username to be benk13, got %v", r.Username)
	}
}

func TestPutUser(t *testing.T) {
	ctx := context.Background()
	store := NewStore(t)
	err := store.Put(ctx, User{Username: "user2"}, "23456")
	require.NoError(t, err)

	r, err := store.Get(ctx, "23456")
	require.NoError(t, err)

	if r.Username != "user2" {
		t.Fatalf("Expected username to be user2, got %v", r.Username)
	}
}
