package util

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/rakyll/statik/fs"
	"github.com/sirupsen/logrus"
)

const (
	configMode       = 0755
	dbName           = "mailman.db"
	LogName          = "mailman.log"
	MailTemplateType = "responsive"
)

var (
	FileLog    = logrus.New()
	ConfigPath = map[string]string{
		"dbPath":       "/.mailman/db",
		"logPath":      "/.mailman/log",
		"tmpPath":      "/.mailman/tmp",
		"htpasswdPath": "/.mailman/.htpasswd",
	}
	DBPath                string
	DefaultLang           = "en"
	ErrSMTPServerNotFound = errors.New("SMTP server not found")
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
	Data    interface{} `json:"data,omitempty"`
}

func GetHomeDir() string {

	return os.Getenv("HOME")
}

func CreateConfigDir() error {

	homeDir := GetHomeDir()
	for _, path := range ConfigPath {
		// dirty fix
		if path == "/.mailman/.htpasswd" {
			continue
		}
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

func GetUser() string {

	return os.Getenv("USER")
}

func CreateDirectory(path string, defaultMode os.FileMode) error {

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(path, defaultMode); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func GetContentFromStatik(name string) (string, error) {

	staticFS, err := fs.New()
	if err != nil {
		return "", err
	}

	hf, err := staticFS.Open(name)
	if err != nil {
		return "", err
	}
	defer hf.Close()

	content, err := io.ReadAll(hf)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
