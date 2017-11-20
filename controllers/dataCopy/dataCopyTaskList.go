/*数据迁移任务列表*/
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

type DCTLController struct{}

var DCTLObj = new(DCTLController)

func (this *DCTLController) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//check ssn
	//TODO
	//parse parament
	//var list *models.TblTaskData
	input := DCTLGetInput{}

	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		public.Write(w, public.ErrParamentErr, nil)
		return
	}

	/*add all parameter*/
	location := queryForm.Get("location")
	input.Para.Location = location
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
	if sort == "" || sort == "date" {
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
	err, list := new(dataCopy.TblDCT).GetDCTL(&input.Para)
	if err != nil {
		panic(fmt.Sprintf("TblDCTL GetLast error:%s", err.Error()))
		return
	}
	//response
	public.Write(w, public.ErrOkErr, list)
}

type DCTLGetInput struct {
	Para dataCopy.TblDCTLSearchPara
}
