package createcases

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/odanaraujo/user-api/cache"
	"github.com/odanaraujo/user-api/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	mock := cache.NewMockCache()
	service := NewCreateUser(mock)

	user := &model.User{
		ID:    "123",
		Name:  "João",
		CPF:   "12345678900",
		Age:   30,
		Email: "joao@gmail.com",
	}

	createdUser, err := service.Execute(context.Background(), user)
	assert.Nil(t, err)
	assert.Equal(t, user, createdUser)

	// verifica se o dado foi realmente salvo no cache
	data, ok := mock.Get(context.Background(), createdUser.ID)
	assert.True(t, ok)

	var cached model.User
	_ = json.Unmarshal(data, &cached)
	assert.Equal(t, user.ID, cached.ID)
	assert.Equal(t, user.Email, cached.Email)
}

func TestCreateUser_OverWrite(t *testing.T) {
	mock := cache.NewMockCache()
	service := NewCreateUser(mock)

	user1 := &model.User{ID: "1", Name: "first", CPF: "12345678910", Age: 30, Email: "1@a.com"}
	user2 := &model.User{ID: "1", Name: "seconds", CPF: "12345678911", Age: 25, Email: "2@a.com"}

	service.Execute(context.Background(), user1)
	service.Execute(context.Background(), user2)

	data, _ := mock.Get(context.Background(), user2.ID)
	var cached model.User
	json.Unmarshal(data, &cached)

	assert.Equal(t, "seconds", cached.Name)
}

func TestCreateUserWithEmailError(t *testing.T) {
	mock := cache.NewMockCache()
	service := NewCreateUser(mock)

	user := &model.User{
		ID:    "123",
		Name:  "João",
		CPF:   "12345678900",
		Age:   30,
		Email: "joao@.com",
	}

	_, err := service.Execute(context.Background(), user)
	assert.NotEmpty(t, err)
}
