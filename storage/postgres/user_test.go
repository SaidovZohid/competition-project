package postgres_test

import (
	"testing"

	"github.com/SaidovZohid/competition-project/storage/repo"
	"github.com/stretchr/testify/require"
	"github.com/bxcodec/faker/v4"
)

func createUser(t *testing.T) *repo.User {
	user, err := strg.User().Create(&repo.User{
		FirstName: faker.FirstName(),
		LastName:  faker.LastName(),
		Email:     faker.Email(),
		Password:  faker.Password(),
	})
	require.NoError(t, err)
	require.NotEmpty(t, user)

	return user
}

func TestGetUser(t *testing.T) {
	c := createUser(t)

	user, err := strg.User().Get(c.Id)
	require.NoError(t, err)
	require.NotEmpty(t, user)
}

func TestCreateUser(t *testing.T) {
	createUser(t)
}

func TestDeleteUser(t *testing.T) {
	c := createUser(t)

	user, err := strg.User().DeleteUser(&repo.User{Id: c.Id})
	require.NoError(t, err)
	require.NotEmpty(t, user)
}

func TestGetAllUsers(t *testing.T) {
	user, err := strg.User().GetAll(&repo.GetAllUsersParams{
		Limit:  3,
		Page:   1,
		Search: "ab",
	})
	require.NoError(t, err)
	require.NotEmpty(t, user)
}
