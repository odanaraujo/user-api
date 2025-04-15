package user

import (
	"context"
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

func (m *mockCache) Get(ctx context.Context, key string) ([]byte, bool) {
	val, ok := m.store[key]
	return val, ok
}

func (m *mockCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) {
	m.store[key] = value
}

func (m *mockCache) Delete(ctx context.Context, key string) {
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

	createdUser, err := service.CreateUser(context.Background(), user)
	assert.Nil(t, err)
	assert.Equal(t, user, createdUser)

	// verifica se o dado foi realmente salvo no cache
	data, ok := mock.Get(context.Background(), "123")
	assert.True(t, ok)

	var cached model.User
	_ = json.Unmarshal(data, &cached)
	assert.Equal(t, user.ID, cached.ID)
	assert.Equal(t, user.Email, cached.Email)
}

func TestCreateUser_OverWrite(t *testing.T) {
	mock := newMockCache()
	service := NewUserService(mock)

	user1 := &model.User{ID: "1", Name: "first", CPF: "12345678910", Age: 30, Email: "1@a.com"}
	user2 := &model.User{ID: "1", Name: "seconds", CPF: "12345678911", Age: 25, Email: "2@a.com"}

	service.CreateUser(context.Background(), user1)
	service.CreateUser(context.Background(), user2)

	data, _ := mock.Get(context.Background(), "1")
	var cached model.User
	json.Unmarshal(data, &cached)

	assert.Equal(t, "seconds", cached.Name)
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

	_, err := service.CreateUser(context.Background(), user)
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
	mock.Set(context.Background(), user.ID, data, time.Hour)

	got, err := service.GetUserByID(context.Background(), user.ID)
	assert.Nil(t, err)
	assert.Equal(t, user.Name, got.Name)
	assert.Equal(t, user.Email, got.Email)
}

func TestGetUserByID_NotFound(t *testing.T) {
	mock := newMockCache()
	service := NewUserService(mock)

	got, err := service.GetUserByID(context.Background(), "não existe")
	assert.NotNil(t, err)
	assert.Nil(t, got)
}

func TestGetUserByID_InvalidJSON(t *testing.T) {
	mock := newMockCache()
	service := NewUserService(mock)

	mock.Set(context.Background(), "abc", []byte("this_not_json"), time.Hour)
	got, err := service.GetUserByID(context.Background(), "abc")
	assert.NotNil(t, err)
	assert.Nil(t, got)
	assert.Equal(t, "error performing unmarshal", err.Message)
}

func TestDeleteUser(t *testing.T) {
	mock := newMockCache()
	service := NewUserService(mock)

	user := model.User{ID: "999"}
	data, _ := json.Marshal(user)
	mock.Set(context.Background(), user.ID, data, time.Hour)

	err := service.DeleteUser(context.Background(), user.ID)
	assert.Nil(t, err)

	_, ok := mock.Get(context.Background(), user.ID)
	assert.False(t, ok)
}

func TestDelete_NotFound(t *testing.T) {
	mock := newMockCache()
	service := NewUserService(mock)
	err := service.DeleteUser(context.Background(), "not found")

	assert.NotNil(t, err)
	assert.Equal(t, "user not found in cache", err.Message)

}

func TestDelete_IDIsRequired(t *testing.T) {
	mock := newMockCache()
	service := NewUserService(mock)
	err := service.DeleteUser(context.Background(), "")

	assert.NotNil(t, err)
	assert.Equal(t, "user ID is required", err.Message)

}
