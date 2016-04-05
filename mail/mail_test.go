package mail

import (
    "testing"
    "github.com/xuqingfeng/mailman/util"
)

func TestSendMail(t *testing.T) {

    exampleMail := Mail{
        "example",
        []string{"to@example.com"},
        []string{},
        "from@example.com",
        "example body",
    }
    err := SendMail(exampleMail)
    if err != util.SMTPServerNotFoundErr {
        t.Error("@example.com SMTP Server exist ")
    }
}

func TestParseMailContent(t *testing.T) {

    if "" == ParseMailContent("") {
        t.Errorf("")
    }
}