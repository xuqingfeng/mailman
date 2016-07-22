package smtp

import (
	"regexp"

	"github.com/xuqingfeng/mailman/util"
)

type SMTPServer struct {
	Address string `json:"address"`
	Server  string `json:"server"`
}

func GetSMTPServer(email string) (string, error) {

	customSMTPServer, _ := GetCustomSMTPServer()
	for _, v := range customSMTPServer {
		ok, err := regexp.MatchString(".*"+v.Address, email)
		if err != nil {
			util.FileLog.Error(err.Error())
		}
		if ok {
			return v.Server, nil
		}
	}
	for k, v := range util.DefaultSMTPServer {
		ok, err := regexp.MatchString(".*"+k, email)
		if err != nil {
			util.FileLog.Error(err.Error())
		}
		if ok {
			return v, nil
		}
	}

	return "", util.SMTPServerNotFoundErr
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
		customSMTPServerList = append(customSMTPServerList, SMTPServer{v, customSMTPServer[v]})
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

	return boltStore.Set([]byte(smtpServer.Address), []byte(smtpServer.Server), util.SmtpBucketName)
}
func DeleteSMTPServer(address string) error {

	boltStore, err := util.NewBoltStore(util.DBPath)
	if err != nil {
		return err
	}
	defer boltStore.Close()

	return boltStore.Delete([]byte(address), util.SmtpBucketName)
}
