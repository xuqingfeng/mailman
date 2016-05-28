package smtp

import (
	"testing"

	"github.com/xuqingfeng/mailman/util"
)

var (
	testSMTPServer = SMTPServer{
		"@example.com",
		"smtp.example.com",
	}
)

func TestGetSMTPServer(t *testing.T) {

	fakeEmailAddress := "test@example.com"
	_, err := GetSMTPServer(fakeEmailAddress)
	t.Logf("err: %v", err)
	if util.SMTPServerNotFoundErr != err {
		t.Error("@example.com SMTP Server exist")
	}
	SaveSMTPServer(testSMTPServer)
	_, err = GetSMTPServer(fakeEmailAddress)
	if err != nil {
		t.Errorf("GetSMTPServer() fail %v", err)
	}
	DeleteSMTPServer(testSMTPServer.Address)
}

func TestSaveSMTPServer(t *testing.T) {

	err := SaveSMTPServer(testSMTPServer)
	if err != nil {
		t.Errorf("SaveSMTPServer() fail %v", err)
	}
}

func TestGetCustomSMTPServer(t *testing.T) {

	customSMTPServerList, err := GetCustomSMTPServer()
	if err != nil {
		t.Errorf("GetCustomSMTPServer() fail %v", err)
	}
	if len(customSMTPServerList) < 1 {
		t.Error("GetCustomSMTPServer() fail")
	}
}

func TestDeleteSMTPServer(t *testing.T) {

	err := DeleteSMTPServer(testSMTPServer.Address)
	if err != nil {
		t.Errorf("DeleteSMTPServer() fail %v", err)
	}
}
