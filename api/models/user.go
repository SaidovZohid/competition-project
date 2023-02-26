package models

import "time"

type User struct {
	ID          int64     `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateUserRequest struct {
	FirstName   string  `json:"first_name" binding:"required,min=2,max=30"`
	LastName    string  `json:"last_name" binding:"required,min=2,max=30"`
	Email       string  `json:"email" binding:"required,email"`
	Password    string  `json:"password" binding:"required,min=6,max=16"`
}

type UpdateUserRequest struct {
	FirstName   string  `json:"first_name" binding:"required,min=2,max=30"`
	LastName    string  `json:"last_name" binding:"required,min=2,max=30"`
}

type GetAllUsersResponse struct {
	Users []*User `json:"users"`
	Count int32   `json:"count"`
}


