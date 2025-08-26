package infrastructure



type ISMTPDialer interface{
	DialAndSend(...*gomail.Message) error
}

type SMTPService struct{
	dialer ISMTPDialer
	EmailFrom string

}

func NewSMTPService (SMTPHost string,SMTPPort int,SMTPUsername string,SMTPPassword string,EmailFrom string) domain.IEmailService{
	d:=gomail.NewDialer(SMTPHost,SMTPPort,SMTPUsername,SMTPPassword)
	return &SMTPService{
		dialer: d,
		EmailFrom: EmailFrom,
	}
}
func (s *SMTPService) SendEmail(to, subject, body string) error {

	m := gomail.NewMessage()
	m.SetHeader("From", s.EmailFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return s.dialer.DialAndSend(m)
}