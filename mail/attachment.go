package mail

import (
	"github.com/xuqingfeng/lift"
	"github.com/xuqingfeng/mailman/util"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
)

func SaveAttachment(f *multipart.FileHeader, token string) error {

	homeDir, _ := util.GetHomeDir()
	file, _ := f.Open()
	dirPath := filepath.Join(homeDir, util.ConfigPath["tmpPath"], token)
	// ModePerm
	err := lift.CreateDirectory(dirPath, os.ModePerm)
	if err != nil {
		util.FileLog.Error(err.Error())
		return err
	}
	path := filepath.Join(homeDir, util.ConfigPath["tmpPath"], token, f.Filename)
	// 中文乱码
	//path := filepath.Join(homeDir, util.ConfigPath["tmpPath"], token, "中文测试")
	buf, _ := ioutil.ReadAll(file)
	err = ioutil.WriteFile(path, buf, os.ModePerm)
	if err != nil {
		util.FileLog.Error(err.Error())
		return err
	}

	return nil
}
