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
		Name:  "João",
		CPF:   "12345678900",
		Age:   30,
		Email: "joao@gmail.com",
	}

	response, err := service.Execute(context.Background(), user)
	assert.Nil(t, err)
	assert.NotEmpty(t, response.Token)
	assert.Equal(t, user.Name, response.User.Name)
	assert.Equal(t, user.Email, response.User.Email)

	// verifica se o dado foi realmente salvo no cache
	data, ok := mock.Get(context.Background(), response.User.ID)
	assert.True(t, ok)

	var cached model.User
	unmarshalErr := json.Unmarshal(data, &cached)
	assert.NoError(t, unmarshalErr)
	assert.Equal(t, response.User.ID, cached.ID)
	assert.Equal(t, response.User.Email, cached.Email)
}

func TestCreateUser_OverWrite(t *testing.T) {
	mock := cache.NewMockCache()
	service := NewCreateUser(mock)

	user1 := &model.User{Name: "first", CPF: "12345678910", Age: 30, Email: "1@a.com"}
	user2 := &model.User{Name: "seconds", CPF: "12345678911", Age: 25, Email: "2@a.com"}

	_, _ = service.Execute(context.Background(), user1)
	response2, _ := service.Execute(context.Background(), user2)

	data, _ := mock.Get(context.Background(), response2.User.ID)
	var cached model.User
	unmarshalErr := json.Unmarshal(data, &cached)
	assert.NoError(t, unmarshalErr)
	assert.Equal(t, "seconds", cached.Name)
}

func TestCreateUserWithEmailError(t *testing.T) {
	mock := cache.NewMockCache()
	service := NewCreateUser(mock)

	user := &model.User{
		Name:  "João",
		CPF:   "12345678900",
		Age:   30,
		Email: "joao@.com",
	}

	_, err := service.Execute(context.Background(), user)
	assert.NotEmpty(t, err)
}
