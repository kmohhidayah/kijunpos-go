package email

import (
	"context"
	"fmt"
	"github/kijunpos/internal/domain"
	"github/kijunpos/internal/pkg/apm"
	"net/smtp"
	"strings"
)

// Config holds the configuration for the email service
type Config struct {
	SMTPHost     string
	SMTPPort     string
	SenderEmail  string
	SenderName   string
	SMTPPassword string
}

// Service implements the domain.EmailService interface
type Service struct {
	config Config
}

// NewEmailService creates a new email service
func NewEmailService(config Config) domain.EmailService {
	return &Service{
		config: config,
	}
}

// SendVerificationCode sends a verification code to the specified email address
func (s *Service) SendVerificationCode(ctx context.Context, email, code string) error {
	ctx, span := apm.GetTracer().Start(ctx, "pkg.email.SendVerificationCode")
	defer span.End()

	// Set up authentication information
	auth := smtp.PlainAuth(
		"",
		s.config.SenderEmail,
		s.config.SMTPPassword,
		s.config.SMTPHost,
	)

	// Compose email
	subject := "KijunPOS - Your Verification Code"
	body := fmt.Sprintf(`
	<html>
		<body>
			<h2>KijunPOS Verification Code</h2>
			<p>Hello,</p>
			<p>Your verification code is: <strong>%s</strong></p>
			<p>This code will expire in 10 minutes.</p>
			<p>If you did not request this code, please ignore this email.</p>
			<p>Thank you,<br>KijunPOS Team</p>
		</body>
	</html>
	`, code)

	// Compose message
	to := []string{email}
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	message := []byte(fmt.Sprintf("To: %s\r\n"+
		"From: %s <%s>\r\n"+
		"Subject: %s\r\n"+
		"%s\r\n"+
		"%s", 
		strings.Join(to, ","), 
		s.config.SenderName,
		s.config.SenderEmail, 
		subject,
		mime,
		body))

	// Send email
	addr := fmt.Sprintf("%s:%s", s.config.SMTPHost, s.config.SMTPPort)
	err := smtp.SendMail(addr, auth, s.config.SenderEmail, to, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
