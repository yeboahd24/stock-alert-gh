package services

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

	"shares-alert-backend/internal/config"
	"shares-alert-backend/internal/models"
)

type EmailService struct {
	config *config.EmailConfig
}

type AlertEmailData struct {
	UserName       string
	StockSymbol    string
	StockName      string
	CurrentPrice   float64
	ThresholdPrice float64
	AlertType      string
}

func NewEmailService(cfg *config.EmailConfig) *EmailService {
	return &EmailService{
		config: cfg,
	}
}

func (s *EmailService) SendAlertEmail(user *models.User, alert *models.Alert) error {
	if s.config.SMTPUser == "" || s.config.SMTPPassword == "" {
		return fmt.Errorf("email service not configured")
	}

	data := AlertEmailData{
		UserName:     user.Name,
		StockSymbol:  alert.StockSymbol,
		StockName:    alert.StockName,
		AlertType:    alert.AlertType,
	}

	if alert.CurrentPrice != nil {
		data.CurrentPrice = *alert.CurrentPrice
	}
	if alert.ThresholdPrice != nil {
		data.ThresholdPrice = *alert.ThresholdPrice
	}

	subject := fmt.Sprintf("Stock Alert: %s (%s)", alert.StockName, alert.StockSymbol)
	body, err := s.generateEmailBody(data)
	if err != nil {
		return fmt.Errorf("failed to generate email body: %w", err)
	}

	return s.sendEmail(user.Email, subject, body)
}

func (s *EmailService) SendWelcomeEmail(user *models.User) error {
	if s.config.SMTPUser == "" || s.config.SMTPPassword == "" {
		return fmt.Errorf("email service not configured")
	}

	subject := "Welcome to Shares Alert Ghana!"
	body := s.generateWelcomeEmailBody(user.Name)

	return s.sendEmail(user.Email, subject, body)
}

func (s *EmailService) sendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.config.SMTPUser, s.config.SMTPPassword, s.config.SMTPHost)

	msg := []byte(fmt.Sprintf(
		"From: %s <%s>\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n"+
			"\r\n"+
			"%s\r\n",
		s.config.FromName, s.config.FromEmail, to, subject, body))

	addr := s.config.SMTPHost + ":" + s.config.SMTPPort
	return smtp.SendMail(addr, auth, s.config.FromEmail, []string{to}, msg)
}

func (s *EmailService) generateEmailBody(data AlertEmailData) (string, error) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Stock Alert</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #2563eb; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background-color: #f9f9f9; }
        .alert-box { background-color: #fff; border-left: 4px solid #2563eb; padding: 15px; margin: 20px 0; }
        .footer { text-align: center; padding: 20px; font-size: 12px; color: #666; }
        .price { font-size: 24px; font-weight: bold; color: #059669; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Stock Alert Triggered!</h1>
        </div>
        <div class="content">
            <p>Hello {{.UserName}},</p>
            
            <div class="alert-box">
                <h3>{{.StockName}} ({{.StockSymbol}})</h3>
                {{if eq .AlertType "price_threshold"}}
                    <p>Your price threshold alert has been triggered!</p>
                    <p>Current Price: <span class="price">GH₵ {{printf "%.2f" .CurrentPrice}}</span></p>
                    <p>Your Threshold: GH₵ {{printf "%.2f" .ThresholdPrice}}</p>
                {{else if eq .AlertType "dividend_announcement"}}
                    <p>A dividend has been announced for {{.StockName}}!</p>
                {{else if eq .AlertType "ipo_alert"}}
                    <p>IPO alert for {{.StockName}} has been triggered!</p>
                {{end}}
            </div>
            
            <p>You can view more details and manage your alerts by logging into your dashboard.</p>
            
            <p>Best regards,<br>The Shares Alert Ghana Team</p>
        </div>
        <div class="footer">
            <p>This is an automated message from Shares Alert Ghana. Please do not reply to this email.</p>
        </div>
    </div>
</body>
</html>
`

	t, err := template.New("alert").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (s *EmailService) generateWelcomeEmailBody(userName string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Welcome to Shares Alert Ghana</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #2563eb; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background-color: #f9f9f9; }
        .footer { text-align: center; padding: 20px; font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome to Shares Alert Ghana!</h1>
        </div>
        <div class="content">
            <p>Hello %s,</p>
            
            <p>Welcome to Shares Alert Ghana! We're excited to have you on board.</p>
            
            <p>With our platform, you can:</p>
            <ul>
                <li>Track Ghana Stock Exchange prices in real-time</li>
                <li>Set up custom price alerts for your favorite stocks</li>
                <li>Receive notifications when your alerts are triggered</li>
                <li>Stay informed about dividend announcements and IPOs</li>
            </ul>
            
            <p>Start by setting up your first stock alert and never miss an important price movement again!</p>
            
            <p>Best regards,<br>The Shares Alert Ghana Team</p>
        </div>
        <div class="footer">
            <p>This is an automated message from Shares Alert Ghana.</p>
        </div>
    </div>
</body>
</html>
`, userName)
}