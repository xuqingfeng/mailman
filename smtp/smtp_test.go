package smtp

import "testing"

var (
	testSMTPServer = SMTPServer{
		"@example.com",
		"smtp.example.com",
	}
)

func TestSaveSMTPServer(t *testing.T) {

	err := SaveSMTPServer(testSMTPServer)
	if err != nil {
		t.Errorf("SaveSMTPServer() fail %v", err)
	}
}

func TestDeleteSMTPServer(t *testing.T) {

	err := DeleteSMTPServer(testSMTPServer.Address)
	if err != nil {
		t.Errorf("DeleteSMTPServer() fail %v", err)
	}
}
