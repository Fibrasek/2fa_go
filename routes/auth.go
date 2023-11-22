package routes

import (
	auth "github.com/fibrasek/2fa_go/controllers"
	"github.com/gin-gonic/gin"
)

type AuthRoute struct {
	authController auth.AuthController
}

func NewAuthRoute(authController auth.AuthController) AuthRoute {
	return AuthRoute{authController}
}

func (route *AuthRoute) AuthRouter(group *gin.RouterGroup) {
	router := group.Group("auth")

	// Sign up/in
	router.POST("/register", route.authController.SignUpUser)
	router.POST("/login", route.authController.LoginUser)

	// OTP
	router.POST("/otp/generate", route.authController.GenerateOTP)
	router.POST("/otp/verify", route.authController.VerifyOTP)
	router.POST("/otp/validate", route.authController.ValidateOTP)
	router.POST("/otp/disable", route.authController.DisableOTP)
}
