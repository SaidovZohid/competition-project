package postgres_test

import (
	"testing"

	"github.com/SaidovZohid/competition-project/storage/repo"
	"github.com/bxcodec/faker/v4"
	"github.com/stretchr/testify/require"
)

func createUrl(t *testing.T) *repo.Url {
	user := createUser(t)

	url, err := strg.Url().Create(&repo.Url{
		UserId: user.Id,
		OriginalUrl: faker.URL(),
		HashedUrl: faker.URL(),
		MaxClicks: faker.RandomUnixTime(),
	})
	require.NoError(t, err)
	require.NotEmpty(t, url)

	return url
}

func TestGetUrl(t *testing.T) {
	c := createUrl(t)

	url, err := strg.Url().Get(c.Id)
	require.NoError(t, err)
	require.NotEmpty(t, url)
}

func TestCreateUrl(t *testing.T) {
	createUrl(t)
}

func TestDeleteUrl(t *testing.T) {
	c := createUrl(t)

	url, err := strg.Url().Delete(&repo.Url{Id: c.Id})
	require.NoError(t, err)
	require.NotEmpty(t, url)
}

func TestGetAllUrls(t *testing.T) {
	url, err := strg.Url().GetAll(&repo.GetAllUrlsParams{
		Limit:  3,
		Page:   1,
		Search: "ab",
	})
	require.NoError(t, err)
	require.NotEmpty(t, url)
}