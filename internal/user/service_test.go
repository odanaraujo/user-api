package user

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/odanaraujo/user-api/internal/model"
	"github.com/stretchr/testify/assert"
)

type mockCache struct {
	store map[string][]byte
}

func newMockCache() *mockCache {
	return &mockCache{store: make(map[string][]byte)}
}

func (m *mockCache) Get(key string) ([]byte, bool) {
	val, ok := m.store[key]
	return val, ok
}

func (m *mockCache) Set(key string, value []byte, ttl time.Duration) {
	m.store[key] = value
}

func (m *mockCache) Delete(key string) {
	delete(m.store, key)
}

func TestCreateUser(t *testing.T) {
	mock := newMockCache()
	service := NewUserService(mock)

	user := &model.User{
		ID:    "123",
		Name:  "João",
		CPF:   "12345678900",
		Age:   30,
		Email: "joao@gmail.com",
	}

	createdUser, err := service.CreateUser(user)
	assert.Nil(t, err)
	assert.Equal(t, user, createdUser)

	// verifica se o dado foi realmente salvo no cache
	data, ok := mock.Get("123")
	assert.True(t, ok)

	var cached model.User
	_ = json.Unmarshal(data, &cached)
	assert.Equal(t, user.ID, cached.ID)
	assert.Equal(t, user.Email, cached.Email)
}

func TestCreateUserWithEmailError(t *testing.T) {
	mock := newMockCache()
	service := NewUserService(mock)

	user := &model.User{
		ID:    "123",
		Name:  "João",
		CPF:   "12345678900",
		Age:   30,
		Email: "joao@.com",
	}

	_, err := service.CreateUser(user)
	assert.NotEmpty(t, err)
}

func TestGetUserByID(t *testing.T) {
	mock := newMockCache()
	service := NewUserService(mock)

	user := model.User{
		ID:    "abc",
		Name:  "Maria",
		CPF:   "98765432100",
		Age:   25,
		Email: "maria@email.com",
	}

	data, _ := json.Marshal(user)
	mock.Set(user.ID, data, time.Hour)

	got, err := service.GetUserByID(user.ID)
	assert.Nil(t, err)
	assert.Equal(t, user.Name, got.Name)
	assert.Equal(t, user.Email, got.Email)
}

func TestDeleteUser(t *testing.T) {
	mock := newMockCache()
	service := NewUserService(mock)

	user := model.User{ID: "999"}
	data, _ := json.Marshal(user)
	mock.Set(user.ID, data, time.Hour)

	err := service.DeleteUser(user.ID)
	assert.Nil(t, err)

	_, ok := mock.Get(user.ID)
	assert.False(t, ok)
}
