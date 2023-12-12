package smtp

import (
	"regexp"

	"encoding/json"
	"mailman/util"
)

type SMTPServer struct {
	Address string `json:"address"`
	Server  string `json:"server"`
	Port    string `json:"port"`
}

var (
	defaultSMTPServer = []SMTPServer{
		{
			"@qq.com", "smtp.qq.com", "25",
		},
		{
			"@hotmail.com", "smtp.live.com", "25",
		},
		{
			"@outlook.com", "smtp.live.com", "25",
		},
		{
			"@gmail.com", "smtp.gmail.com", "25",
		},
	}
)

func GetSMTPServer(email string) (SMTPServer, error) {

	customSMTPServer, _ := GetCustomSMTPServer()
	for _, v := range customSMTPServer {
		ok, err := regexp.MatchString(".*"+v.Address, email)
		if err != nil {
			util.FileLog.Error(err.Error())
		}
		if ok {
			return v, nil
		}
	}
	for _, v := range defaultSMTPServer {
		ok, err := regexp.MatchString(".*"+v.Address, email)
		if err != nil {
			util.FileLog.Error(err.Error())
		}
		if ok {
			return v, nil
		}
	}

	return SMTPServer{}, util.SMTPServerNotFoundErr
}
func GetCustomSMTPServer() ([]SMTPServer, error) {

	boltStore, err := util.NewBoltStore(util.DBPath)
	if err != nil {
		return nil, err
	}
	defer boltStore.Close()

	customSMTPServer, order, err := boltStore.GetRange(util.SmtpBucketName)
	if err != nil {
		util.FileLog.Error(err.Error())
		return nil, err
	}
	var customSMTPServerList []SMTPServer
	for _, v := range order {
		smtpServer := customSMTPServer[v]
		var server SMTPServer
		err = json.Unmarshal([]byte(smtpServer), &server)
		if err != nil {
			return customSMTPServerList, err
		}
		customSMTPServerList = append(customSMTPServerList, server)
	}

	return customSMTPServerList, nil
}

// user setting first
func SaveSMTPServer(smtpServer SMTPServer) error {

	boltStore, err := util.NewBoltStore(util.DBPath)
	if err != nil {
		return err
	}
	defer boltStore.Close()

	encoded, err := json.Marshal(smtpServer)
	if err != nil {
		return err
	}
	return boltStore.Set([]byte(smtpServer.Address), encoded, util.SmtpBucketName)
}
func DeleteSMTPServer(address string) error {

	boltStore, err := util.NewBoltStore(util.DBPath)
	if err != nil {
		return err
	}
	defer boltStore.Close()

	return boltStore.Delete([]byte(address), util.SmtpBucketName)
}
