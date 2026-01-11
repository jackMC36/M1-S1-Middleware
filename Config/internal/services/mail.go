package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	// "os"
)

const mailAPIURL = "https://mail-api.edu.forestier.re"

type MailRequest struct {
	Recipient string `json:"recipient"`
	Subject   string `json:"subject"`
	Content   string `json:"content"`
}

func SendMail(recipient, subject, content string) error {
	token := "tEUBejkaAwKsBzxVTyaJkvNcUMLJbnKwqWIePukb"	//os.Getenv("MAIL_API_TOKEN")
	if token == "" {
		return fmt.Errorf("MAIL_API_TOKEN environment variable not set")
	}

	mailReq := MailRequest{
		Recipient: recipient,
		Subject:   subject,
		Content:   content,
	}

	jsonData, err := json.Marshal(mailReq)
	if err != nil {
		return fmt.Errorf("failed to marshal mail request: %w", err)
	}

	req, err := http.NewRequest("POST", mailAPIURL+"/mail", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("mail API returned status %d", resp.StatusCode)
	}

	return nil
}
