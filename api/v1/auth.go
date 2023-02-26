package v1

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/SaidovZohid/competition-project/api/models"
	emailPkg "github.com/SaidovZohid/competition-project/pkg/email"
	"github.com/SaidovZohid/competition-project/pkg/token"
	"github.com/SaidovZohid/competition-project/pkg/utils"
	"github.com/SaidovZohid/competition-project/storage/repo"
	"github.com/gin-gonic/gin"
)

const (
	RegisterCodeKey   = "register_code_"
	ForgotPasswordKey = "forgot_password_code_"
)

// @Router /auth/register [post]
// @Summary Register a user
// @Description Register a user
// @Tags auth
// @Accept json
// @Produce json
// @Param data body models.RegisterRequest true "Data"
// @Success 200 {object} models.ResponseOK
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) Register(ctx *gin.Context) {
	var (
		req models.RegisterRequest
	)
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		h.logger.WithError(err).Error("failed to bind json to user")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !validatePassword(req.Password) {
		h.logger.WithError(err).Error("failed to validate password")
		ctx.JSON(http.StatusBadRequest, errorResponse(ErrWeakPassword))
		return
	}

	res, _ := h.storage.User().GetByEmail(req.Email)
	if res != nil {
		h.logger.WithError(err).Error("failed to check user by email")
		ctx.JSON(http.StatusBadRequest, errorResponse(ErrEmailExists))
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		h.logger.WithError(err).Error("failed to hash password")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	user := repo.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  hashedPassword,
	}

	userData, err := json.Marshal(user)
	if err != nil {
		h.logger.WithError(err).Error("failed to marshal json")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = h.inMemory.Set("user_"+user.Email, string(userData), 10*time.Minute)
	if err != nil {
		h.logger.WithError(err).Error("failed to set user data to redis")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	go func() {
		err := h.sendVerificationCode(RegisterCodeKey, req.Email)
		if err != nil {
			h.logger.WithError(err).Error("failed to send verfication code")
			fmt.Printf("Failed to send verification code: %v", err)
		}
	}()

	ctx.JSON(http.StatusCreated, models.ResponseOK{
		Message: "Verification code has been sent!",
	})

}

func validatePassword(password string) bool {
	var smallLetter, number bool
	if len(password) < 6 {
		return false
	}

	for i := 0; i < len(password); i++ {
		if password[i] >= 97 && password[i] <= 122 {
			smallLetter = true
		} else if password[i] >= 48 && password[i] <= 57 {
			number = true
		}
	}

	return smallLetter && number
}

func (h *handlerV1) sendVerificationCode(key, email string) error {
	code, err := utils.GenerateRandomCode(6)
	if err != nil {
		return err
	}

	err = h.inMemory.Set(key+email, code, time.Minute*2)
	if err != nil {
		return err
	}

	err = emailPkg.SendEmail(h.cfg, &emailPkg.SendEmailRequest{
		To:      []string{email},
		Subject: "Verification email",
		Body: map[string]string{
			"code": code,
		},
		Type: emailPkg.VerificationEmail,
	})
	if err != nil {
		return err
	}

	return nil
}

// @Router /auth/verify [post]
// @Summary Verify email
// @Description Verify your email which you have used to register
// @Tags auth
// @Accept json
// @Produce json
// @Param data body models.VerifyRequest true "Data"
// @Success 200 {object} models.AuthResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) Verify(c *gin.Context) {
	var (
		req models.VerifyRequest
	)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		h.logger.WithError(err).Error("failed bind json to user")
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	userData, err := h.inMemory.Get("user_" + req.Email)
	if err != nil {
		h.logger.WithError(err).Error("failed to get user data from redis")
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var user repo.User
	err = json.Unmarshal([]byte(userData), &user)
	if err != nil {
		h.logger.WithError(err).Error("failed to unmarshal user")
		c.JSON(http.StatusForbidden, errorResponse(err))
	}

	code, err := h.inMemory.Get(RegisterCodeKey + user.Email)
	if err != nil {
		h.logger.WithError(err).Error("failed to get verfication code from redis")
		c.JSON(http.StatusForbidden, errorResponse(ErrCodeExpired))
		return
	}

	if req.Code != code {
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	result, err := h.storage.User().Create(&user)
	if err != nil {
		h.logger.WithError(err).Error("failed to create user")
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	token, _, err := h.tokenMaker.CreateToken(&token.TokenParams{
		UserID:   result.Id,
		Email:    result.Email,
		Duration: h.cfg.AccessTokenDuration,
	})
	if err != nil {
		h.logger.WithError(err).Error("failed to create token")
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusCreated, models.AuthResponse{
		ID:          user.Id,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt,
		AccessToken: token,
	})
}

// @Router /auth/login [post]
// @Summary Login user
// @Description Login to the service
// @Tags auth
// @Accept json
// @Produce json
// @Param data body models.LoginRequest true "Data"
// @Success 200 {object} models.AuthResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) Login(c *gin.Context) {
	var (
		req models.LoginRequest
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		h.logger.WithError(err).Error("failed to bind json")
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := h.storage.User().GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.logger.WithError(err).Error("failed to get user by email")
			c.JSON(http.StatusForbidden, errorResponse(ErrWrongEmailOrPass))
			return
		}
		h.logger.WithError(err).Error("failed get user by email")
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = utils.CheckPassword(req.Password, user.Password)
	if err != nil {
		h.logger.WithError(err).Error("failed on checking password")
		c.JSON(http.StatusForbidden, errorResponse(ErrWrongEmailOrPass))
		return
	}
	accessToken, _, err := h.tokenMaker.CreateToken(&token.TokenParams{
		UserID:   user.Id,
		Email:    user.Email,
		Duration: h.cfg.AccessTokenDuration,
	})
	if err != nil {
		h.logger.WithError(err).Error("failed access create token")
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, models.LoginRes{
		User: models.User{
			ID:        user.Id,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
		AccessToken: accessToken,
	})
}
