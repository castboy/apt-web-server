/********获取紧急事件详情********/
package controllers

import (
	"apt-web-server/models"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type OLPOController struct{}

var OLPOObj = new(OLPOController)

func (this *OLPOController) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//check ssn
	//TODO
	//parse parament
	var list *models.CMDResult
	input := OLPOGetInput{}

	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		Write(w, ErrParamentErr, nil)
		return
	}
	name := queryForm.Get("name")
	input.Para.Name = name

	time, err := strconv.Atoi(queryForm.Get("time"))
	if err != nil {
		//time = 0
	}
	input.Para.Time = int64(time)
	cmdtype := queryForm.Get("type")
	input.Para.Type = cmdtype

	offlinetag := queryForm.Get("offlinetag")
	if offlinetag == "" {
		input.Para.OfflineTag = "offline"
	} else {
		input.Para.OfflineTag = offlinetag
	}

	start := queryForm.Get("start")
	if err != nil {
		//start = 0
	}
	input.Para.Start = start
	end := queryForm.Get("end")
	if err != nil {
		//end = 0
	}
	input.Para.End = end
	weight, err := strconv.Atoi(queryForm.Get("weight"))
	if err != nil {
		//weight = 0
	}
	input.Para.Weight = int(weight)
	details := queryForm.Get("details")
	if details == "" {
		input.Para.Details = fmt.Sprintf("%s offline dispatch", input.Para.Type)
	} else {
		input.Para.Details = name
	}

	command := queryForm.Get("cmd")
	fmt.Println(input)
	switch command {
	case "creat":
		err, list = new(models.TblOLA).CreatAssignment(&input.Para)
	case "delete":
		err, list = new(models.TblOLA).DeleteAssignment(&input.Para)
	case "start":
		err, list = new(models.TblOLA).StartAssignment(&input.Para)
		/*	case "stop":
			err, list = new(models.TblOLA).StopAssignment(&input.Para)
		*/
	case "shutdown":
		err, list = new(models.TblOLA).ShutDownAssignment(&input.Para)

	}

	if err != nil {
		panic(fmt.Sprintf("TblUrgency GetLast error:%s", err.Error()))
		return
	}
	//response
	//Write(w, ErrOkErr, queryForm)
	Write(w, ErrOkErr, list)
}

type OLPOGetInput struct {
	Para models.TblOLASearchPara
}
