package v1

import (
	"net/http"
	"strconv"

	"github.com/SaidovZohid/competition-project/api/models"
	"github.com/SaidovZohid/competition-project/storage/repo"
	"github.com/gin-gonic/gin"
)

// @Security ApiKeyAuth
// @Router /users/{id} [get]
// @Summary Get user by id
// @Description Get user by id
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.User
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	resp, err := h.storage.User().Get(int64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, parseUserModel(resp))
}

// @Router /users [get]
// @Summary Get all users
// @Description Get all users
// @Tags user
// @Accept json
// @Produce json
// @Param filter query models.GetAllParams false "Filter"
// @Success 200 {object} models.GetAllUsersResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetAllUsers(c *gin.Context) {
	req, err := validateGetAllParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := h.storage.User().GetAll(&repo.GetAllUsersParams{
		Page:   req.Page,
		Limit:  req.Limit,
		Search: req.Search,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, getUsersResponse(result))
}

// @Router /users/email/{email} [get]
// @Summary Get user by email
// @Description Get user by email
// @Tags user
// @Accept json
// @Produce json
// @Param email path string true "Email"
// @Success 200 {object} models.User
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")

	resp, err := h.storage.User().GetByEmail(email)
	if err != nil {
		h.logger.WithError(err).Error("failed to get user by email")
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, parseUserModel(resp))
}

// @Security ApiKeyAuth
// @Router /user/{id} [delete]
// @Summary Delete user by id
// @Description Delete user by id
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 201 {object} models.ResponseO
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = h.storage.User().Delete(int64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, models.ResponseOK{
		Message: "success",
	})
}

// @Security ApiKeyAuth
// @Router /users [put]
// @Summary Update a user
// @Description Update a user
// @Tags user
// @Accept json
// @Produce json
// @Param user body models.UpdateUserRequest true "User"
// @Success 201 {object} models.User
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) UpdateUser(c *gin.Context) {
	var (
		req models.UpdateUserRequest
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	resp, err := h.storage.User().Update(&repo.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, parseUserModel(resp))
}

func getUsersResponse(data *repo.GetAllUsersResult) *models.GetAllUsersResponse {
	response := models.GetAllUsersResponse{
		Users: make([]*models.User, 0),
		Count: data.Count,
	}

	for _, user := range data.Users {
		u := parseUserModel(user)
		response.Users = append(response.Users, &u)
	}

	return &response
}

func parseUserModel(user *repo.User) models.User {
	return models.User{
		ID:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}
