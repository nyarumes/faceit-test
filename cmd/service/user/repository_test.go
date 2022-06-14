package user

import (
	"os"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRepositoryUserInsert(t *testing.T) {
	dbName := uuid.New().String() + ".db"
	defer os.Remove(dbName)

	db, err := bolt.Open(dbName, 0600, nil)
	assert.NoError(t, err)

	defer db.Close()

	r, err := NewRepository(db)
	assert.NoError(t, err)

	var (
		user = Model{
			ID:       uuid.New(),
			Nickname: "test",
		}
	)

	err = r.Insert(user)
	assert.NoError(t, err)

	userUpdated, err := r.Get(user.ID)
	assert.NoError(t, err)

	assert.Equal(t, user.Nickname, userUpdated.Nickname)
}

func TestRepositoryUserUpdate(t *testing.T) {
	dbName := uuid.New().String() + ".db"
	defer os.Remove(dbName)

	db, err := bolt.Open(dbName, 0600, nil)
	assert.NoError(t, err)

	defer db.Close()

	r, err := NewRepository(db)
	assert.NoError(t, err)

	var (
		user = Model{
			ID:       uuid.New(),
			Nickname: "test",
		}
	)

	err = r.Insert(user)
	assert.NoError(t, err)

	userUpdated, err := r.Get(user.ID)
	assert.NoError(t, err)

	assert.Equal(t, user.Nickname, userUpdated.Nickname)

	user.Nickname = "test2"

	err = r.Update(user)
	assert.NoError(t, err)

	userUpdated, err = r.Get(user.ID)
	assert.NoError(t, err)

	assert.Equal(t, user.Nickname, userUpdated.Nickname)
}

func TestRepositoryFindPaginatedWithFilter(t *testing.T) {
	dbName := uuid.New().String() + ".db"
	defer os.Remove(dbName)

	db, err := bolt.Open(dbName, 0600, nil)
	assert.NoError(t, err)

	defer db.Close()

	r, err := NewRepository(db)
	assert.NoError(t, err)

	var (
		user = Model{
			ID:       uuid.New(),
			Nickname: "test",
		}
	)

	err = r.Insert(user)
	assert.NoError(t, err)

	users, err := r.FindPaginatedWithFilter(0, 10)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(users))

	users, err = r.FindPaginatedWithFilter(2, 10)
	assert.NoError(t, err)

	assert.Equal(t, 0, len(users))

	users, err = r.FindPaginatedWithFilter(0, 0)
	assert.NoError(t, err)

	assert.Equal(t, 0, len(users))
}
