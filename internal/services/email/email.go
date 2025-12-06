package email

import (
	"fmt"
	"main/internal/pkg/config"
	"math/rand"
	"net/smtp"
)

type EmailService struct{}

func NewEmailService() *EmailService { return &EmailService{} }

func (es EmailService) SendMailSimple(subject, body string, to []string) error {
	auth := smtp.PlainAuth(
		"",
		config.GetConfig().SenderEmail,
		config.GetConfig().AppPassword,
		"smtp.gmail.com",
	)
	msg := "Subject" + subject + "\n" + body
	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"abduslamov.adhmabek@gmail.com",
		to,
		[]byte(msg),
	)
	if err != nil {
		return err
	}
	return nil
}

func (es EmailService) GenerateCode(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "000000"
	}

	code := ""
	for i := 0; i < n; i++ {
		code += fmt.Sprintf("%d", int(b[i]%10))
	}
	return code
}
