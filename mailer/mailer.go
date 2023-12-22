package mailer

import (
	"bytes"
	"embed"
	"text/template"

	"github.com/go-gomail/gomail"
)

//go:embed "templates/*"
var templateFS embed.FS

type Mailer struct {
	dialer *gomail.Dialer
	sender string
}

func New(host string, port int, username string, password string, sender string) Mailer {
	// Initialize a new mail.Dialer instance with the given SMTP server settings.
	dialer := gomail.NewDialer(host, port, username, password)

	// Return a Mailer instance containing the dialer and sender information.
	return Mailer{
		dialer: dialer,
		sender: sender,
	}
}

// Define a Send() method on the Mailer type.
func (mailer Mailer) Send(recipient string, templateFile string, data interface{}) error {
	// Use the ParseFS() method to parse the required template file from the embedded
	// file system.
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	// Execute the named template "subject", passing in the dynamic data and storing the
	// result in a bytes.Buffer variable.
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	// Use the mail.NewMessage() function to initialize a new mail.Message instance.
	message := gomail.NewMessage()
	message.SetHeader("To", recipient)
	message.SetHeader("From", mailer.sender)
	message.SetHeader("Subject", subject.String())
	message.SetBody("text/plain", plainBody.String())
	message.AddAlternative("text/html", htmlBody.String())

	// Call the DialAndSend() method on the dialer, passing in the message to send. This
	// opens a connection to the SMTP server, sends the message, then closes the
	// connection.
	err = mailer.dialer.DialAndSend(message)
	if err != nil {
		return err
	}

	return nil
}
