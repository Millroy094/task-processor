package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"

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

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	from := "your-email@gmail.com"
	password := "your-app-password"

	to := []string{data.Email}
	subject := "Subject: " + data.Subject
	body := "Message: " + data.Body

	msg := []byte(subject + "\r\n\r\n" + body)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	if err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg); err != nil {
		log.Printf("Failed to send email to %s: %v\n", data.Email, err)
		return err
	}

	// Log success
	fmt.Printf("Email sent to %s\n", data.Email)
	return nil

}
