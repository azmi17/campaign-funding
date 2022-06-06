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

func (handler *userHandler) Login(c *gin.Context) {
	/*
		user memasukan input (email & password)
		input ditangkap handler
		mapping dari input user ke input struct
		input struct passing ke service
		di service mencari dg bantuan repository user menuju ke email.user
		mencocokan password
	*/
	var input user.LoginInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		// Errors validations
		errors := helper.FormatValidationError(err)
		errorMesage := gin.H{"errors": errors}

		response := helper.ApiResponse("Login failed", http.StatusUnprocessableEntity, "failed", errorMesage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedInUser, err := handler.userService.Login(input)
	if err != nil {
		// Errors validations
		errorMesage := gin.H{"errors": err.Error()}

		response := helper.ApiResponse("Login failed", http.StatusUnprocessableEntity, "failed", errorMesage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(loggedInUser, "u73bD9RD9gTXG8W")

	response := helper.ApiResponse("Successfuly logged in", http.StatusOK, "Success", formatter)

	c.JSON(http.StatusOK, response)
}
