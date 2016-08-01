package util

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"

	"github.com/Sirupsen/logrus"
)

const (
	configMode       = 0755
	dbName           = "mailman.db"
	LogName          = "mailman.log"
	MailTemplatePath = "./ui/mail-template"
	MailTemplateType = "responsive"
	//MailTemplateType = "stationery"
)

var (
	FileLog    = logrus.New()
	ConfigPath = map[string]string{
		"dbPath":  "/.mailman/db",
		"logPath": "/.mailman/log",
		"tmpPath": "/.mailman/tmp",
	}
	DBPath                string
	DefaultLang           = "en"
	SMTPServerNotFoundErr = errors.New("SMTP Server Not Found")
)

func init() {

	homeDir := GetHomeDir()
	CreateConfigDir()
	logFilePath := filepath.Join(homeDir, ConfigPath["logPath"], LogName)
	logFile, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE, configMode)
	if err != nil {
		// mailman.log not exist
		FileLog.Fatal(err.Error())
		//log.Fatal(err)
	}
	FileLog.Out = logFile
	FileLog.Formatter = &logrus.TextFormatter{DisableColors: true}
	DBPath = filepath.Join(homeDir, ConfigPath["dbPath"], dbName)
}

type Msg struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,emitempty"`
}

func GetHomeDir() string {

	return os.Getenv("HOME")
}

func CreateConfigDir() error {

	homeDir := GetHomeDir()
	for _, path := range ConfigPath {
		var p = homeDir + path
		if _, err := os.Stat(p); err != nil {
			if os.IsNotExist(err) {
				if err := os.MkdirAll(p, configMode); err != nil {
					FileLog.Fatal(err.Error())
					return err
				}
			} else {
				FileLog.Fatal(err.Error())
				return err
			}
		}
	}
	return nil
}

func GetTmpDir() string {

	homeDir := GetHomeDir()
	return filepath.Join(homeDir, ConfigPath["tmpPath"])
}

func GetUserName() string {

	u, _ := user.Current()
	return u.Username
}
