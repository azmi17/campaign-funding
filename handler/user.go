package handler

import (
	"go-campaign-funding/helper"
	"go-campaign-funding/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHanlder(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (handler *userHandler) RegisterUser(c *gin.Context) {

	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {

		// Errors validations
		errors := helper.FormatValidationError(err)
		errorMesage := gin.H{"errors": errors}

		response := helper.ApiResponse("Register account failed", http.StatusUnprocessableEntity, "failed", errorMesage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newUser, err := handler.userService.RegisterUser(input)

	// if true // debug check
	if err != nil {
		response := helper.ApiResponse("Register account failed", http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// token, err := h.jwtService.GenerateToken()

	formatter := user.FormatUser(newUser, "u73bD9RD9gTXG8W")

	response := helper.ApiResponse("Account has been registered", http.StatusOK, "Success", formatter)

	c.JSON(http.StatusOK, response)
}
