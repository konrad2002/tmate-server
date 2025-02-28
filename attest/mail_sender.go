package attest

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
)

func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func SendAttestEmail(firstName string, lastName string, email string, date string, subject string, tmpl string) error {

	//allowedMails := []string{"konrad@schwimmteamerzgebirge.de", "konrad2002@arcor.de", "dr.weiss@arcor.de", "johann2005@arcor.de", "ubuntovka@gmail.com"}
	allowedMails := []string{"konrad@schwimmteamerzgebirge.de"}

	if !contains(allowedMails, email) {
		return errors.New("illegal receiver in testing")
	}

	// SMTP Server details
	smtpServer := "smtp.strato.de"
	smtpPort := "465" // or "587"
	username := "attest@schwimmteamerzgebirge.de"
	password := os.Getenv("TMATE_SMTP_PASSWORD")

	// Setup authentication
	auth := smtp.PlainAuth("", username, password, smtpServer)

	// TLS Config for SSL (Port 465)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpServer,
	}

	// Connect to SMTP Server
	conn, err := tls.Dial("tcp", smtpServer+":"+smtpPort, tlsConfig)
	if err != nil {
		fmt.Println("TLS Connection Error:", err)
		return err
	}

	// Create SMTP client
	client, err := smtp.NewClient(conn, smtpServer)
	if err != nil {
		fmt.Println("SMTP Client Error:", err)
		return err
	}

	// Authenticate
	if err := client.Auth(auth); err != nil {
		fmt.Println("Authentication Error:", err)
		return err
	}

	// Set sender and recipient
	from := username
	to := email

	if err := client.Mail(from); err != nil {
		fmt.Println("Mail Error:", err)
		return err
	}
	if err := client.Rcpt(to); err != nil {
		fmt.Println("Recipient Error:", err)
		return err
	}

	// Write message
	wc, err := client.Data()
	if err != nil {
		fmt.Println("Data Error:", err)
		return err
	}

	// BODY

	t, _ := template.ParseFiles(tmpl)

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: %s \n%s\n\n", subject, mimeHeaders)))

	err = t.Execute(&body, struct {
		FirstName string
		LastName  string
		Date      string
	}{
		FirstName: firstName,
		LastName:  lastName,
		Date:      date,
	})
	if err != nil {
		fmt.Println("template Error:", err)
		return err
	}

	_, err = wc.Write(body.Bytes())
	if err != nil {
		fmt.Println("Write Error:", err)
		return err
	}
	wc.Close()

	// Quit the session
	client.Quit()

	fmt.Println("Email sent successfully!")

	return nil
}
