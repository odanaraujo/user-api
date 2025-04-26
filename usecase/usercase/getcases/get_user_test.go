package getcases

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/odanaraujo/user-api/cache"
	"github.com/odanaraujo/user-api/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByID(t *testing.T) {
	mock := cache.NewMockCache()
	service := NewGetUserCase(mock)

	user := model.User{
		ID:    "abc",
		Name:  "Maria",
		CPF:   "98765432100",
		Age:   25,
		Email: "maria@email.com",
	}

	data, _ := json.Marshal(user)
	mock.Set(context.Background(), user.ID, data, time.Hour)

	got, err := service.Execute(context.Background(), user.ID)
	assert.Nil(t, err)
	assert.Equal(t, user.Name, got.Name)
	assert.Equal(t, user.Email, got.Email)
}

func TestGetUserByID_NotFound(t *testing.T) {
	mock := cache.NewMockCache()
	service := NewGetUserCase(mock)

	got, err := service.Execute(context.Background(), "não existe")
	assert.NotNil(t, err)
	assert.Nil(t, got)
}

func TestGetUserByID_InvalidJSON(t *testing.T) {
	mock := cache.NewMockCache()
	service := NewGetUserCase(mock)

	mock.Set(context.Background(), "abc", []byte("this_not_json"), time.Hour)
	got, err := service.Execute(context.Background(), "abc")
	assert.NotNil(t, err)
	assert.Nil(t, got)
	assert.Equal(t, "error performing unmarshal", err.Message)
}
