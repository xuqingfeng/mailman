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
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/abbot/go-http-auth"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/urfave/cli"
	"github.com/xuqingfeng/mailman/account"
	"github.com/xuqingfeng/mailman/contacts"
	"github.com/xuqingfeng/mailman/lang"
	"github.com/xuqingfeng/mailman/mail"
	"github.com/xuqingfeng/mailman/smtp"
	"github.com/xuqingfeng/mailman/util"
)

const (
	name    = "mailman"
	version = "0.4.1"

	SPINNER_CHAR_INDEX = 14
	READ_LOG_FILE_GAP  = 5 // second
	HI_THERE           = "HI THERE !"
	MIN_TCP_PORT       = 0
	MEX_TCP_PORT       = 65535
	//maxReservedTCPPort = 1024
	// 15M
	MAX_MEMORY = 1024 * 1024 * 15

	ASSETS_PREFIX = "ui"
)

var (
	msg             util.Msg
	enableBasicAuth = false
	previewContent  = ""
	upgrader        = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	ErrDataIsNotJson = errors.New("data is not json formated")
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

		// FIXME: 16/8/31
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

		if runtime.GOOS == "darwin" {
			_, err := exec.Command("open", "http://127.0.0.1:"+strconv.Itoa(portInUse)).Output()
			if err != nil {
				log.Fatalf("darwin open fail: %s", err.Error())
			}
		} else {
			log.Info("Open 127.0.0.1:" + strconv.Itoa(portInUse) + " in browser")
		}

		s := spinner.New(spinner.CharSets[SPINNER_CHAR_INDEX], 100*time.Millisecond)
		s.Color("cyan")
		s.Start()

		// util init
		util.CreateConfigDir()

		// router
		router := mux.NewRouter()

		apiSubRouter := router.PathPrefix("/api").Subrouter()
		apiSubRouter.HandleFunc("/ping", PingHandler)
		apiSubRouter.HandleFunc("/lang", LangHandler)
		apiSubRouter.HandleFunc("/mail", MailHandler)
		apiSubRouter.HandleFunc("/file", FileHandler)
		apiSubRouter.HandleFunc("/account", AccountHandler)
		apiSubRouter.HandleFunc("/contacts", ContactsHandler)
		apiSubRouter.HandleFunc("/smtpServer", SMTPServerHandler)
		apiSubRouter.HandleFunc("/preview", PreviewHandler)
		apiSubRouter.HandleFunc("/wslog", WSLogHandler)

		// / /index /setting /log
		rootSubRouter := router.PathPrefix("/").Subrouter()
		rootSubRouter.HandleFunc("/", IndexHandler)
		rootSubRouter.HandleFunc("/index", IndexHandler)
		rootSubRouter.HandleFunc("/setting", SettingHandler)
		rootSubRouter.HandleFunc("/log", LogHandler)
		rootSubRouter.HandleFunc("/robots.txt", RobotsHandler)

		// /assets
		router.HandleFunc("/assets/"+`{path:\S+}`, AssetsHandler)

		http.ListenAndServe(":"+strconv.Itoa(portInUse), router)
	}
	app.Commands = []cli.Command{
		{
			Name:        "clean",
			Usage:       "clean up tmp dir",
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
			Usage:       "enable basic auth for mailman",
			Destination: &enableBasicAuth,
		},
	}

	app.Run(os.Args)
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "pong")
}

func LangHandler(w http.ResponseWriter, r *http.Request) {

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

func MailHandler(w http.ResponseWriter, r *http.Request) {

	if "GET" == r.Method {

		sendSuccess(w, struct{}{}, HI_THERE)
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

func FileHandler(w http.ResponseWriter, r *http.Request) {

	if "GET" == r.Method {

		sendSuccess(w, struct{}{}, HI_THERE)
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

func AccountHandler(w http.ResponseWriter, r *http.Request) {

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

func ContactsHandler(w http.ResponseWriter, r *http.Request) {

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

func SMTPServerHandler(w http.ResponseWriter, r *http.Request) {

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

func PreviewHandler(w http.ResponseWriter, r *http.Request) {

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

func WSLogHandler(w http.ResponseWriter, r *http.Request) {

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

func IndexHandler(w http.ResponseWriter, r *http.Request) {

	if !basicAuth(w, r) {
		http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
		return
	}

	asset, err := Asset(ASSETS_PREFIX + "/index.html")
	if err != nil {
		fmt.Fprint(w, http.StatusNotFound)
	}
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, string(asset))
}

func SettingHandler(w http.ResponseWriter, r *http.Request) {

	if !basicAuth(w, r) {
		http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
		return
	}

	asset, err := Asset(ASSETS_PREFIX + "/setting.html")
	if err != nil {
		fmt.Fprint(w, http.StatusNotFound)
	}
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, string(asset))
}

func LogHandler(w http.ResponseWriter, r *http.Request) {

	if !basicAuth(w, r) {
		http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
		return
	}

	asset, err := Asset(ASSETS_PREFIX + "/log.html")
	if err != nil {
		fmt.Fprint(w, http.StatusNotFound)
	}
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, string(asset))
}

func RobotsHandler(w http.ResponseWriter, r *http.Request) {

	asset, err := Asset(ASSETS_PREFIX + "/robots.txt")
	if err != nil {
		fmt.Fprint(w, http.StatusNotFound)
	}
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, string(asset))
}

func AssetsHandler(w http.ResponseWriter, r *http.Request) {

	if !basicAuth(w, r) {
		http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	asset, err := Asset(ASSETS_PREFIX + "/assets/" + vars["path"])
	if err != nil {
		fmt.Fprint(w, err.Error())
	}
	if strings.HasSuffix(vars["path"], ".css") {
		// fixed (Stylesheet)
		w.Header().Set("Content-Type", "text/css")
	} else if strings.HasSuffix(vars["path"], ".js") {
		// fixed
		w.Header().Set("Content-Type", "text/javascript")
	} else if strings.HasSuffix(vars["path"], "png") {
		w.Header().Set("Content-Type", "image/png")
	} else if strings.HasSuffix(vars["path"], "ico") {
		w.Header().Set("Content-Type", "image/x-icon")
	} else if strings.HasSuffix(vars["path"], "xml") || strings.HasSuffix(vars["path"], "svg") {
		w.Header().Set("Content-Type", "text/xml")
	} else if strings.HasSuffix(vars["path"], "json") {
		w.Header().Set("Content-Type", "application/json")
	}
	fmt.Fprint(w, string(asset))
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

	if port < MIN_TCP_PORT || port > MEX_TCP_PORT {
		return false
	}
	conn, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(port))
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
