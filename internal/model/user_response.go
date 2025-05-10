package model

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   uint64 `json:"age"`
}

type CreateUserResponse struct {
	User  *UserResponse `json:"user"`
	Token string        `json:"token"`
}

func NewUserResponse(user *User) *UserResponse {
	return &UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	}
}
