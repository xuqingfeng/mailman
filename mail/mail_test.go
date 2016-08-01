package mail

import (
	"os"
	"testing"

	"github.com/xuqingfeng/mailman/account"
	"github.com/xuqingfeng/mailman/smtp"
)

func TestMain(m *testing.M) {

	account4mail := account.Account{
		"from@example.com",
		"password",
	}
	smtp4mail := smtp.SMTPServer{
		"@example.com",
		"smtp.example.com",
		"25",
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
		"1461655086246",
		[]string{},
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
