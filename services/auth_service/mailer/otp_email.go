package authmailer

import (
	"fmt"
	"geko/internal/helper"
)

type OtpEmailTemplateData struct {
	Email      string
	Otp        string
	AppName    string
	Expiration int
}

func (authMailer *AuthMailer) SendOTPEmail(templData OtpEmailTemplateData) {

	emailBody, err := helper.LoadHtmlTemplateToString(FS, "templates/otp.templ", templData)
	if err != nil {
		// @TODO implement error logging
		fmt.Println("Template parsing error:", err)
		return
	}

	// Send the email to user
	if err := (*authMailer.mailer).SendHTML(authMailer.managerEmail, []string{templData.Email}, "OTP", emailBody); err != nil {
		// @TODO implement error logging
		fmt.Println(err)
	}
}
