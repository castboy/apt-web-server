/********获取紧急事件详情********/
package controllers

import (
	"apt-web-server/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func (this *OLPOController) Post(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var list *models.CMDResult
	var err error

	input := Params(r)

	switch input.Para.Cmmand {
	case "creat":
		//		RuleFile(r, "/tmp/rules/", "rule")

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
	var params models.TblOLASearchPara

	bytes := []byte(GetDataString(r))

	json.Unmarshal(bytes, &params)

	input.Para.Name = params.Name

	input.Para.Time = time.Now().Unix()

	input.Para.Type = params.Type

	input.Para.Start = params.Start

	input.Para.End = params.End

	input.Para.Cmmand = params.Cmmand

	input.Para.OfflineTag = "rule"

	input.Para.Weight = params.Weight

	details := params.Details
	if details == "" {
		input.Para.Details = fmt.Sprintf("%s offline dispatch", input.Para.Type)
	} else {
		input.Para.Details = input.Para.Name
	}

	fmt.Println(input)

	return input
}

func GetDataString(req *http.Request) string {
	result, err := ioutil.ReadAll(req.Body)
	if err != nil {
	} else {

	}
	fmt.Println(bytes.NewBuffer(result).String())
	return bytes.NewBuffer(result).String()
}
