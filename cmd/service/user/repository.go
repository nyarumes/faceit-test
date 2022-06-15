package user

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/google/uuid"
)

const (
	defaultUserBucketName = "users"
)

var (
	ErrUserNotFound = errors.New("user isn't found")
)

type Repository struct {
	storage *bolt.DB
}

// NewRepository inits repository based bolt db
func NewRepository(db *bolt.DB) (*Repository, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(defaultUserBucketName))

		return err
	})

	return &Repository{
		storage: db,
	}, err
}

// Insert adds a new record with user
func (r Repository) Insert(user Model) error {
	return r.storage.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(defaultUserBucketName))

		buf, err := json.Marshal(user)
		if err != nil {
			return err
		}

		return b.Put([]byte(user.ID.String()), buf)
	})
}

// Update updates existing user
func (r Repository) Update(updatedUser Model) error {
	return r.storage.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(defaultUserBucketName))

		data := b.Get([]byte(updatedUser.ID.String()))
		if len(data) == 0 {
			return ErrUserNotFound
		}

		var user Model
		err := json.Unmarshal(data, &user)
		if err != nil {
			return err
		}

		user.FirstName = updatedUser.FirstName
		user.LastName = updatedUser.LastName
		user.Nickname = updatedUser.Nickname
		user.Country = updatedUser.Country
		user.UpdatedAt = time.Now()

		buf, err := json.Marshal(user)
		if err != nil {
			return err
		}

		return b.Put([]byte(user.ID.String()), buf)
	})
}

// Get gets a user by id
func (r Repository) Get(id uuid.UUID) (user Model, err error) {
	err = r.storage.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(defaultUserBucketName))

		data := b.Get([]byte(id.String()))
		if len(data) == 0 {
			return ErrUserNotFound
		}

		return json.Unmarshal(data, &user)
	})

	return
}

func (r Repository) Delete(id uuid.UUID) error {
	return r.storage.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(defaultUserBucketName))

		return b.Delete([]byte(id.String()))
	})
}

type userFilter func(user Model) bool

func WithUserNickNameFilter(nickname string) userFilter {
	return func(user Model) bool {
		return strings.Contains(user.Nickname, nickname)
	}
}

// FindPaginatedWithFilter finds scope of users with pagination and filters for fields
func (r Repository) FindPaginatedWithFilter(offset int, count int, filters ...userFilter) (users []Model, err error) {
	err = r.storage.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(defaultUserBucketName))

		c := b.Cursor()
		i := 0

		for k, v := c.First(); k != nil; k, v = c.Next() {
			i++

			if offset > i {
				continue
			}

			if i > count {
				break
			}

			var user Model

			err := json.Unmarshal(v, &user)
			if err != nil {
				return err
			}

			ok := true
			for _, filter := range filters {
				ok = filter(user)
			}

			if !ok {
				continue
			}

			users = append(users, user)
		}

		return nil
	})

	return
}
