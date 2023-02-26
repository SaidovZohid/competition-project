package v1

import (
	"errors"
	"strconv"

	"github.com/SaidovZohid/competition-project/api/models"
	"github.com/SaidovZohid/competition-project/config"
	"github.com/SaidovZohid/competition-project/pkg/logger"
	"github.com/SaidovZohid/competition-project/pkg/token"
	"github.com/SaidovZohid/competition-project/storage"
	"github.com/gin-gonic/gin"
)

var (
	ErrWrongEmailOrPassword = errors.New("INVALID_EMAIL_OR_PASSWORD")
	ErrInternalServer       = errors.New("INTERNAL_SERVER_ERROR")
	ErrNotFound             = errors.New("NOT_FOUND")
	ErrUnknownUser          = errors.New("UNKNOWN_USER")
	ErrEmailExists          = errors.New("EMAIL_EXISTS")
	ErrMemberExists         = errors.New("MEMBER_IN_COMPANY")
	ErrIncorrectCode        = errors.New("INCORRECT_VERIFICATION_CODE")
	ErrCodeExpired          = errors.New("VERIFICATION_CODE_IS_EXPIRED")
	ErrForbidden            = errors.New("FORBIDDEN")
	ErrBadRequest           = errors.New("BAD_REQUEST")
	ErrFailedDowload        = errors.New("FAILED_TO_DOWNLOAD")
	ErrSize                 = errors.New("INVALID_SIZE")
	ErrFileType             = errors.New("WRONG_TYPE")
	ErrUnauthorized         = errors.New("UNAUTHORIZED")
	ErrMethodNotAllowed     = errors.New("METHOD_NOT_ALLOWED")
	ErrNotAllowed           = errors.New("NOT_ALLOWED")
	ErrWrongEmailOrPass     = errors.New("wrong email or password")
	ErrWeakPassword         = errors.New("password must contain at least one small letter, one number and be at least 6 characters long")
)

type handlerV1 struct {
	cfg        *config.Config
	storage    storage.StorageI
	inMemory   storage.InMemoryStorageI
	tokenMaker token.Maker
	logger     *logger.Logger
}

type HandlerV1Options struct {
	Cfg        *config.Config
	Storage    storage.StorageI
	InMemory   storage.InMemoryStorageI
	TokenMaker token.Maker
	Logger     *logger.Logger
}

func New(options *HandlerV1Options) *handlerV1 {
	return &handlerV1{
		cfg:        options.Cfg,
		storage:    options.Storage,
		inMemory:   options.InMemory,
		tokenMaker: options.TokenMaker,
		logger:     options.Logger,
	}
}

func errorResponse(err error) *models.ErrorResponse {
	return &models.ErrorResponse{
		Error: err.Error(),
	}
}

func validateUrlParams(ctx *gin.Context) (*models.CreateShortUrlRequest, error) {
	var (
		maxClicks int = 0
		err       error
	)

	if ctx.Query("max_clicks") != "" {
		maxClicks, err = strconv.Atoi(ctx.Query("max_clicks"))
		if err != nil {
			return nil, err
		}
	}

	return &models.CreateShortUrlRequest{
		OriginalUrl: ctx.Query("original_url"),
		Duration:    ctx.Query("duration"),
		MaxClicks:   int64(maxClicks),
	}, nil
}

func validateGetAllParams(c *gin.Context) (*models.GetAllParams, error) {
	var (
		limit int = 10
		page  int = 1
		err   error
	)

	if c.Query("limit") != "" {
		limit, err = strconv.Atoi(c.Query("limit"))
		if err != nil {
			return nil, err
		}
	}

	if c.Query("page") != "" {
		page, err = strconv.Atoi(c.Query("page"))
		if err != nil {
			return nil, err
		}
	}

	return &models.GetAllParams{
		Limit:  int32(limit),
		Page:   int32(page),
		Search: c.Query("search"),
	}, nil
}
