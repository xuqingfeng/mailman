package mail

import (
	"github.com/xuqingfeng/mailman/util"
	"testing"
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
	if err != util.KeyNotFoundErr {
		t.Error("@example.com key exist ")
	}
}

func TestParseMailContent(t *testing.T) {

    // mail-template not exist
	if "" != ParseMailContent("") {
		t.Error("ParseMailContent() fail")
	}
}
