package controllers

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/mail"
	"net/smtp"

	"github.com/Bulusideng/go-crm/models"
)

func SendEmail(to, title, body string) bool {
	cfg := models.GetConfig()

	from := cfg.Mailaddr
	pwd := cfg.Mailpwd
	smtpaddr := cfg.Smtpaddr

	msg := "From: " + "Admin" + "\n" +
		"To: " + to + "\n" +
		"Subject: " + title + "\n\n" +
		body

	err := smtp.SendMail(smtpaddr+":25",
		smtp.PlainAuth("", from, pwd, smtpaddr),
		from, []string{to}, []byte(msg))

	if err != nil {
		fmt.Printf("smtp error: %s\n", err)
		return false
	}
	fmt.Printf("Send mail success: %s %s\n", to, body)
	return true
}

func send1() {

	from := mail.Address{"", "xianjun.deng@gmail.com"}
	to := mail.Address{"", "xianjun.deng@gmail.com"}
	subj := "This is the email subject"
	body := "This is an example body.\n With two lines."

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the SMTP Server
	servername := "smtp.gmail.com:587"

	host, _, _ := net.SplitHostPort(servername)

	auth := smtp.PlainAuth("", "xianjun.deng@gmail.com", ".", host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		log.Panic(err)
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Panic(err)
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Panic(err)
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		log.Panic(err)
	}

	if err = c.Rcpt(to.Address); err != nil {
		log.Panic(err)
	}

	// Data
	w, err := c.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	c.Quit()

}
