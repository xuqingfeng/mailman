package mail

import (
	"github.com/xuqingfeng/lift"
	"github.com/xuqingfeng/mailman/util"
	"io/ioutil"
	"os"
	"path/filepath"
)

// SaveAttachment(f *multipart.FileHeader, token string) 不知道如何测试,更改方法
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
	// 中文乱码
	//attachmentPath := filepath.Join(homeDir, util.ConfigPath["tmpPath"], token, "中文测试")
	err = ioutil.WriteFile(attachmentPath, fileContent, os.ModePerm)
	if err != nil {
		util.FileLog.Error(err.Error())
		return err
	}

	return nil
}
