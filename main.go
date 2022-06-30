package main

import (
	"fmt"
	"go-campaign-funding/auth"
	"go-campaign-funding/campaign"
	"go-campaign-funding/handler"
	"go-campaign-funding/helper"
	"go-campaign-funding/user"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:@tcp(127.0.0.1:3317)/campaign_startup?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	// Service TEST
	campaignRepository := campaign.NewRepository(db)
	campaignService := campaign.NewService(campaignRepository)
	// Get All Campaigns
	campaigns, _ := campaignService.FindCampaigns(0)
	fmt.Println("DEBUG MODE..")
	fmt.Println("DEBUG MODE..")
	fmt.Println("Service TESTS..")
	if len(campaigns) > 0 {
		fmt.Println("Total Campaign Data:", len(campaigns))
	} else {
		fmt.Println("Data not found")
	}

	authService := auth.NewService()
	userHandler := handler.NewUserHanlder(userService, authService)

	router := gin.Default()
	api := router.Group("api/v1")

	// routing path API
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)

	router.Run(":3000")
}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.ApiResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// Bearer - TokenByGenerate
		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		// Token Validation
		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.ApiResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		payload, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.ApiResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// Get userId from user's payload
		userId := int(payload["user_id"].(float64))
		user, err := userService.GetUserByID(userId) // <= searching userId into Db
		if err != nil {                              // failed
			response := helper.ApiResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// Success
		c.Set("currentUser", user)
	}

}
