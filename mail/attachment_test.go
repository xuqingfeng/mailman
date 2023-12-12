package mail

import (
	"os"
	"path/filepath"
	"testing"

	"mailman/util"
)

func TestSaveAttachment(t *testing.T) {

	fileContent := []byte("this is content")
	token := "0"
	fileName := "test"

	err := SaveAttachment(fileContent, token, fileName)
	if err != nil {
		t.Errorf("SaveAttachment() fail %v", err)
	}

	homeDir := util.GetHomeDir()
	attachmentPath := filepath.Join(homeDir, util.ConfigPath["tmpPath"], token, fileName)
	if _, err = os.Stat(attachmentPath); err != nil {
		t.Errorf("SaveAttachment() fail %v", err)
	}
}
