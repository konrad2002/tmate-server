package service

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/konrad2002/tmate-server/dto"
	"github.com/konrad2002/tmate-server/misc"
	"github.com/konrad2002/tmate-server/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"html/template"
	"net/smtp"
)

type EmailService struct {
	configService  ConfigService
	memberService  MemberService
	historyService HistoryService
}

func NewEmailService(cs ConfigService, ms MemberService, hs HistoryService) EmailService {
	return EmailService{
		configService:  cs,
		memberService:  ms,
		historyService: hs,
	}
}

func (ems *EmailService) GetEmailSenders() (*[]dto.EmailSenderDto, error) {
	senders, err := ems.configService.GetMailConfigs()
	if err != nil {
		return nil, err
	}

	var senderDtos []dto.EmailSenderDto
	for _, config := range *senders {
		senderDtos = append(senderDtos, dto.EmailSenderDto{
			Address: config.Address,
			Name:    config.Name,
		})
	}

	return &senderDtos, nil
}

func (ems *EmailService) sendEmail(sender string, receiver string, subject string, content *bytes.Buffer, member model.Member) error {
	fmt.Printf("try to send mail '%s' from '%s' to '%s'\n", subject, sender, receiver)
	fmt.Printf("Body: %s\n", content)

	mailConfig, err := ems.configService.GetMailConfig(sender)
	if err != nil {
		return err
	}

	allowedMails := []string{"konrad@schwimmteamerzgebirge.de", "konrad2002@arcor.de", "dr.weiss@arcor.de", "johann2005@arcor.de", "ubuntovka@gmail.com"}
	//allowedMails := []string{"konrad@schwimmteamerzgebirge.de"}

	if !misc.Contains(allowedMails, receiver) {
		return errors.New("illegal receiver in testing")
	}

	// SMTP Server details
	smtpServer := mailConfig.Smtp.Host
	smtpPort := mailConfig.Smtp.Port
	username := mailConfig.Smtp.Username
	password := mailConfig.Smtp.Password

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
	from := mailConfig.Address
	to := receiver

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
	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: %s \n%s\n\n", subject, mimeHeaders)))

	body.Write(content.Bytes())

	_, err = wc.Write(body.Bytes())
	if err != nil {
		fmt.Println("Write Error:", err)
		return err
	}
	wc.Close()

	// Quit the session
	client.Quit()

	fmt.Printf("Email sent successfully to %s!", receiver)

	ems.historyService.LogEMailAction(primitive.NilObjectID, member.Identifier, body.String())

	return nil
}

func (ems *EmailService) SendEmailFromTemplate(sender string, receivers []primitive.ObjectID, subject string, tmpl string) error {
	t, _ := template.New("template").Parse(tmpl)

	specialFields, err := ems.configService.GetSpecialFields()
	if err != nil {
		fmt.Println("Failed to load special fields in SendEmailFromTemplate:", err)
		return err
	}

	var errs []error

	for _, receiver := range receivers {
		member, err := ems.memberService.GetById(receiver)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		email := member.Data[specialFields.EMail].(string)

		if email == "" {
			errs = append(errs, errors.New(fmt.Sprintf("failed to load email address for user: %s\n", receiver)))
		}

		var body bytes.Buffer
		// TODO: preprocess data to parse time strings correctly
		err = t.Execute(&body, member.Data)
		if err != nil {
			errs = append(errs, errors.New(fmt.Sprintf("failed to parse template user: %s; %s\n", receiver, err.Error())))
		}

		err = ems.sendEmail(sender, email, subject, &body, member)
		if err != nil {
			errs = append(errs, errors.New(fmt.Sprintf("failed to send mail for user: %s; %s\n", receiver, err.Error())))
		}
	}

	for _, err2 := range errs {
		println(err2.Error())
	}

	if len(errs) > 0 {
		return errors.New(fmt.Sprintf("failed to send %d emails!", len(errs)))
	}
	return nil
}

func (ems *EmailService) SendAttestEmail(firstName string, lastName string, email string, date string, subject string, tmpl string, member model.Member) error {
	// BODY

	t, _ := template.ParseFiles(tmpl)

	var body bytes.Buffer
	err := t.Execute(&body, struct {
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

	err = ems.sendEmail("attest@schwimmteamerzgebirge.de", email, subject, &body, member)
	if err != nil {
		return err
	}

	return nil
}
