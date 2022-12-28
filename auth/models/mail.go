package models

import mail "github.com/xhit/go-simple-mail/v2"

// Model for a mail
type Mail struct {
	From        string
	To          string
	CC          string
	Subject     string
	Body        string
	Attachment  string
	Credentials MailCredentials
	Server      MailServer
}

// Model for the mail credentials
type MailCredentials struct {
	Username string
	Password string
}

// Model for the mail server
type MailServer struct {
	Host       string
	Port       int
	Encryption mail.Encryption
}
