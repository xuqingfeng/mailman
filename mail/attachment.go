package mail

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/xuqingfeng/lift"
	"github.com/xuqingfeng/mailman/util"
)

func SaveAttachment(fileContent []byte, token, fileName string) error {

	homeDir, _ := util.GetHomeDir()
	dirPath := filepath.Join(homeDir, util.ConfigPath["tmpPath"], token)
	// ModePerm
	err := lift.CreateDirectory(dirPath, os.ModePerm)
	if err != nil {
		util.FileLog.Error(err.Error())
		return err
	}
	attachmentPath := filepath.Join(homeDir, util.ConfigPath["tmpPath"], token, fileName)
	err = ioutil.WriteFile(attachmentPath, fileContent, os.ModePerm)
	if err != nil {
		util.FileLog.Error(err.Error())
		return err
	}

	return nil
}
