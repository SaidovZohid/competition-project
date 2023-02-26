package v1

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/SaidovZohid/competition-project/api/models"
	"github.com/SaidovZohid/competition-project/pkg/utils"
	"github.com/SaidovZohid/competition-project/storage/repo"
	"github.com/gin-gonic/gin"
	_ "github.com/skip2/go-qrcode"
)

// @Security ApiKeyAuth
// @Router /urls/make-short-url [post]
// @Summary Make short url
// @Description Make your long url short
// @Tags url
// @Accept json
// @Produce json
// @Param data query models.CreateShortUrlRequest true "Data"
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
		return
	}
	if req.Duration != "" {
		duration, err = time.ParseDuration(req.Duration)
		if err != nil {
			h.logger.WithError(err).Error("failed to parse string to time.Duration")
			ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInternalServer))
			return
		}
	}
	var shortUrl string
	if req.CustomUrl == "" {
		str := utils.RandomString(6)
		shortUrl = fmt.Sprintf("http://localhost%s/v1/urls/%s", h.cfg.HttpPort, str)
		err = h.inMemory.Set(shortUrl, req.OriginalUrl, duration)
		if err != nil {
			h.logger.WithError(err).Error("failed to set url to redis db")
			ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInternalServer))
			return
		}
	} else if req.CustomUrl != "" {
		customUrl := fmt.Sprintf("http://localhost%s/v1/urls/%s", h.cfg.HttpPort, req.CustomUrl)
		url, err := h.storage.Url().Get(customUrl)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				shortUrl = fmt.Sprintf("http://localhost%s/v1/urls/%s", h.cfg.HttpPort, req.CustomUrl)
				err = h.inMemory.Set(customUrl, req.OriginalUrl, duration)
				if err != nil {
					h.logger.WithError(err).Error("failed to set url to redis db")
					ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInternalServer))
					return
				}
			}
		} else if url.HashedUrl == customUrl {
			h.logger.WithError(err).Error("urls are identical")
			ctx.JSON(http.StatusBadRequest, errorResponse(ErrUrlUnavailable))
			return
		} else {
			if err != nil {
				h.logger.WithError(err).Error("failed to get url")
				ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInternalServer))
				return
			}
		}
	}

	expiresAt := time.Now().Add(duration)
	url, err := h.storage.Url().Create(&repo.Url{
		UserId:      payload.UserID,
		OriginalUrl: req.OriginalUrl,
		HashedUrl:   shortUrl,
		MaxClicks:   req.MaxClicks,
		ExpiresAt:   &expiresAt,
	})
	if err != nil {
		h.logger.WithError(err).Error("failed create user")
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInternalServer))
		return
	}

	ctx.JSON(http.StatusOK, parseUrlModel(url))
}

func parseUrlModel(data *repo.Url) *models.Url {
	var clicks *int64
	if data.MaxClicks != 0 {
		clicks = &data.MaxClicks
	}
	return &models.Url{
		Id:          data.Id,
		UserId:      data.UserId,
		OriginalUrl: data.OriginalUrl,
		HashedUrl:   data.HashedUrl,
		MaxClicks:   clicks,
		ExpiresAt:   data.ExpiresAt,
		CreatedAt:   data.CreatedAt.Format(time.RFC3339),
	}
}

// @Security ApiKeyAuth
// @Router /urls/generate-qr-code [post]
// @Summary Generate QR code
// @Description Generate QR code to store your url in it
// @Tags url
// @Accept json
// @Produce json
// @Param data query models.CreateShortUrlRequest true "Data"
// @Success 200 {object} models.Url
// @Failure 500 {object} models.ErrorResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
func (h *handlerV1) GenerateQRCode(ctx *gin.Context) {

}
