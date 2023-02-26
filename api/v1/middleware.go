package v1

import (
	"errors"
	"net/http"
	"time"

	"github.com/SaidovZohid/competition-project/api/models"
	"github.com/gin-gonic/gin"
)

type Payload struct {
	Id        string `json:"id"`
	UserID    int64  `json:"user_id"`
	Email     string `json:"email"`
	IssuedAt  string `json:"issued_at"`
	ExpiredAt string `json:"expired_at"`
}

func (h *handlerV1) AuthMiddleware(ctx *gin.Context) {

	accessToken := ctx.GetHeader(h.cfg.AuthHeaderKey)

	if len(accessToken) == 0 {
		err := errors.New("authorization header is not provided")
		h.logger.Error(err)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	payload, err := h.VerifyToken(accessToken)
	if err != nil {
		h.logger.WithError(err).Error("failed to verify token")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.Set(h.cfg.AuthPayloadKey, Payload{
		Id:        payload.ID,
		UserID:    payload.UserID,
		Email:     payload.Email,
		IssuedAt:  payload.IssuedAt.Format(time.RFC3339),
		ExpiredAt: payload.ExpiredAt.Format(time.RFC3339),
	})
	ctx.Next()
}

func (h *handlerV1) GetAuthPayload(ctx *gin.Context) (*Payload, error) {
	i, exists := ctx.Get(h.cfg.AuthPayloadKey)
	if !exists {
		h.logger.Error("not found")
		return nil, errors.New("not found")
	}

	payload, ok := i.(Payload)
	if !ok {
		h.logger.Error("unknown user")
		return nil, errors.New("unknown user")
	}
	return &Payload{
		Id:        payload.Id,
		UserID:    payload.UserID,
		Email:     payload.Email,
		IssuedAt:  payload.IssuedAt,
		ExpiredAt: payload.ExpiredAt,
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
