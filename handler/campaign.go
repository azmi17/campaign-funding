package handler

import (
	"go-campaign-funding/campaign"
	"go-campaign-funding/helper"
	"go-campaign-funding/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// tangkap parameter di handler
// handler ke service
// di service menentukan: repository mana yang di CALL (Get all data atau Get Spesifik data berdasarkan user tertentu)
// Repository: FindAll, FindByUserID
// Db

type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

// api/v1/campaigns
func (handler *campaignHandler) GetCampaigns(c *gin.Context) {

	// convert str to int (w ignoring an err)
	userID, _ := strconv.Atoi(c.Query("user_id"))

	var campaigns []campaign.Campaign
	campaigns, err := handler.service.GetCampaigns(userID)
	if err != nil {
		response := helper.ApiResponse("Error to get campaigns", http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.ApiResponse("List of campaigns", http.StatusOK, "success", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)
}

func (handler *campaignHandler) GetCampaign(c *gin.Context) {
	// api/v1/campaigns/1
	// handler: mapping id yang di url ke struct input => service, call formatter
	// service: struct inputt => menangkap id di url, memanggil repo
	// repository: get campaign by Id

	var input campaign.GetCampaignDetailInput
	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.ApiResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	campaignDetail, err := handler.service.GetCampaignByID(input)
	if err != nil {
		response := helper.ApiResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.ApiResponse("Campaign detail", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)
}

func (handler *campaignHandler) CreateCampaign(c *gin.Context) {

	var input campaign.CreateCampaignInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		// Errors validations
		errors := helper.FormatValidationError(err)
		errorMesage := gin.H{"errors": errors}

		response := helper.ApiResponse("Failed to create campaign", http.StatusUnprocessableEntity, "failed", errorMesage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// Get current User ID
	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	newCampaign, err := handler.service.CreateCampaign(input)
	if err != nil {
		response := helper.ApiResponse("Failed to create campaign", http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.ApiResponse("Success to create campaign", http.StatusOK, "success", campaign.FormatCampaign(newCampaign))
	c.JSON(http.StatusOK, response)

}
