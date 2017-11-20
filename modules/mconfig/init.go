package mconfig

import (
	"apt-web-server_v2/modules/mlog"
	"os"
	"strconv"

	"github.com/larspensjo/config"
)

var Conf *config.Config

func init() {
	conf, err := config.ReadDefault("conf/app.conf")
	if err != nil {
		os.Exit(1)
	}
	Conf = conf
	return
}

func Writepid() {
	f, err := os.Create("/tmp/apt-web-server.pid")
	defer f.Close()
	if err != nil {
		mlog.Debug("create PIDfile error:", err)
	}
	_, err = f.WriteString(strconv.Itoa(os.Getpid()))
	if err != nil {
		mlog.Debug("write PID error:", err)
	}
}
