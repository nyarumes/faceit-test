package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type MockHandler struct {
	users map[uuid.UUID]Model
}

func (m *MockHandler) Add(firstName, lastName, nickName, password, email, country string) error {
	id := uuid.New()

	m.users[id] = Model{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Nickname:  nickName,
		Password:  password,
		Email:     email,
		Country:   country,
	}

	return nil
}

func (m *MockHandler) Update(id uuid.UUID, firstName, lastName, nickName, country string) error {
	user, ok := m.users[id]
	if !ok {
		return ErrUserNotFound
	}

	user.FirstName = firstName
	user.LastName = lastName
	user.Nickname = nickName
	user.Country = country

	m.users[id] = user

	return nil
}

func (m *MockHandler) Remove(id uuid.UUID) error {
	delete(m.users, id)

	return nil
}

func (m *MockHandler) FindAllWithPaginationAndFilter(offset, count int, filters map[string]string) ([]Model, error) {
	var users []Model
	for _, user := range m.users {
		users = append(users, user)
	}

	return users, nil
}

func executeRequest(router *mux.Router, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

func TestHttp(t *testing.T) {
	r := HttpRouter(&MockHandler{users: make(map[uuid.UUID]Model)})

	request := Request{
		Nickname: "test",
	}

	data, err := json.Marshal(request)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/user", bytes.NewReader(data))
	assert.NoError(t, err)

	resp := executeRequest(r, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	req, err = http.NewRequest("GET", "/user", nil)
	assert.NoError(t, err)

	resp = executeRequest(r, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var users []Model
	err = json.NewDecoder(resp.Body).Decode(&users)
	assert.NoError(t, err)

	requestUpdate := Request{
		Nickname: "test2",
	}

	data, err = json.Marshal(requestUpdate)
	assert.NoError(t, err)

	req, err = http.NewRequest("PUT", "/user/"+users[0].ID.String(), bytes.NewReader(data))
	assert.NoError(t, err)

	resp = executeRequest(r, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	req, err = http.NewRequest("GET", "/user", nil)
	assert.NoError(t, err)

	resp = executeRequest(r, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	err = json.NewDecoder(resp.Body).Decode(&users)
	assert.NoError(t, err)

	assert.Equal(t, requestUpdate.Nickname, users[0].Nickname)

	req, err = http.NewRequest("DELETE", "/user/"+users[0].ID.String(), nil)
	assert.NoError(t, err)

	resp = executeRequest(r, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	req, err = http.NewRequest("GET", "/user", nil)
	assert.NoError(t, err)

	resp = executeRequest(r, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	err = json.NewDecoder(resp.Body).Decode(&users)
	assert.NoError(t, err)

	assert.Equal(t, 0, len(users))
}
