package smtp

import (
	"testing"

	"mailman/util"
)

var (
	testSMTPServer = SMTPServer{
		"@example.com",
		"smtp.example.com",
		"25",
	}
)

func TestGetSMTPServer(t *testing.T) {

	fakeEmailAddress0 := "test@example.net"
	_, err := GetSMTPServer(fakeEmailAddress0)
	if util.SMTPServerNotFoundErr != err {
		t.Error("@example.net SMTP Server exists")
	}

	fakeEmailAddress1 := "test@example.com"
	SaveSMTPServer(testSMTPServer)
	_, err = GetSMTPServer(fakeEmailAddress1)
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
