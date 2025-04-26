package deletecases

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/odanaraujo/user-api/cache"
	"github.com/odanaraujo/user-api/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestDeleteUser(t *testing.T) {
	mock := cache.NewMockCache()
	service := NewDeleteUser(mock)

	user := model.User{ID: "999"}
	data, _ := json.Marshal(user)
	mock.Set(context.Background(), user.ID, data, time.Hour)

	err := service.Execute(context.Background(), user.ID)
	assert.Nil(t, err)

	_, ok := mock.Get(context.Background(), user.ID)
	assert.False(t, ok)
}

func TestDelete_NotFound(t *testing.T) {
	mock := cache.NewMockCache()
	service := NewDeleteUser(mock)
	err := service.Execute(context.Background(), "not found")

	assert.NotNil(t, err)
	assert.Equal(t, "user not found in cache", err.Message)

}

func TestDelete_IDIsRequired(t *testing.T) {
	mock := cache.NewMockCache()
	service := NewDeleteUser(mock)
	err := service.Execute(context.Background(), "")

	assert.NotNil(t, err)
	assert.Equal(t, "user ID is required", err.Message)

}
