package smtp

import (
	"github.com/xuqingfeng/mailman/util"
	"testing"
)

var (
	testSMTPServer = SMTPServer{
		"@example.com",
		"smtp.example.com",
	}
)

func TestGetSMTPServer(t *testing.T) {

	_, err := GetSMTPServer("test@fakedomain.com")
	t.Logf("err: %v", err)
	if util.SMTPServerNotFoundErr != err {
		t.Error("@example.com SMTP Server exist")
	}
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
