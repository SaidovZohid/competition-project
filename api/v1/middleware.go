package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/SaidovZohid/competition-project/api/models"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

type Payload struct {
	Id        string `json:"id"`
	UserID    int64  `json:"user_id"`
	Email     string `json:"email"`
	UserType  string `json:"user_type"`
	IssuedAt  string `json:"issued_at"`
	ExpiredAt string `json:"expired_at"`
}

func (h *handlerV1) AuthMiddleware(ctx *gin.Context) {

	authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

	if len(authorizationHeader) == 0 {
		err := errors.New("authorization header is not provided")
		h.logger.WithError(err).Error("authorization header is not provided")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	fields := strings.Fields(authorizationHeader)
	if len(fields) < 2 {
		err := errors.New("invalid authorization header format")
		h.logger.WithError(err).Error("invalid authorization header format")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
	}

	authorizationType := strings.ToLower(fields[0])
	if authorizationType != authorizationTypeBearer {
		err := fmt.Errorf("unsupported authorization type %s", authorizationType)
		h.logger.WithError(err).Error("unsupported authorization type")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	accessToken := fields[1]
	payload, err := h.VerifyToken(accessToken)
	if err != nil {
		h.logger.WithError(err).Error("failed to verify token")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.Set(authorizationPayloadKey, payload)
	ctx.Next()
}

func (h *handlerV1) GetAuthPayload(ctx *gin.Context) (*Payload, error) {
	i, exists := ctx.Get(authorizationPayloadKey)
	if !exists {
		h.logger.Error("not found")
		return nil, errors.New("not found")
	}

	payload, ok := i.(*models.AuthPayload)
	if !ok {
		h.logger.Error("unknown user")
		return nil, errors.New("unknown user")
	}
	return &Payload{
		Id:        payload.ID,
		UserID:    payload.UserID,
		Email:     payload.Email,
		IssuedAt:  payload.IssuedAt.Format(time.RFC3339),
		ExpiredAt: payload.ExpiredAt.Format(time.RFC3339),
	}, nil
}

func (h *handlerV1) VerifyToken(accessToken string) (*models.AuthPayload, error) {
	payload, err := h.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		h.logger.WithError(err).Error("failed to verify token")
		return nil, err
	}

	return &models.AuthPayload{
		ID:        payload.ID.String(),
		UserID:    payload.UserID,
		Email:     payload.Email,
		IssuedAt:  payload.IssuedAt,
		ExpiredAt: payload.ExpiresAt,
	}, nil
}
