//go:generate statik -src=./ui

package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/abbot/go-http-auth"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rakyll/statik/fs"
	"github.com/urfave/cli"
	"github.com/xuqingfeng/mailman/account"
	"github.com/xuqingfeng/mailman/contacts"
	"github.com/xuqingfeng/mailman/lang"
	"github.com/xuqingfeng/mailman/mail"
	"github.com/xuqingfeng/mailman/smtp"
	_ "github.com/xuqingfeng/mailman/statik"
	"github.com/xuqingfeng/mailman/util"
)

const (
	SPINNER_CHAR_INDEX = 14
	READ_LOG_FILE_GAP  = 5 // second
	MAILMAN_IS_AWESOME = "mailman is awesome !"
	MIN_TCP_PORT       = 0
	MAX_TCP_PORT       = 65535
	//maxReservedTCPPort = 1024
	// 15M
	MAX_MEMORY = 1024 * 1024 * 15

	ASSETS_PREFIX = "ui"
)

var (
	name    = "mailman"
	version = "master"

	msg             util.Msg
	enableBasicAuth = false
	previewContent  = ""
	unauthorized    = "401 Unauthorized"
	loopback        = "127.0.0.1"
	upgrader        = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	ErrDataIsNotJson = errors.New("data is not json format")
)

type Key struct {
	Key string `json:"key"`
}

func main() {

	cyan := color.New(color.FgCyan).SprintFunc()

	colorName := cyan("NAME:")
	colorUsage := cyan("USAGE:")
	colorVersion := cyan("VERSION:")
	colorAuthor := cyan("AUTHOR")
	colorCommands := cyan("COMMANDS")
	colorGlobalOptions := cyan("GLOBAL OPTIONS")

	cli.AppHelpTemplate = colorName + `
    {{.Name}} - {{.Usage}}
` + colorUsage + `
{{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .Flags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}
{{if .Version}}
` + colorVersion + `
{{.Version}}
{{end}}{{if len .Authors}}
` + colorAuthor + `
{{range .Authors}}{{ . }}{{end}}
{{end}}{{if .Commands}}
` + colorCommands + `
{{range .Commands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
{{end}}{{end}}{{if .Flags}}
` + colorGlobalOptions + `
{{range .Flags}}{{.}}
{{end}}{{end}}{{if .Copyright }}
COPYRIGHT:
{{.Copyright}}
{{end}}
`

	app := cli.NewApp()
	app.Name = name
	app.Usage = "Web email client supporting HTML template and SMTP"
	app.Version = version
	app.Author = "xuqingfeng"
	app.Action = func(c *cli.Context) {

		portInUse := -1
		portStart := 8000
		portEnd := 8100
		for portStart <= portEnd {
			if isTCPPortAvailable(portStart) {
				portInUse = portStart
				break
			}
			portStart++
		}
		if -1 == portInUse {
			log.Fatal("can't find available port")
		}

		localIP := getLocalIP()

		if runtime.GOOS == "darwin" {
			_, err := exec.Command("open", "http://"+localIP+":"+strconv.Itoa(portInUse)).Output()
			if err != nil {
				log.Fatalf("darwin open fail: %s", err.Error())
			}
		} else {
			log.Info("Open " + localIP + ":" + strconv.Itoa(portInUse) + " in browser")
		}

		s := spinner.New(spinner.CharSets[SPINNER_CHAR_INDEX], 100*time.Millisecond)
		s.Color("cyan")
		s.Start()

		// util init
		util.CreateConfigDir()

		// router
		router := mux.NewRouter()

		apiSubRouter := router.PathPrefix("/api").Subrouter()
		apiSubRouter.HandleFunc("/ping", pingHandler)
		apiSubRouter.HandleFunc("/lang", langHandler)
		apiSubRouter.HandleFunc("/mail", mailHandler)
		apiSubRouter.HandleFunc("/file", fileHandler)
		apiSubRouter.HandleFunc("/account", accountHandler)
		apiSubRouter.HandleFunc("/contacts", contactsHandler)
		apiSubRouter.HandleFunc("/smtpServer", smtpServerHandler)
		apiSubRouter.HandleFunc("/preview", previewHandler)
		apiSubRouter.HandleFunc("/wslog", wsLogHandler)

		statikFS, err := fs.New()
		if err != nil {
			log.Fatal(err)
		}

		router.PathPrefix("/").Handler(http.FileServer(statikFS))

		http.ListenAndServe(":"+strconv.Itoa(portInUse), router)
	}
	app.Commands = []cli.Command{
		{
			Name:        "clean",
			Usage:       "clean up tmp directory",
			Description: "mailman clean",
			Action: func(c *cli.Context) {
				homeDir := util.GetHomeDir()
				tmpPath := filepath.Join(homeDir, util.ConfigPath["tmpPath"])
				err := os.RemoveAll(tmpPath)
				if err != nil {
					log.Error(err)
				}
				util.CreateConfigDir()
			},
		},
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "basic-auth",
			Usage:       "enable basic auth (~/.mailman/.htpasswd)",
			Destination: &enableBasicAuth,
		},
	}

	app.Run(os.Args)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "pong")
}

func langHandler(w http.ResponseWriter, r *http.Request) {

	if "GET" == r.Method {
		lg, _ := lang.GetLang()

		switch lg {
		case "en", "zh":
			sendSuccess(w, lg, "I! get lang success")
		default:
			sendSuccess(w, "en", "I! get lang success")
		}

	} else if "POST" == r.Method {

		var lg lang.Lang
		err := json.NewDecoder(r.Body).Decode(&lg)
		if err != nil {

			sendError(w, ErrDataIsNotJson.Error())
		} else if err = lang.SaveLang(lg); err != nil {
			sendError(w, "E! save lang fail: "+err.Error())
		} else {
			l, err := lang.GetLang()
			if err != nil {

				sendError(w, "E! get lang fail: "+err.Error())
			} else {

				sendSuccess(w, l, "I! save lang success")
			}
		}
	}
}

func mailHandler(w http.ResponseWriter, r *http.Request) {

	if "GET" == r.Method {

		sendSuccess(w, struct{}{}, MAILMAN_IS_AWESOME)
	} else if "POST" == r.Method {

		var m mail.Mail
		err := json.NewDecoder(r.Body).Decode(&m)
		if err != nil {

			sendError(w, "E! "+ErrDataIsNotJson.Error())
		} else if err = mail.SendMail(m); err != nil {

			sendError(w, "E! send mail fail: "+err.Error())
		} else {
			// empty struct
			sendSuccess(w, struct{}{}, "I! send mail success")
		}
	}
}

func fileHandler(w http.ResponseWriter, r *http.Request) {

	if "GET" == r.Method {

		sendSuccess(w, struct{}{}, MAILMAN_IS_AWESOME)
	} else if "POST" == r.Method {

		if err := r.ParseMultipartForm(MAX_MEMORY); err != nil {
			sendError(w, "E! parse posted file fail: "+err.Error())
		}

		token := ""
		for k, vs := range r.MultipartForm.Value {
			for _, v := range vs {
				if "token" == k {
					token += v
				}
			}
		}

		for _, fileHeaders := range r.MultipartForm.File {
			for _, fileHeader := range fileHeaders {
				f, _ := fileHeader.Open()
				fileContent, _ := ioutil.ReadAll(f)
				err := mail.SaveAttachment(fileContent, token, fileHeader.Filename)
				if err != nil {
					sendError(w, "E! save attachment fail")
					// todo multi
					break
				}
			}
		}
	}
}

func accountHandler(w http.ResponseWriter, r *http.Request) {

	if "GET" == r.Method {

		emails, err := account.GetAccountEmail()
		if err != nil {
			sendError(w, "E! get account email fail: "+err.Error())
		} else {
			// empty []string
			sendSuccess(w, emails, "I! get account email success")
		}
	} else if "POST" == r.Method {

		var at account.Account
		err := json.NewDecoder(r.Body).Decode(&at)
		if err != nil {

			sendError(w, "E! "+ErrDataIsNotJson.Error())
		} else if err = account.SaveAccount(at); err != nil {

			sendError(w, "E! save account fail: "+err.Error())
		} else {
			emails, err := account.GetAccountEmail()
			if err != nil {

				sendError(w, "E! get account email fail: "+err.Error())
			} else {

				sendSuccess(w, emails, "I! save account success")
			}
		}
	} else if "DELETE" == r.Method {

		var k Key
		err := json.NewDecoder(r.Body).Decode(&k)
		if err != nil {

			sendError(w, "E! "+ErrDataIsNotJson.Error()+" "+err.Error())
		} else if err = account.DeleteAccount(k.Key); err != nil {

			sendError(w, "E! delete account fail: "+err.Error())
		} else {
			emails, err := account.GetAccountEmail()
			if err != nil {

				sendError(w, "E! get account email fail: "+err.Error())
			} else {

				sendSuccess(w, emails, "I! delete account success")
			}
		}
	}
}

func contactsHandler(w http.ResponseWriter, r *http.Request) {

	if "GET" == r.Method {

		c, err := contacts.GetContacts()
		if err != nil {

			sendError(w, "E! get contacts fail: "+err.Error())
		} else {

			sendSuccess(w, c, "I! get contacts success")
		}
	} else if "POST" == r.Method {

		var ct contacts.Contacts
		err := json.NewDecoder(r.Body).Decode(&ct)
		if err != nil {

			sendError(w, ErrDataIsNotJson.Error())
		} else if err = contacts.SaveContacts(ct); err != nil {

			sendError(w, "E! save contacts fail: "+err.Error())
		} else {
			c, err := contacts.GetContacts()
			if err != nil {

				sendError(w, "E! get contacts fail: "+err.Error())
			} else {

				sendSuccess(w, c, "I! save contacts success")
			}
		}
	} else if "DELETE" == r.Method {
		var k Key
		err := json.NewDecoder(r.Body).Decode(&k)
		if err != nil {

			sendError(w, ErrDataIsNotJson.Error()+" "+err.Error())
		} else if err = contacts.DeleteContacts(k.Key); err != nil {

			sendError(w, "E! delete contacts fail: "+err.Error())
		} else {
			c, err := contacts.GetContacts()
			if err != nil {

				sendError(w, "E! get contacts fail: "+err.Error())
			} else {

				sendSuccess(w, c, "I! delete contacts success")
			}
		}
	}
}

func smtpServerHandler(w http.ResponseWriter, r *http.Request) {

	if "GET" == r.Method {

		customSMTPServer, err := smtp.GetCustomSMTPServer()
		if err != nil {

			sendError(w, "E! get custom SMTP server fail: "+err.Error())
		} else {

			sendSuccess(w, customSMTPServer, "I! get custom SMTP Server success")
		}
	} else if "POST" == r.Method {

		var smtpServer smtp.SMTPServer
		err := json.NewDecoder(r.Body).Decode(&smtpServer)
		if err != nil {

			sendError(w, "E! "+ErrDataIsNotJson.Error())
		} else if err = smtp.SaveSMTPServer(smtpServer); err != nil {

			sendError(w, "E! "+err.Error())
		} else {

			customSMTPServer, err := smtp.GetCustomSMTPServer()
			if err != nil {

				sendError(w, "E! "+err.Error())
			} else {

				sendSuccess(w, customSMTPServer, "I! save SMTP Server success")
			}
		}
	} else if "DELETE" == r.Method {
		var k Key
		err := json.NewDecoder(r.Body).Decode(&k)
		if err != nil {

			sendError(w, "E! "+ErrDataIsNotJson.Error()+" "+err.Error())
		} else if err = smtp.DeleteSMTPServer(k.Key); err != nil {

			sendError(w, "E! delete SMTPServer fail: "+err.Error())
		} else {
			server, err := smtp.GetCustomSMTPServer()
			if err != nil {

				sendError(w, "E! get custom SMTP server fail: "+err.Error())
			} else {

				sendSuccess(w, server, "I! delete SMTP server success")
			}
		}
	}
}

func previewHandler(w http.ResponseWriter, r *http.Request) {

	if "GET" == r.Method {

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(previewContent))
	} else if "POST" == r.Method {

		type Body struct {
			Body string `json:"body"`
		}
		var body Body
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {

			sendError(w, "E! "+ErrDataIsNotJson.Error())
		} else {

			previewContent = mail.ParseMailContent(body.Body)
			sendSuccess(w, struct{}{}, previewContent)
		}
	}
}

func wsLogHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err.Error())
	}

	homeDir := util.GetHomeDir()
	logFilePath := filepath.Join(homeDir, util.ConfigPath["logPath"], util.LogName)
	logFile, err := os.Open(logFilePath)
	if err != nil {
		log.Error(err.Error())
	}
	reader := bufio.NewReader(logFile)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			log.Fatal(err.Error())
		} else if err == io.EOF {
			// wait
			time.Sleep(READ_LOG_FILE_GAP * time.Second)
		} else {
			if err = conn.WriteMessage(1, []byte(line)); err != nil {
				log.Error(err.Error())
			}
		}
	}
}

func basicAuth(w http.ResponseWriter, r *http.Request) bool {

	if enableBasicAuth {
		ba := auth.NewBasicAuthenticator(fmt.Sprintf("%s / %s", name, version), auth.HtpasswdFileProvider(filepath.Join(util.GetHomeDir(), util.ConfigPath["htpasswdPath"])))
		if ba.CheckAuth(r) == "" {
			w.Header().Set("WWW-Authenticate", `Basic realm="realm"`)
			return false
		}
	}

	return true
}

func sendSuccess(w http.ResponseWriter, data interface{}, message string) {

	msg = util.Msg{
		Success: true,
		Data:    data,
		Message: message,
	}

	msgInByteSlice, _ := json.Marshal(msg)
	w.Header().Set("Content-Type", "application/json")
	w.Write(msgInByteSlice)
}
func sendError(w http.ResponseWriter, message string) {

	msg = util.Msg{
		Success: false,
		Message: message,
	}

	msgInByteSlice, _ := json.Marshal(msg)
	w.Header().Set("Content-Type", "application/json")
	w.Write(msgInByteSlice)
}

func isTCPPortAvailable(port int) bool {

	if port < MIN_TCP_PORT || port > MAX_TCP_PORT {
		return false
	}
	conn, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func getLocalIP() (ip string) {

	ip = loopback
	ifaces, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String()
				}
			}
		}
	}

	return
}
