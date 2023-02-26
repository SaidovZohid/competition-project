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

	identifier := utils.GenerateShortLink(req.OriginalUrl, strconv.Itoa(int(payload.UserID)))
	shortURL := fmt.Sprintf("http://localhost%s/v1/urls/%s", h.cfg.HttpPort, identifier)
	err = h.inMemory.Set(shortURL, req.OriginalUrl, duration)
	if err != nil {
		h.logger.WithError(err).Error("failed to set url to redis db")
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInternalServer))
		return
	}
	url, err := h.storage.Url().Create(&repo.Url{
		UserId:      payload.UserID,
		OriginalUrl: req.OriginalUrl,
		HashedUrl:   shortURL,
		MaxClicks:   &req.MaxClicks,
	})
	if err != nil {
		h.logger.WithError(err).Error("failed create user")
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInternalServer))
		return
	}

	ctx.JSON(http.StatusOK, parseUrlModel(url))
}

// @Router /urls/{shorturl} [get]
// @Summary Redirect short url
// @Description Redirect url by giving short url to original url
// @Tags url
// @Accept json
// @Param shorturl path string true "ShortUrl"
// @Success 302 {object} models.Url
// @Failure 500 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
func (h *handlerV1) RedirectUrl(ctx *gin.Context) {
	url := fmt.Sprintf("http://localhost%s", h.cfg.HttpPort+ctx.Request.URL.Path)
	url1, err := h.storage.Url().Get(url)
	if err != nil {
		h.logger.WithError(err).Error("failed to get url")
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInternalServer))
		return
	}

	if url1.ExpiresAt != nil {
		if time.Now().After(*url1.ExpiresAt) {
			h.logger.WithError(err).Error("time expired")
			ctx.JSON(http.StatusNotFound, errorResponse(ErrNotFound))
			return
		}
	}
	if url1.MaxClicks != nil {
		if *url1.MaxClicks <= 0 {
			h.logger.WithError(err).Error("max click is over")
			ctx.JSON(http.StatusNotFound, errorResponse(ErrNotFound))
			return
		}
	}
	err = h.storage.Url().DecrementClick(url)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInternalServer))
		return
	}

	ctx.Redirect(http.StatusFound, url1.OriginalUrl)
}

// @Security ApiKeyAuth
// @Router /urls/hashed-url{id} [get]
// @Summary Get url by hashed url
// @Description Get url by hashed url
// @Tags url
// @Accept json
// @Produce json
// @Param url path int true "URL"
// @Success 200 {object} models.Url
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetUrl(c *gin.Context) {
	url := c.Param("url")

	resp, err := h.storage.Url().Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, parseUrlModel(resp))
}

// @Security ApiKeyAuth
// @Router /user/{id} [delete]
// @Summary Delete user by id
// @Description Delete user by id
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param user_id path int true "UserID"
// @Success 201 {object} models.ResponseOK
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) DeleteUrl(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user_id, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = h.storage.Url().Delete(int64(id), int64(user_id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, models.ResponseOK{
		Message: "success",
	})
}

// @Security ApiKeyAuth
// @Router /urls [put]
// @Summary Update a url
// @Description Update a url
// @Tags url
// @Accept json
// @Produce json
// @Param url body models.UpdateUrlRequest true "Url"
// @Success 201 {object} models.Url
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) UpdateUrl(c *gin.Context) {
	var (
		req models.UpdateUrlRequest
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	resp, err := h.storage.Url().Update(&repo.Url{
		HashedUrl: req.HashedUrl,
		MaxClicks: req.MaxClicks,
		ExpiresAt: req.ExpiresAt,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, parseUrlModel(resp))
}

func parseUrlModel(data *repo.Url) *models.Url {
	return &models.Url{
		Id:          data.Id,
		UserId:      data.UserId,
		OriginalUrl: data.OriginalUrl,
		HashedUrl:   data.HashedUrl,
		MaxClicks:   data.MaxClicks,
		ExpiresAt:   data.ExpiresAt,
		CreatedAt:   data.CreatedAt.Format(time.RFC3339),
	}
}
