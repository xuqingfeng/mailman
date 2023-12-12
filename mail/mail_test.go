package mail

import (
	"os"
	"testing"

	"mailman/account"
	"mailman/smtp"
)

func TestMain(m *testing.M) {

	account4mail := account.Account{
		Email:    "from@example.com",
		Password: "password",
	}
	smtp4mail := smtp.SMTPServer{
		Address: "@example.com",
		Server:  "smtp.example.com",
		Port:    "25",
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
	content := ParseMailContent("test")
	if len(content) == 0 {
		t.Error("ParseMailContent() fail")
	}
}
