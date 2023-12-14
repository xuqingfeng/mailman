package mail

import (
	"os"
	"path/filepath"

	"mailman/util"
)

func SaveAttachment(fileContent []byte, token, fileName string) error {

	homeDir := util.GetHomeDir()
	dirPath := filepath.Join(homeDir, util.ConfigPath["tmpPath"], token)
	// ModePerm
	err := util.CreateDirectory(dirPath, os.ModePerm)
	if err != nil {
		util.FileLog.Error(err.Error())
		return err
	}
	attachmentPath := filepath.Join(homeDir, util.ConfigPath["tmpPath"], token, fileName)
	err = os.WriteFile(attachmentPath, fileContent, os.ModePerm)
	if err != nil {
		util.FileLog.Error(err.Error())
		return err
	}

	return nil
}
