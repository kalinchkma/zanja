package authcontroller

import (
	authstore "geko/internal/store/auth_store"
	authmailer "geko/services/auth_service/mailer"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RegisterPayload struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (s *AuthController) Register(ctx *gin.Context) {
	var registerBody RegisterPayload

	if err := ctx.ShouldBindJSON(&registerBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errors": err.Error(),
		})
		return
	}

	userStore := s.serverContext.Store.UserStore
	// Find the user by email, if already exist
	// Return error if user already  exist
	if _, err := userStore.FindByEmail(registerBody.Email); err == nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"errors": "User already exists",
		})
		return
	}

	// Hash password
	hashedPassword, err := userStore.HashPassword(registerBody.Password)
	// Check error
	if err != nil {
		// Return error
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"errors": err.Error(),
		})
		return
	}

	// Create new user
	user := authstore.User{
		Name:     registerBody.Name,
		Email:    registerBody.Email,
		Password: hashedPassword,
	}

	// Store user to database
	err = userStore.Create(user)
	// Check error
	if err != nil {
		// Return error
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"errors": err.Error(),
		})
		return
	}

	// Get Update user
	user, err = userStore.FindByEmail(user.Email)
	if err != nil {
		// If created user not found return error
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"errors": err.Error(),
		})
		return
	}

	// generate and store otp
	otpStore := s.serverContext.Store.OTPStore

	// New otp code
	newOTPCode := otpStore.GenerateOTP(6)

	// New otp object
	otp := authstore.OTP{
		Code:      newOTPCode,
		UserId:    user.ID,
		ExpiresAt: time.Now().Add(time.Duration(s.serverContext.Config.OTPValidationTime) * time.Minute),
	}

	// Store to database
	err = otpStore.Create(otp)

	// Only send the email if otp has been created on database
	if err == nil {
		templData := authmailer.OtpEmailTemplateData{
			Email:      user.Email,
			Otp:        otp.Code,
			AppName:    s.serverContext.Config.AppName,
			Expiration: s.serverContext.Config.OTPValidationTime,
		}

		s.mailer.SendOTPEmail(templData)
	}

	// Send
	ctx.SecureJSON(http.StatusCreated, gin.H{
		"message": "Register successfully",
		"data":    s.serverContext.Store.UserStore.Normalize(user),
	})

}
