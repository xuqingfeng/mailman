package mail

import (
	"github.com/xuqingfeng/mailman/account"
	"github.com/xuqingfeng/mailman/smtp"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	account4mail := account.Account{
		"from@example.com",
		"password",
	}
	smtp4mail := smtp.SMTPServer{
		"@example.com",
		"smtp.example.com",
	}
	account.SaveAccount(account4mail)
	smtp.SaveSMTPServer(smtp4mail)
	m.Run()
	account.DeleteAccount(account4mail.Email)
	smtp.DeleteSMTPServer(smtp4mail.Address)
	os.Exit(0)
}

func TestSendMail(t *testing.T) {

	exampleMail := Mail{
		"example",
		[]string{"to@example.com"},
		[]string{},
		"from@example.com",
		false,
		"example body",
	}
	err := SendMail(exampleMail)
	if err == nil {
		t.Error("SendMail() should fail")
	}
}

func TestParseMailContent(t *testing.T) {

	// mail-template not exist
	if "" != ParseMailContent("") {
		t.Error("ParseMailContent() fail")
	}
}
