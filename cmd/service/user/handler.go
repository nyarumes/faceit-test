package user

import (
	"time"

	"github.com/google/uuid"
)

const (
	// FilterUserNickname represents a filter name that filters nickname by value
	FilterUserNickname = "nickname"
)

// UniqueRepository represents any storage fit user capabilities
type UniqueRepository interface {
	Insert(user Model) error
	Update(user Model) error
	Delete(id uuid.UUID) error
	FindPaginatedWithFilter(offset int, count int, filters ...userFilter) (users []Model, err error)
}

type UniqueHandler interface {
	Add(firstName, lastName, nickName, password, email, country string) error
	Update(id uuid.UUID, firstName, lastName, nickName, country string) error
	Remove(id uuid.UUID) error
	FindAllWithPaginationAndFilter(offset, count int, filters map[string]string) ([]Model, error)
}

// Handler represents user handler
type Handler struct {
	repository UniqueRepository
}

// NewHandler inits a new handler that unique (rpc / http)
func NewHandler(repository UniqueRepository) *Handler {
	return &Handler{
		repository: repository,
	}
}

// Add adds a new user to storage
func (h *Handler) Add(firstName, lastName, nickName, password, email, country string) error {
	id := uuid.New()

	return h.repository.Insert(Model{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Nickname:  nickName,
		Password:  password,
		Email:     email,
		Country:   country,
		CreatedAt: time.Now(),
	})
}

// Update updates user with allowed fields
func (h Handler) Update(id uuid.UUID, firstName, lastName, nickName, country string) error {
	return h.repository.Update(Model{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Nickname:  nickName,
		Country:   country,
	})
}

// Remove deletes user
func (h Handler) Remove(id uuid.UUID) error {
	return h.repository.Delete(id)
}

// FindAllWithPaginationAndFilter finds all users paginated and applied filters by schema
func (h Handler) FindAllWithPaginationAndFilter(offset, count int, filters map[string]string) ([]Model, error) {
	var userFilters []userFilter
	for name, value := range filters {
		switch name {
		case FilterUserNickname:
			userFilters = append(userFilters, WithUserNickNameFilter(value))
		}
	}

	return h.repository.FindPaginatedWithFilter(offset, count, userFilters...)
}
