package auth

import (
	"net/http"
	"os"
	"strings"

	"github.com/fibrasek/2fa_go/models"
	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(DB *gorm.DB) AuthController {
	return AuthController{DB}
}

func (conn *AuthController) SignUpUser(ctx *gin.Context) {
	var payload *models.RegisterUserInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err})
		return
	}

	newUser := models.User{
		Name:     payload.Name,
		Email:    strings.ToLower(payload.Email),
		Password: payload.Password,
	}

	result := conn.DB.Create(&newUser)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "E-mail already in use"})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "ok", "message": "Registered with success"})
}

func (conn *AuthController) LoginUser(ctx *gin.Context) {
	var payload *models.LoginUserInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
	}

	var user models.User
	result := conn.DB.First(&user, "email = ?", strings.ToLower(payload.Email))
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid credentials"})
		return
	}

	userResponse := gin.H{
		"id":          user.ID.String(),
		"name":        user.Name,
		"email":       user.Email,
		"otp_enabled": user.OTPEnabled,
	}

	ctx.JSON(http.StatusOK, gin.H{"stauts": "ok", "user": userResponse})
}

func (conn *AuthController) GenerateOTP(ctx *gin.Context) {
	var payload *models.OTPInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      os.Getenv("OTP_ISSUER"),
		AccountName: os.Getenv("OTP_ACCOUNT_NAME"),
		SecretSize:  15,
	})

	if err != nil {
		panic(err)
	}

	var user models.User
	result := conn.DB.First(&user, "id = ?", payload.UserId)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid credentials"})
		return
	}

	toUpdate := models.User{
		OTPSecret:  key.Secret(),
		OTPAuthUrl: key.URL(),
	}

	conn.DB.Model(&user).Updates(toUpdate)

	response := gin.H{
		"base32":      key.Secret(),
		"otpauth_url": key.URL(),
	}
	ctx.JSON(http.StatusOK, response)
}

func (conn *AuthController) VerifyOTP(ctx *gin.Context) {
	var payload *models.OTPInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	message := "Invalid token"

	var user models.User
	result := conn.DB.First(&user, "id = ?", payload.UserId)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": message})
		return
	}

	toUpdate := models.User{
		OTPEnabled:  true,
		OTPVerified: true,
	}

	conn.DB.Model(&user).Updates(toUpdate)

	response := gin.H{
		"id":          user.ID.String(),
		"name":        user.Name,
		"email":       user.Email,
		"otp_enabled": user.OTPEnabled,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "otp_verified": true, "user": response})
}

func (conn *AuthController) ValidateOTP(ctx *gin.Context) {
	var payload *models.OTPInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	message := "Invalid token"

	var user models.User
	result := conn.DB.First(&user, "id = ?", payload.UserId)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": message})
		return
	}

	valid := totp.Validate(payload.Token, user.OTPSecret)
	if !valid {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": message})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "otp_valid": true})
}

func (conn *AuthController) DisableOTP(ctx *gin.Context) {
	var payload *models.OTPInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	message := "Invalid token"

	var user models.User
	result := conn.DB.First(&user, "id = ?", payload.UserId)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": message})
		return
	}

	user.OTPEnabled = false
	conn.DB.Save(&user)

	response := gin.H{
		"id":          user.ID.String(),
		"name":        user.Name,
		"email":       user.Email,
		"otp_enabled": user.OTPEnabled,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "otp_disabled": true, "user": response})
}
