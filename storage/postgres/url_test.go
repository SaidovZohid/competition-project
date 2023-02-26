package postgres_test

import (
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/SaidovZohid/competition-project/storage/repo"
	"github.com/bxcodec/faker/v4"
	"github.com/stretchr/testify/require"
)

func createUrl(t *testing.T) *repo.Url {
	user := createUser(t)
	click := int64(100)
	tm := time.Now().Add(time.Hour * 24)
	url1 := repo.Url{
		UserId:      user.Id,
		OriginalUrl: faker.URL(),
		HashedUrl:   faker.URL(),
		MaxClicks:   &click,
		ExpiresAt:   &tm,
	}
	url2, err := strg.Url().Create(&url1)
	require.NoError(t, err)
	require.NotZero(t, url1.Id)
	require.Equal(t, url1.OriginalUrl, url2.OriginalUrl)
	require.Equal(t, url1.HashedUrl, url2.HashedUrl)
	require.Equal(t, url1.MaxClicks, url2.MaxClicks)
	require.NotZero(t, url2.ExpiresAt)
	require.NotZero(t, url2.CreatedAt)

	return url2
}

func deleteUrl(t *testing.T, id, userID int64) {
	err := strg.Url().Delete(id, userID)
	require.NoError(t, err)
}
func TestCreateUrl(t *testing.T) {
	url := createUrl(t)
	deleteUser(t, url.UserId)
}

func TestGetUrl(t *testing.T) {
	url := createUrl(t)
	url2, err := strg.Url().Get(url.HashedUrl)
	require.NoError(t, err)
	require.Equal(t, url.Id, url2.Id)
	require.Equal(t, url.UserId, url2.UserId)
	require.Equal(t, url.HashedUrl, url2.HashedUrl)
	require.Equal(t, url.OriginalUrl, url.OriginalUrl)
	require.Equal(t, url.MaxClicks, url2.MaxClicks)
	require.NotZero(t, url2.ExpiresAt)
	require.NotZero(t, url2.CreatedAt)
	log.Println(url2.UserId)
	deleteUser(t, url2.UserId)
}

func TestDeleteUrl(t *testing.T) {
	url := createUrl(t)
	deleteUrl(t, url.Id, url.UserId)
	url = createUrl(t)
	err := strg.Url().Delete(-1, url.UserId)
	require.Error(t, err, sql.ErrNoRows)
	deleteUser(t, url.UserId)
}

func TestUpdateUrl(t *testing.T) {
	url := createUrl(t)
	click := int64(100)
	url2, err := strg.Url().Update(&repo.Url{
		Id:        url.Id,
		UserId:    url.UserId,
		HashedUrl: faker.URL(),
		MaxClicks: &click,
	})
	require.NoError(t, err)
	require.NotEmpty(t, url2)
	deleteUser(t, url.UserId)
}

func TestDecrementMaxClick(t *testing.T) {
	url := createUrl(t)
	log.Println(url.MaxClicks)
	err := strg.Url().DecrementClick(url.HashedUrl)
	require.NoError(t, err)
	deleteUser(t, url.UserId)
}

func TestGetAllUrl(t *testing.T) {
	userId := make([]int64, 0)
	for i := 0; i < 10; i++ {
		url := createUrl(t)
		userId = append(userId, url.UserId)
	}
	urls, err := strg.Url().GetAll(&repo.GetAllUrlsParams{
		Limit: 10,
		Page:  1,
	})
	require.NoError(t, err)
	require.GreaterOrEqual(t, urls.Count, int32(10))

	urls2, _ := strg.Url().GetAll(&repo.GetAllUrlsParams{
		Limit:  10,
		Page:   1,
		UserID: -1,
	})
	require.Empty(t, urls2.Urls)
	require.Equal(t, urls2.Count, int32(0))
	for i := 0; i < 10; i++ {
		deleteUser(t, userId[i])
	}
}
