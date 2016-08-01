package mail

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/russross/blackfriday"
	"github.com/xuqingfeng/mailman/account"
	"github.com/xuqingfeng/mailman/smtp"
	"github.com/xuqingfeng/mailman/util"
	"gopkg.in/gomail.v2"
	"strconv"
)

type Mail struct {
	Subject             string   `json:"subject"`
	To                  []string `json:"to"`
	Cc                  []string `json:"cc"`
	From                string   `json:"from"`
	Priority            bool     `json:"priority"`
	Body                string   `json:"body"`
	Token               string   `json:"token"`
	AttachmentFileNames []string `json:"attachmentFileNames"`
}

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

	// Cc is empty
	if len(mail.Cc) > 0 {
		// multiple cc
		m.SetHeader("Cc", mail.Cc...)
	}
	// todo with name - SetAddressHeader
	//m.SetAddressHeader("Cc", mail.Cc[0].Email, mail.Cc[0].Name)
	m.SetHeader("From", account.Email)

	// priority
	if mail.Priority {
		m.SetHeader("X-Priority", "1")
	}

	// attachment
	tmpDir := util.GetTmpDir()
	attachmentDir := filepath.Join(tmpDir, mail.Token)
	if len(mail.AttachmentFileNames) > 0 {
		for _, v := range mail.AttachmentFileNames {
			attachmentPath := filepath.Join(attachmentDir, v)
			m.Attach(attachmentPath)
		}
	}

	content := ParseMailContent(mail.Body)

	m.SetBody("text/html", content)

	port, err := strconv.Atoi(SMTPServer.Port)
	if err != nil {
		return err
	}
	d := gomail.NewDialer(SMTPServer.Server, port, account.Email, account.Password)
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
