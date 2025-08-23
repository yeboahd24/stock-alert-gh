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

func (s *EmailService) SendIPOAlertEmail(user *models.User, alert *models.Alert, ipo *models.IPOAnnouncement, eventType string) error {
	if s.config.SMTPUser == "" || s.config.SMTPPassword == "" {
		return fmt.Errorf("email service not configured")
	}

	var subject string
	var body string
	var err error

	if eventType == "announced" {
		subject = fmt.Sprintf("New IPO Announced: %s (%s)", ipo.CompanyName, ipo.Symbol)
		body, err = s.generateIPOAnnouncementEmail(user.Name, ipo)
	} else {
		subject = fmt.Sprintf("IPO Now Listed: %s (%s)", ipo.CompanyName, ipo.Symbol)
		body, err = s.generateIPOListingEmail(user.Name, ipo)
	}

	if err != nil {
		return fmt.Errorf("failed to generate email body: %w", err)
	}

	return s.sendEmail(user.Email, subject, body)
}

func (s *EmailService) generateIPOAnnouncementEmail(userName string, ipo *models.IPOAnnouncement) (string, error) {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>New IPO Announced</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #059669; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background-color: #f9f9f9; }
        .ipo-box { background-color: #fff; border-left: 4px solid #059669; padding: 15px; margin: 20px 0; }
        .footer { text-align: center; padding: 20px; font-size: 12px; color: #666; }
        .price { font-size: 20px; font-weight: bold; color: #059669; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🎉 New IPO Announced!</h1>
        </div>
        <div class="content">
            <p>Hello %s,</p>
            
            <div class="ipo-box">
                <h3>%s (%s)</h3>
                <p><strong>Sector:</strong> %s</p>
                <p><strong>Offer Price:</strong> <span class="price">GH₵ %.2f</span></p>
                <p><strong>Expected Listing Date:</strong> %s</p>
            </div>
            
            <p>A new company is going public on the Ghana Stock Exchange! This is your opportunity to invest in %s before it starts trading.</p>
            
            <p>Keep an eye on the listing date and be ready to trade when it goes live.</p>
            
            <p>Best regards,<br>The Shares Alert Ghana Team</p>
        </div>
        <div class="footer">
            <p>This is an automated IPO alert from Shares Alert Ghana.</p>
        </div>
    </div>
</body>
</html>
`, userName, ipo.CompanyName, ipo.Symbol, ipo.Sector, ipo.OfferPrice, ipo.ListingDate.Format("January 2, 2006"), ipo.CompanyName), nil
}

func (s *EmailService) generateIPOListingEmail(userName string, ipo *models.IPOAnnouncement) (string, error) {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>IPO Now Listed</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #2563eb; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background-color: #f9f9f9; }
        .ipo-box { background-color: #fff; border-left: 4px solid #2563eb; padding: 15px; margin: 20px 0; }
        .footer { text-align: center; padding: 20px; font-size: 12px; color: #666; }
        .price { font-size: 20px; font-weight: bold; color: #2563eb; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>📈 IPO Now Trading!</h1>
        </div>
        <div class="content">
            <p>Hello %s,</p>
            
            <div class="ipo-box">
                <h3>%s (%s)</h3>
                <p><strong>Sector:</strong> %s</p>
                <p><strong>Original Offer Price:</strong> <span class="price">GH₵ %.2f</span></p>
                <p><strong>Listed On:</strong> %s</p>
            </div>
            
            <p>Great news! %s is now officially listed and trading on the Ghana Stock Exchange.</p>
            
            <p>You can now buy and sell shares of %s through your broker or trading platform.</p>
            
            <p>Best regards,<br>The Shares Alert Ghana Team</p>
        </div>
        <div class="footer">
            <p>This is an automated IPO listing notification from Shares Alert Ghana.</p>
        </div>
    </div>
</body>
</html>
`, userName, ipo.CompanyName, ipo.Symbol, ipo.Sector, ipo.OfferPrice, ipo.ListingDate.Format("January 2, 2006"), ipo.CompanyName, ipo.CompanyName), nil
}

func (s *EmailService) SendDividendAlertEmail(user *models.User, alert *models.Alert, dividend *models.DividendAnnouncement, eventType string) error {
	if s.config.SMTPUser == "" || s.config.SMTPPassword == "" {
		return fmt.Errorf("email service not configured")
	}

	var subject string
	var body string
	var err error

	if eventType == "announced" {
		subject = fmt.Sprintf("Dividend Announced: %s (%s)", dividend.StockName, dividend.StockSymbol)
		body, err = s.generateDividendAnnouncementEmail(user.Name, dividend)
	} else {
		subject = fmt.Sprintf("Dividend Paid: %s (%s)", dividend.StockName, dividend.StockSymbol)
		body, err = s.generateDividendPaymentEmail(user.Name, dividend)
	}

	if err != nil {
		return fmt.Errorf("failed to generate email body: %w", err)
	}

	return s.sendEmail(user.Email, subject, body)
}

func (s *EmailService) generateDividendAnnouncementEmail(userName string, dividend *models.DividendAnnouncement) (string, error) {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Dividend Announced</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #16a34a; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background-color: #f9f9f9; }
        .dividend-box { background-color: #fff; border-left: 4px solid #16a34a; padding: 15px; margin: 20px 0; }
        .footer { text-align: center; padding: 20px; font-size: 12px; color: #666; }
        .amount { font-size: 24px; font-weight: bold; color: #16a34a; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>💰 Dividend Announced!</h1>
        </div>
        <div class="content">
            <p>Hello %s,</p>
            
            <div class="dividend-box">
                <h3>%s (%s)</h3>
                <p><strong>Dividend Type:</strong> %s</p>
                <p><strong>Amount:</strong> <span class="amount">%s %.2f</span></p>
                <p><strong>Ex-Dividend Date:</strong> %s</p>
                <p><strong>Payment Date:</strong> %s</p>
            </div>
            
            <p>Great news! %s has announced a dividend payment. Make sure you own shares before the ex-dividend date to be eligible.</p>
            
            <p>Best regards,<br>The Shares Alert Ghana Team</p>
        </div>
        <div class="footer">
            <p>This is an automated dividend alert from Shares Alert Ghana.</p>
        </div>
    </div>
</body>
</html>
`, userName, dividend.StockName, dividend.StockSymbol, dividend.DividendType, dividend.Currency, dividend.Amount, dividend.ExDate.Format("January 2, 2006"), dividend.PaymentDate.Format("January 2, 2006"), dividend.StockName), nil
}

func (s *EmailService) generateDividendPaymentEmail(userName string, dividend *models.DividendAnnouncement) (string, error) {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Dividend Payment</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #059669; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background-color: #f9f9f9; }
        .dividend-box { background-color: #fff; border-left: 4px solid #059669; padding: 15px; margin: 20px 0; }
        .footer { text-align: center; padding: 20px; font-size: 12px; color: #666; }
        .amount { font-size: 24px; font-weight: bold; color: #059669; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>✅ Dividend Paid!</h1>
        </div>
        <div class="content">
            <p>Hello %s,</p>
            
            <div class="dividend-box">
                <h3>%s (%s)</h3>
                <p><strong>Dividend Type:</strong> %s</p>
                <p><strong>Amount:</strong> <span class="amount">%s %.2f</span></p>
                <p><strong>Payment Date:</strong> %s</p>
            </div>
            
            <p>The dividend for %s has been paid! If you were eligible, the payment should appear in your account soon.</p>
            
            <p>Best regards,<br>The Shares Alert Ghana Team</p>
        </div>
        <div class="footer">
            <p>This is an automated dividend payment notification from Shares Alert Ghana.</p>
        </div>
    </div>
</body>
</html>
`, userName, dividend.StockName, dividend.StockSymbol, dividend.DividendType, dividend.Currency, dividend.Amount, dividend.PaymentDate.Format("January 2, 2006"), dividend.StockName), nil
}

func (s *EmailService) SendDividendYieldAlertEmail(user *models.User, alert *models.Alert, currentYield float64) error {
	if s.config.SMTPUser == "" || s.config.SMTPPassword == "" {
		return fmt.Errorf("email service not configured")
	}

	var subject string
	var body string
	var err error

	switch alert.AlertType {
	case models.AlertTypeHighDividendYield:
		subject = fmt.Sprintf("High Dividend Yield Alert: %s (%s)", alert.StockName, alert.StockSymbol)
		body, err = s.generateHighDividendYieldEmail(user.Name, alert, currentYield)
	case models.AlertTypeTargetDividendYield:
		subject = fmt.Sprintf("Target Dividend Yield Reached: %s (%s)", alert.StockName, alert.StockSymbol)
		body, err = s.generateTargetDividendYieldEmail(user.Name, alert, currentYield)
	case models.AlertTypeDividendYieldChange:
		subject = fmt.Sprintf("Dividend Yield Change Alert: %s (%s)", alert.StockName, alert.StockSymbol)
		body, err = s.generateDividendYieldChangeEmail(user.Name, alert, currentYield)
	default:
		return fmt.Errorf("unsupported dividend yield alert type: %s", alert.AlertType)
	}

	if err != nil {
		return fmt.Errorf("failed to generate email body: %w", err)
	}

	return s.sendEmail(user.Email, subject, body)
}

func (s *EmailService) generateHighDividendYieldEmail(userName string, alert *models.Alert, currentYield float64) (string, error) {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>High Dividend Yield Alert</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #16a34a; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background-color: #f9f9f9; }
        .yield-box { background-color: #fff; border-left: 4px solid #16a34a; padding: 15px; margin: 20px 0; }
        .footer { text-align: center; padding: 20px; font-size: 12px; color: #666; }
        .yield { font-size: 28px; font-weight: bold; color: #16a34a; }
        .threshold { font-size: 18px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>📈 High Dividend Yield Alert!</h1>
        </div>
        <div class="content">
            <p>Hello %s,</p>
            
            <div class="yield-box">
                <h3>%s (%s)</h3>
                <p><strong>Current Dividend Yield:</strong> <span class="yield">%.2f%%</span></p>
                <p class="threshold">Your threshold: %.2f%%</p>
            </div>
            
            <p>Great news! %s has reached your high dividend yield threshold of %.2f%%. The current dividend yield is %.2f%%.</p>
            
            <p>This could be a good opportunity to consider this stock for dividend income investing.</p>
            
            <p>Best regards,<br>The Shares Alert Ghana Team</p>
        </div>
        <div class="footer">
            <p>This is an automated dividend yield alert from Shares Alert Ghana.</p>
        </div>
    </div>
</body>
</html>
`, userName, alert.StockName, alert.StockSymbol, currentYield, *alert.ThresholdYield, alert.StockName, *alert.ThresholdYield, currentYield), nil
}

func (s *EmailService) generateTargetDividendYieldEmail(userName string, alert *models.Alert, currentYield float64) (string, error) {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Target Dividend Yield Reached</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #2563eb; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background-color: #f9f9f9; }
        .yield-box { background-color: #fff; border-left: 4px solid #2563eb; padding: 15px; margin: 20px 0; }
        .footer { text-align: center; padding: 20px; font-size: 12px; color: #666; }
        .yield { font-size: 28px; font-weight: bold; color: #2563eb; }
        .target { font-size: 18px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🎯 Target Dividend Yield Reached!</h1>
        </div>
        <div class="content">
            <p>Hello %s,</p>
            
            <div class="yield-box">
                <h3>%s (%s)</h3>
                <p><strong>Current Dividend Yield:</strong> <span class="yield">%.2f%%</span></p>
                <p class="target">Your target: %.2f%%</p>
            </div>
            
            <p>Excellent! %s has reached your target dividend yield of %.2f%%. The current yield is %.2f%%.</p>
            
            <p>This might be the perfect time to consider your investment strategy for this stock.</p>
            
            <p>Best regards,<br>The Shares Alert Ghana Team</p>
        </div>
        <div class="footer">
            <p>This is an automated target dividend yield alert from Shares Alert Ghana.</p>
        </div>
    </div>
</body>
</html>
`, userName, alert.StockName, alert.StockSymbol, currentYield, *alert.TargetYield, alert.StockName, *alert.TargetYield, currentYield), nil
}

func (s *EmailService) generateDividendYieldChangeEmail(userName string, alert *models.Alert, currentYield float64) (string, error) {
	change := currentYield - *alert.LastYield
	changeDirection := "increased"
	if change < 0 {
		changeDirection = "decreased"
		change = -change
	}

	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Dividend Yield Change Alert</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #f59e0b; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background-color: #f9f9f9; }
        .yield-box { background-color: #fff; border-left: 4px solid #f59e0b; padding: 15px; margin: 20px 0; }
        .footer { text-align: center; padding: 20px; font-size: 12px; color: #666; }
        .yield { font-size: 24px; font-weight: bold; color: #f59e0b; }
        .change { font-size: 18px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>📊 Dividend Yield Change Alert!</h1>
        </div>
        <div class="content">
            <p>Hello %s,</p>
            
            <div class="yield-box">
                <h3>%s (%s)</h3>
                <p><strong>Current Dividend Yield:</strong> <span class="yield">%.2f%%</span></p>
                <p class="change">Previous yield: %.2f%%</p>
                <p class="change">Change: %.2f%% (%s)</p>
                <p class="change">Your threshold: %.2f%%</p>
            </div>
            
            <p>The dividend yield for %s has %s by %.2f%%, which exceeds your change threshold of %.2f%%.</p>
            
            <p>This significant change might warrant a review of your investment position in this stock.</p>
            
            <p>Best regards,<br>The Shares Alert Ghana Team</p>
        </div>
        <div class="footer">
            <p>This is an automated dividend yield change alert from Shares Alert Ghana.</p>
        </div>
    </div>
</body>
</html>
`, userName, alert.StockName, alert.StockSymbol, currentYield, *alert.LastYield, change, changeDirection, *alert.YieldChangeThreshold, alert.StockName, changeDirection, change, *alert.YieldChangeThreshold), nil
}