package pkg

import (
	"fmt"
	"strconv"

	"gopkg.in/gomail.v2"
)

type Mailer interface {
	SendResetPasswordOTP(to string, otpCode string) error
}

type gomailMailer struct {
	host     string
	port     int
	username string
	password string
	from     string
}

func NewGomailMailer(host, portStr, username, password, from string) Mailer {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 587
	}
	return &gomailMailer{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (m *gomailMailer) SendResetPasswordOTP(to string, otpCode string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", "E-Wallet - Password Reset Request")

	body := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #e5e7eb; border-radius: 8px;">
			<h2 style="color: #6366f1; text-align: center;">Reset Password Request</h2>
			<p style="color: #4b5563; font-size: 16px;">Hello,</p>
			<p style="color: #4b5563; font-size: 16px;">We received a request to reset the password for your E-Wallet account. Please use the following OTP code to proceed:</p>
			
			<div style="background-color: #e0e7ff; padding: 16px; border-radius: 8px; text-align: center; margin: 24px 0;">
				<h1 style="color: #4338ca; letter-spacing: 8px; margin: 0; font-size: 32px;">%s</h1>
			</div>
			
			<p style="color: #4b5563; font-size: 14px;">This code is valid for <strong>5 minutes</strong>. If you did not request a password reset, please ignore this email or contact support immediately.</p>
			<hr style="border: none; border-top: 1px solid #e5e7eb; margin: 24px 0;" />
			<p style="color: #9ca3af; font-size: 12px; text-align: center;">This is an automated email from E-Wallet, please do not reply.</p>
		</div>`, otpCode)

	msg.SetBody("text/html", body)

	dialer := gomail.NewDialer(m.host, m.port, m.username, m.password)
	if err := dialer.DialAndSend(msg); err != nil {
		return fmt.Errorf("gomail failed send reset password email: %w", err)
	}

	return nil
}
