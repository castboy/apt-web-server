/********获取紧急事件详情********/
package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

const (
	RULESDIR = "/tmp/rules"
)

type RULEController struct{}

var RULEObj = new(RULEController)

func (this *RULEController) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()

	var topic string

	if val, ok := r.Form["topic"]; ok {
		topic = val[0]
	}

	cont := RdFile(RULESDIR + "/" + topic)
	res := struct {
		Rule string
		Cont string
	}{topic, cont}

	byte, _ := json.Marshal(res)

	io.WriteString(w, string(byte))
}

func RdFile(file string) string {
	fHdl, err := os.Open(file)
	if nil != err {
		if pe, ok := err.(*os.PathError); ok {
			fmt.Printf("Path Error: %s (op=%s, path=%s)\n", pe.Err, pe.Op, pe.Path)
		} else {
			fmt.Printf("Unknow Error: %s\n", err)
		}
	}

	defer fHdl.Close()

	fd, _ := ioutil.ReadAll(fHdl)
	return string(fd)
}
