package mail

import (
	"github.com/xuqingfeng/mailman/util"
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAttachment(t *testing.T) {

	fileContent := []byte("this is content")
	token := "0"
	fileName := "test"

	err := SaveAttachment(fileContent, token, fileName)
	if err != nil {
		t.Errorf("SaveAttachment() fail %v", err)
	}

	homeDir, _ := util.GetHomeDir()
	attachmentPath := filepath.Join(homeDir, util.ConfigPath["tmpPath"], token, fileName)
	if _, err = os.Stat(attachmentPath); err != nil {
		t.Errorf("SaveAttachment() fail %v", err)
	}
}
