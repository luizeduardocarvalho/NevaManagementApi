package email

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"time"

	mailjet "github.com/mailjet/mailjet-apiv3-go/v4"
)

const invitationTemplate = `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<style>
		body {
			font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
			line-height: 1.6;
			color: #333;
			max-width: 600px;
			margin: 0 auto;
			padding: 20px;
		}
		.container {
			background-color: #f9f9f9;
			border-radius: 8px;
			padding: 30px;
			box-shadow: 0 2px 4px rgba(0,0,0,0.1);
		}
		h1 {
			color: #2563eb;
			margin-bottom: 20px;
		}
		.button {
			display: inline-block;
			background-color: #2563eb;
			color: white;
			padding: 12px 30px;
			text-decoration: none;
			border-radius: 6px;
			margin: 20px 0;
			font-weight: 600;
		}
		.details {
			background-color: #fff;
			padding: 15px;
			border-left: 4px solid #2563eb;
			margin: 20px 0;
		}
		.footer {
			font-size: 12px;
			color: #666;
			margin-top: 30px;
			padding-top: 20px;
			border-top: 1px solid #ddd;
		}
	</style>
</head>
<body>
	<div class="container">
		<h1>You're Invited to LabFlux!</h1>

		<p>Hi{{if .FirstName}} {{.FirstName}}{{end}},</p>

		<p>
			<strong>{{.InvitedByName}}</strong> has invited you to join
			<strong>{{.LaboratoryName}}</strong> on LabFlux as a
			<strong>{{.Role}}</strong>.
		</p>

		<a href="{{.InvitationLink}}" class="button">
			Accept Invitation
		</a>

		<div class="details">
			<p><strong>Laboratory:</strong> {{.LaboratoryName}}</p>
			<p><strong>Role:</strong> {{.Role}}</p>
			<p><strong>Invited by:</strong> {{.InvitedByName}}</p>
			<p><strong>Expires:</strong> {{.ExpirationDate}}</p>
		</div>

		<p>
			This invitation will expire on <strong>{{.ExpirationDate}}</strong>.
			Please accept it before then to join the laboratory.
		</p>

		<p class="footer">
			If you didn't expect this invitation, you can safely ignore this email.
			<br>
			This is an automated message from LabFlux.
		</p>
	</div>
</body>
</html>
`

type InvitationEmailData struct {
	FirstName      string
	LaboratoryName string
	Role           string
	InvitedByName  string
	InvitationLink string
	ExpirationDate string
}

type EmailService interface {
	SendInvitationEmail(to, firstName, laboratoryName, role, invitedByName, token string, expiresAt time.Time) error
}

type MailjetEmailService struct {
	client      *mailjet.Client
	fromEmail   string
	fromName    string
	frontendURL string
}

func NewMailjetEmailService() (*MailjetEmailService, error) {
	apiKey := os.Getenv("MAILJET_API_KEY")
	apiSecret := os.Getenv("MAILJET_API_SECRET")

	if apiKey == "" || apiSecret == "" {
		return nil, fmt.Errorf("MAILJET_API_KEY and MAILJET_API_SECRET environment variables must be set")
	}

	fromEmail := os.Getenv("MAILJET_FROM_EMAIL")
	if fromEmail == "" {
		fromEmail = "luiz@labflux.app"
	}

	fromName := os.Getenv("MAILJET_FROM_NAME")
	if fromName == "" {
		fromName = "LabFlux"
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5174"
	}

	client := mailjet.NewMailjetClient(apiKey, apiSecret)

	return &MailjetEmailService{
		client:      client,
		fromEmail:   fromEmail,
		fromName:    fromName,
		frontendURL: frontendURL,
	}, nil
}

func (s *MailjetEmailService) SendInvitationEmail(to, firstName, laboratoryName, role, invitedByName, token string, expiresAt time.Time) error {
	// Parse the email template
	tmpl, err := template.New("invitation").Parse(invitationTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	// Prepare template data
	data := InvitationEmailData{
		FirstName:      firstName,
		LaboratoryName: laboratoryName,
		Role:           role,
		InvitedByName:  invitedByName,
		InvitationLink: fmt.Sprintf("%s/accept-invitation?token=%s", s.frontendURL, token),
		ExpirationDate: expiresAt.Format("January 2, 2006"),
	}

	// Execute template
	var htmlBody bytes.Buffer
	if err := tmpl.Execute(&htmlBody, data); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	// Create email message
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: s.fromEmail,
				Name:  s.fromName,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: to,
					Name:  firstName,
				},
			},
			Subject:  fmt.Sprintf("You're invited to join %s on LabFlux", laboratoryName),
			HTMLPart: htmlBody.String(),
		},
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}

	// Send email
	res, err := s.client.SendMailV31(&messages)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	if res.ResultsV31 != nil && len(res.ResultsV31) > 0 {
		result := res.ResultsV31[0]
		if result.Status != "success" {
			return fmt.Errorf("mailjet API error: %+v", result)
		}
		log.Printf("Invitation email sent to %s (status: %s)", to, result.Status)
	}

	return nil
}

// MockEmailService for development/testing
type MockEmailService struct{}

func NewMockEmailService() *MockEmailService {
	return &MockEmailService{}
}

func (s *MockEmailService) SendInvitationEmail(to, firstName, laboratoryName, role, invitedByName, token string, expiresAt time.Time) error {
	log.Printf("[MOCK EMAIL] Invitation sent to %s for %s as %s in %s (token: %s)",
		to, firstName, role, laboratoryName, token)
	return nil
}

// GetEmailService returns the appropriate email service based on configuration
func GetEmailService() EmailService {
	apiKey := os.Getenv("MAILJET_API_KEY")
	apiSecret := os.Getenv("MAILJET_API_SECRET")

	if apiKey == "" || apiSecret == "" {
		log.Println("MAILJET_API_KEY or MAILJET_API_SECRET not set, using mock email service")
		return NewMockEmailService()
	}

	service, err := NewMailjetEmailService()
	if err != nil {
		log.Printf("Failed to create Mailjet service: %v, using mock", err)
		return NewMockEmailService()
	}

	return service
}
