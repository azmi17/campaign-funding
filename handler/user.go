package handler

import (
	"fmt"
	"go-campaign-funding/auth"
	"go-campaign-funding/helper"
	"go-campaign-funding/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHanlder(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
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

	token, err := handler.authService.GenerateToken(newUser.ID)
	if err != nil {
		response := helper.ApiResponse("Register account failed", http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(newUser, token)

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

	token, err := handler.authService.GenerateToken(loggedInUser.ID)
	if err != nil {
		response := helper.ApiResponse("Login failed", http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(loggedInUser, token)

	response := helper.ApiResponse("Successfuly logged in", http.StatusOK, "Success", formatter)

	c.JSON(http.StatusOK, response)
}

func (handler *userHandler) CheckEmailAvailability(c *gin.Context) {
	/*
		(Check apakah email sudah terdaftar / belum ?)
		ada input email dari user
		input email di mapping ke struct input
		struct input di passing ke service
		service akan memanggil repository - email sudah ada atau belum
		repository - db
	*/
	var input user.CheckEmailInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		// Errors validations
		errors := helper.FormatValidationError(err)
		errorMesage := gin.H{"errors": errors}

		response := helper.ApiResponse("Email checking failed", http.StatusUnprocessableEntity, "failed", errorMesage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isEmailAvailable, err := handler.userService.IsEmailAvailable(input)
	if err != nil {
		errorMesage := gin.H{"errors": "Server error"}
		response := helper.ApiResponse("Email checking failed", http.StatusUnprocessableEntity, "failed", errorMesage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// key-val email status message
	data := gin.H{
		"is_available": isEmailAvailable,
	}

	// Meta messagse reference to bool-value
	var metaMessage string
	if isEmailAvailable {
		metaMessage = "Email is available"
	} else {
		metaMessage = "Email has already taken"
	}
	response := helper.ApiResponse(metaMessage, http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) UploadAvatar(c *gin.Context) {
	/*
	 input dari user dalam bentuk form-file
	 simpan gambarnya di folder "images/"
	 di service memangil repo
	 JWT (Sementara Hardcode, seakan-akan user yang login ID = 1)
	 repo ambil data user yg ID = 1
	 repo update data user simpan lokasi file
	*/
	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.ApiResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	// nantinya dapat dari JWT (masih hardcode)
	userID := 1

	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.ApiResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.userService.SaveAvatar(userID, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.ApiResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.ApiResponse("Avatar successfully uploaded", http.StatusOK, "Success", data)

	c.JSON(http.StatusOK, response)
}
