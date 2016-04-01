package main

import (
	"github.com/xuqingfeng/mailman/util"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	util.CreateConfigDir()
	os.Exit(m.Run())
}
