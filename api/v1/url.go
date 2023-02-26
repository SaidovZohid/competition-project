package v1

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/SaidovZohid/competition-project/api/models"
	"github.com/SaidovZohid/competition-project/pkg/utils"
	"github.com/SaidovZohid/competition-project/storage/repo"
	"github.com/gin-gonic/gin"
)

// @Security ApiKeyAuth
// @Router /urls/make-short-url [post]
// @Summary Make short url
// @Description Make your long url short
// @Tags url
// @Accept json
// @Produce json
// @Param data body models.CreateShortUrlRequest true "Data"
// @Success 200 {object} models.Url
// @Failure 500 {object} models.ErrorResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
func (h *handlerV1) MakeShortUrl(ctx *gin.Context) {
	var (
		duration time.Duration
		err      error
	)
	req, err := validateUrlParams(ctx)
	if err != nil {
		h.logger.WithError(err).Error("failed to validate url params")
		ctx.JSON(http.StatusBadRequest, errorResponse(ErrBadRequest))
		return
	}

	payload, err := h.GetAuthPayload(ctx)
	if err != nil {
		h.logger.WithError(err).Error("failed to get authorization payload")
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrUnauthorized))
	}
	if req.Duration != "" {
		duration, err = time.ParseDuration(req.Duration)
		if err != nil {
			h.logger.WithError(err).Error("failed to parse string to time.Duration")
			ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInternalServer))
			return
		}
	}

	identifier := utils.GenerateShortLink(req.OriginalUrl, strconv.Itoa(int(payload.UserID)))
	shortURL := fmt.Sprintf("http://localhost%s/%s", h.cfg.HttpPort, identifier)
	err = h.inMemory.Set(identifier, req.OriginalUrl, duration)
	if err != nil {
		h.logger.WithError(err).Error("failed to set url to redis db")
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInternalServer))
		return
	}

	url, err := h.storage.Url().Create(&repo.Url{
		UserId:      payload.UserID,
		OriginalUrl: req.OriginalUrl,
		HashedUrl:   shortURL,
		MaxClicks:   req.MaxClicks,
	})
	if err != nil {
		h.logger.WithError(err).Error("failed create user")
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInternalServer))
		return
	}

	ctx.JSON(http.StatusOK, parseUrlModel(url))
}

func parseUrlModel(data *repo.Url) *models.Url {
	expiresAt := data.ExpiresAt.Format(time.RFC3339)
	return &models.Url{
		Id:          data.Id,
		UserId:      data.UserId,
		OriginalUrl: data.OriginalUrl,
		HashedUrl:   data.HashedUrl,
		MaxClicks:   data.MaxClicks,
		ExpiresAt:   &expiresAt,
		CreatedAt:   data.CreatedAt.Format(time.RFC3339),
	}
}
