package repo

import "time"

type User struct {
	Id        int64
	FirstName string
	LastName  string
	Email     string
	Password  string
	CreatedAt time.Time
}

type GetAllUsersResult struct {
	Users []*User
	Count int32
}

type GetAllUsersParams struct {
	Limit  int32
	Page   int32
	Search string
}

type UpdatePassword struct {
	UserId int64
	Password string
}

type UserStorageI interface {
	Create(u *User) (*User, error)
	Get(id int64) (*User, error)
	GetByEmail(email string) (*User, error)
	GetAll(params *GetAllUsersParams) (*GetAllUsersResult, error)
	UpdateUser(u *User) (*User, error)
	DeleteUser(u *User) (*User, error)
}
