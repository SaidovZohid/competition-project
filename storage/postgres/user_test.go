package postgres_test

import (
	"database/sql"
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
	err := strg.User().Delete(id)
	require.NoError(t, err)
}
func TestCreateUser(t *testing.T) {
	u := createUser(t)
	deleteUser(t, u.Id)
}

func TestGetUser(t *testing.T) {
	c := createUser(t)
	user, err := strg.User().Get(c.Id)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	user2, err := strg.User().Get(-1)
	require.Error(t, err, sql.ErrNoRows)
	require.Empty(t, user2)
	deleteUser(t, user.Id)
}

func TestDeleteUser(t *testing.T) {
	c := createUser(t)
	deleteUser(t, c.Id)
	err := strg.User().Delete(-1)
	require.Error(t, err, sql.ErrNoRows)
}

func TestGetAllUsers(t *testing.T) {
	ids := make([]int64, 0)
	for i := 0; i < 10; i++ {
		u := createUser(t)
		ids = append(ids, u.Id)
	}

	users, err := strg.User().GetAll(&repo.GetAllUsersParams{
		Limit: 10,
		Page:  1,
	})
	require.NoError(t, err)
	require.NotEmpty(t, users)
	require.GreaterOrEqual(t, users.Count, int32(10))
	for i := 0; i < 10; i++ {
		deleteUser(t, ids[i])
	}
}

func TestGetByEmail(t *testing.T) {
	user := createUser(t)
	user2, err := strg.User().GetByEmail(user.Email)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	deleteUser(t, user.Id)
}

func TestUpdateUser(t *testing.T) {
	user := createUser(t)
	user2, err := strg.User().Update(&repo.User{
		Id:        user.Id,
		FirstName: faker.FirstName(),
		LastName:  faker.LastName(),
	})
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	deleteUser(t, user.Id)
	user3, err := strg.User().Update(&repo.User{
		Id:        user.Id,
		FirstName: faker.FirstName(),
		LastName:  faker.LastName(),
	})
	require.Error(t, err, sql.ErrNoRows)
	require.Nil(t, user3)
}
