package service

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
	"text/template"
)

type EmailConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	From     string
	To       string
}

type EmailService struct {
	config *EmailConfig
}

func NewEmailService(config *EmailConfig) *EmailService {
	if config.From == "" {
		config.From = config.User
	}
	return &EmailService{config: config}
}

// Send 发送邮件
func (s *EmailService) Send(to, subject, body string) error {
	if s.config.Host == "" || s.config.User == "" || s.config.Password == "" {
		return fmt.Errorf("邮件配置不完整")
	}

	auth := smtp.PlainAuth("", s.config.User, s.config.Password, s.config.Host)

	msg := fmt.Sprintf("From: %s\r\n", s.config.From)
	msg += fmt.Sprintf("To: %s\r\n", to)
	msg += fmt.Sprintf("Subject: %s\r\n", subject)
	msg += "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"
	msg += body

	addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)

	// 使用 TLS
	if s.config.Port == "465" || s.config.Port == "587" {
		return s.sendWithTLS(addr, auth, s.config.From, []string{to}, []byte(msg))
	}

	return smtp.SendMail(addr, auth, s.config.From, []string{to}, []byte(msg))
}

func (s *EmailService) sendWithTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	host := strings.Split(addr, ":")[0]

	conn, err := tls.Dial("tcp", addr, &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	})
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer client.Close()

	if err = client.Auth(auth); err != nil {
		return err
	}

	if err = client.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = w.Write(msg)
	return err
}

// SendWithTemplate 使用模板发送邮件
func (s *EmailService) SendWithTemplate(to, subject, templateName string, data interface{}) error {
	tmpl, err := template.New("email").Parse(getEmailTemplate(templateName))
	if err != nil {
		return err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	return s.Send(to, subject, body.String())
}

// 发送通知邮件
func (s *EmailService) SendNotification(subject, content string) error {
	if s.config.To == "" {
		return fmt.Errorf("未配置收件人")
	}
	return s.Send(s.config.To, subject, content)
}

// 获取邮件模板
func getEmailTemplate(name string) string {
	templates := map[string]string{
		"message": `
<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family: Arial, sans-serif; padding: 20px; background: #f5f5f5;">
<div style="max-width: 600px; margin: 0 auto; background: #fff; padding: 30px; border-radius: 8px;">
	<h2 style="color: #333; margin-bottom: 20px;">您收到一条新留言</h2>
	<div style="background: #f9f9f9; padding: 15px; border-radius: 4px; margin-bottom: 20px;">
		<p><strong>姓名：</strong>{{.Name}}</p>
		<p><strong>邮箱：</strong>{{.Email}}</p>
		<p><strong>内容：</strong></p>
		<p style="white-space: pre-wrap;">{{.Content}}</p>
	</div>
	<a href="{{.AdminURL}}" style="display: inline-block; padding: 10px 20px; background: #6366f1; color: #fff; text-decoration: none; border-radius: 4px;">前往管理后台</a>
</div>
</body>
</html>`,
		"article": `
<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family: Arial, sans-serif; padding: 20px; background: #f5f5f5;">
<div style="max-width: 600px; margin: 0 auto; background: #fff; padding: 30px; border-radius: 8px;">
	<h2 style="color: #333; margin-bottom: 20px;">{{.Title}}</h2>
	<div style="color: #666; line-height: 1.8;">
		{{.Content}}
	</div>
	<a href="{{.URL}}" style="display: inline-block; padding: 10px 20px; background: #6366f1; color: #fff; text-decoration: none; border-radius: 4px; margin-top: 20px;">阅读全文</a>
</div>
</body>
</html>`,
	}
	return templates[name]
}
