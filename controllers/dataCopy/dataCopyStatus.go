/*数据迁移任务状态*/
package dataCopy

import (
	"apt-web-server_v2/controllers/public"
	"apt-web-server_v2/models/dataCopy"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type DCTSController struct{}

var DCTSObj = new(DCTSController)

func (this *DCTSController) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//check ssn
	//TODO
	//parse parament
	input := DCTSGetInput{}

	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		public.Write(w, public.ErrParamentErr, nil)
		return
	}
	location := queryForm.Get("location")
	input.Para.Location = location

	date, err := strconv.Atoi(queryForm.Get("date"))
	if err != nil {
		//time = 0
	}
	input.Para.Date = int64(date)

	rate, err := strconv.Atoi(queryForm.Get("rate"))
	if err != nil {
		//rate = 0
	}
	input.Para.Rate = rate

	status := queryForm.Get("status")
	input.Para.Status = status

	details := queryForm.Get("details")
	input.Para.Details = details

	command := queryForm.Get("cmd")
	var list *dataCopy.DCTSList
	var lastTime *dataCopy.DCTSLastTime
	switch command {
	case "get":
		err, list = new(dataCopy.TblDCT).GetTaskStatus(&input.Para)
	case "update":
		err, list = new(dataCopy.TblDCT).UpdateTaskStatus(&input.Para)
	case "lasttime":
		err, lastTime = new(dataCopy.TblDCT).GetLastTime(&input.Para)
	default:
		public.Write(w, public.ErrOkErr, "command error")
	}

	if err != nil {
		panic(fmt.Sprintf("TblDCT %s error:%s", command, err.Error()))
		return
	}
	//response
	switch command {
	case "get":
		public.Write(w, public.ErrOkErr, list)
	case "update":
		public.Write(w, public.ErrOkErr, list)
	case "lasttime":
		public.Write(w, public.ErrOkErr, lastTime)
	}
	//Write(w, ErrOkErr, list)
}

type DCTSGetInput struct {
	Para dataCopy.DCTSPara
}
