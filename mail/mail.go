package mail

import (
	"github.com/russross/blackfriday"
	"github.com/xuqingfeng/mailman/account"
	"github.com/xuqingfeng/mailman/smtp"
	"github.com/xuqingfeng/mailman/util"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"strings"
)

type Mail struct {
	Subject string   `json:"subject"`
	To      []string `json:"to"`
	//Cc      []Contacts
	Cc   []string `json:"cc"`
	From string   `json:"from"`
	Body string   `json:"body"`
}

//****Mail START****
func SendMail(mail Mail) error {

	account, err := account.GetAccountInfo(mail.From)
	if err != nil {
		return err
	}
	SMTPServer, err := smtp.GetSMTPServer(mail.From)
	if err != nil {
		return err
	}
	m := gomail.NewMessage()
	m.SetHeader("Subject", mail.Subject)
	m.SetHeader("To", mail.To...)
	// with name - SetAddressHeader
	//m.SetAddressHeader("Cc", mail.Cc[0].Email, mail.Cc[0].Name)
	// multiple cc
	m.SetHeader("Cc", mail.Cc...)
	m.SetHeader("From", account.Email)

	content := ParseMailContent(mail.Body)

	m.SetBody("text/html", content)

	d := gomail.NewDialer(SMTPServer, util.DefaultSMTPPort, account.Email, account.Password)
	if err = d.DialAndSend(m); err != nil {
		util.FileLog.Error(err.Error())
		return err
	}
	return nil
}
func ParseMailContent(body string) string {

	var content = ""
	// markdown parse
	parsedContent := blackfriday.MarkdownCommon([]byte(body))
	mailTemplateContent, err := ioutil.ReadFile(util.MailTemplatePath + "/" + util.MailTemplateType + ".html")
	if err != nil {
		util.FileLog.Warn(err.Error())
		content = string(parsedContent)
	} else {
		content = strings.Replace(string(mailTemplateContent), "{{MAIL_BODY_"+util.MailBodyKey+"}}", string(parsedContent), -1)
	}

	return content
}

//****Mail END****
