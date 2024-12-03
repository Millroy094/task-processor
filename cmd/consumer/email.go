package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/millroy094/task-processor/pkg/task"
	mail "github.com/xhit/go-simple-mail/v2"
)

type EmailPayload struct {
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func sendEmail(task task.Task) error {
	var data EmailPayload

	if err := json.Unmarshal([]byte(task.Payload), &data); err != nil {
		log.Printf("Failed to unmarshal email payload: %v", err)
		return err
	}

	smtpHost := os.Getenv("MAILHOG_HOST")
	smtpPortStr := os.Getenv("MAILHOG_PORT")
	from := os.Getenv("MAIL_FROM")

	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		log.Printf("Failed to convert smtpPort to int: %v", err)
		return err
	}

	smtpClient := mail.NewSMTPClient()

	smtpClient.Host = smtpHost
	smtpClient.Port = smtpPort

	smtpClient.Username = ""
	smtpClient.Password = ""
	smtpClient.Encryption = mail.EncryptionNone

	client, err := smtpClient.Connect()
	if err != nil {
		log.Printf("Failed to connect to SMTP server: %v", err)
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(from).
		AddTo(data.Email).
		SetSubject(data.Subject).
		SetBody(mail.TextPlain, data.Body)

	if err := email.Send(client); err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	fmt.Printf("Email sent to %s\n", data.Email)
	return nil
}
