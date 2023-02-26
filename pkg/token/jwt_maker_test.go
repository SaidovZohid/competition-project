package token

import (
	"testing"
	"time"

	"github.com/SaidovZohid/competition-project/pkg/utils"

	"github.com/bxcodec/faker/v4"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(utils.RandomString(32))
	require.NoError(t, err)

	userID := utils.RandomInt(1, 10)
	userEmail := faker.Email()
	userType := "user"
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(&TokenParams{
		UserID:   userID,
		Email:    userEmail,
		UserType: userType,
		Duration: duration,
	})
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, userID, payload.UserID)
	require.Equal(t, userEmail, payload.Email)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiresAt, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(utils.RandomString(32))
	require.NoError(t, err)

	userID := utils.RandomInt(1, 10)
	userEmail := faker.Email()
	userType := "user"

	token, payload, err := maker.CreateToken(&TokenParams{
		UserID:   userID,
		Email:    userEmail,
		UserType: userType,
		Duration: -time.Minute,
	})
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	userID := utils.RandomInt(1, 10)
	userEmail := faker.Email()
	userType := "user"

	payload, err := NewPayload(&TokenParams{
		UserID:   userID,
		Email:    userEmail,
		UserType: userType,
		Duration: time.Minute,
	})
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(utils.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
