package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/millroy094/task-processor/pkg/task"
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
	smtpPort := os.Getenv("MAILHOG_PORT")
	from := os.Getenv("MAIL_FROM")
	password := ""

	to := []string{data.Email}
	subject := "Subject: " + data.Subject
	body := "Message: " + data.Body

	msg := []byte(subject + "\r\n\r\n" + body)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	if err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg); err != nil {
		log.Printf("Failed to send email to %s: %v\n", data.Email, err)
		return err
	}

	fmt.Printf("Email sent to %s\n", data.Email)
	return nil
}
