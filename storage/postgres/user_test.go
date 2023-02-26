package postgres_test

import (
	"testing"

	"github.com/SaidovZohid/competition-project/pkg/utils"
	"github.com/SaidovZohid/competition-project/storage/repo"
	"github.com/bxcodec/faker/v4"
	"github.com/stretchr/testify/require"
)

func createUser(t *testing.T) *repo.User {
	hashedPassword, err := utils.HashPassword(faker.Password())
	require.NoError(t, err)
	u := repo.User{
		FirstName: faker.FirstName(),
		LastName:  faker.LastName(),
		Email:     faker.Email(),
		Password:  hashedPassword,
	}
	user, err := strg.User().Create(&u)
	require.NoError(t, err)
	require.NotZero(t, user.Id)
	require.Equal(t, u.FirstName, u.FirstName)
	require.Equal(t, u.LastName, u.LastName)
	require.Equal(t, u.Email, u.Email)
	require.NotZero(t, user.CreatedAt)

	return user
}

func deleteUser(t *testing.T, id int64) {
	err := strg.User().DeleteUser(id)
	require.NoError(t, err)
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

	deleteUser(t, c.Id)
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
