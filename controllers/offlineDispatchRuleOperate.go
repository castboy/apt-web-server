/********获取紧急事件详情********/
package controllers

import (
	"apt-web-server/models"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

func (this *OLPOController) Post(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var list *models.CMDResult
	var err error

	input := Params(r)

	switch input.Para.Cmmand {
	case "creat":
		RuleFile(r, "/tmp/rules/", "rule")

		err, list = new(models.TblOLA).CreatAssignment(&input.Para)
	case "delete":
		err, list = new(models.TblOLA).DeleteAssignment(&input.Para)
	case "start":
		err, list = new(models.TblOLA).StartAssignment(&input.Para)

	case "shutdown":
		err, list = new(models.TblOLA).ShutDownAssignment(&input.Para)

	}

	if err != nil {
		panic(fmt.Sprintf("TblUrgency GetLast error:%s", err.Error()))
		return
	}

	Write(w, ErrOkErr, list)
}

func Params(r *http.Request) (input OLPOGetInput) {
	input.Para.Name = r.PostFormValue("name")

	input.Para.Time = time.Now().Unix()

	input.Para.Type = r.PostFormValue("type")

	input.Para.Start = r.PostFormValue("start")

	input.Para.End = r.PostFormValue("end")

	input.Para.Cmmand = r.PostFormValue("cmd")

	input.Para.OfflineTag = "rule"

	weight, err := strconv.Atoi(r.PostFormValue("weight"))
	if nil != err {
		weight = 1
	}
	input.Para.Weight = int(weight)

	details := r.PostFormValue("details")
	if details == "" {
		input.Para.Details = fmt.Sprintf("%s offline dispatch", input.Para.Type)
	} else {
		input.Para.Details = input.Para.Name
	}

	return input
}
func RuleFile(r *http.Request, dstDir, formElement string) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile(formElement)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	f, err := os.OpenFile(dstDir+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
}
