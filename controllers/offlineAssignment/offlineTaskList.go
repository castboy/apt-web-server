/*离线任务列表*/
package offlineAssignment

import (
	"apt-web-server_v2/controllers/public"
	"apt-web-server_v2/models/offlineAssignment"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type TLController struct{}

var TLObj = new(TLController)

func (this *TLController) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//check ssn
	//TODO
	//parse parament
	//var list *models.TblTaskData
	input := OTLGetInput{}

	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		public.Write(w, public.ErrParamentErr, nil)
		return
	}

	/*add all parameter*/
	name := queryForm.Get("name")
	input.Para.Name = name
	creatStart, err := strconv.Atoi(queryForm.Get("start"))
	if err != nil {
		//	start = 0
	}
	input.Para.PField.Start = int64(creatStart)
	end, err := strconv.Atoi(queryForm.Get("end"))
	if err != nil {
		//	end = 0
	}
	input.Para.PField.End = int64(end)
	//	ugctype, type_err := strconv.Atoi(queryForm.Get("type"))
	taskType := queryForm.Get("type")
	input.Para.Type = taskType

	time, err := strconv.Atoi(queryForm.Get("time"))
	if err != nil {
		//page = 1
	}
	input.Para.Time = int64(time)

	status := queryForm.Get("status")
	input.Para.Status = status

	page, err := strconv.Atoi(queryForm.Get("page"))
	if err != nil {
		page = 1
	}
	input.Para.Page = int32(page)
	count, err := strconv.Atoi(queryForm.Get("count"))
	if err != nil {
		//count = 1000
	}
	input.Para.Count = int32(count)

	sort := queryForm.Get("sort")
	if sort == "" {
		sort = "time"
	}
	input.Para.Sort = sort

	order := queryForm.Get("order")
	input.Para.Order = order

	lastCount, err := strconv.Atoi(queryForm.Get("lastCount"))
	if err != nil {
		//		Write(w, ErrParamentErr, nil)
		//		return
	}
	input.Para.LastCount = int32(lastCount)
	//get list
	//	err, list := new(models.TblUrgency).GetLast(input)
	err, list := new(offlineAssignment.TblOLA).GetTaskDetails(&input.Para)
	if err != nil {
		panic(fmt.Sprintf("TblFTD GetLast error:%s", err.Error()))
		return
	}
	//response
	public.Write(w, public.ErrOkErr, list)
}

type OTLGetInput struct {
	Para offlineAssignment.TblTaskSearchPara
}
