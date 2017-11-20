/********获取紧急事件详情********/
package urgency

import (
	"apt-web-server_v2/controllers/public"
	"apt-web-server_v2/models/urgency"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type UrgencyController struct{}

var UrgencyObj = new(UrgencyController)

func (this *UrgencyController) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//check ssn
	//TODO
	//parse parament
	//input := models.TblUrgencySearchPara{}
	input := UrgencyGetInput{}

	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		public.Write(w, public.ErrParamentErr, nil)
		return
	}

	if err != nil {
		//		Write(w, ErrParamentErr, nil)
		//		return
	}
	/*add all parameter*/
	start, err := strconv.Atoi(queryForm.Get("start"))
	if err != nil {
		//	start = 0
	}
	input.Para.PField.Start = int64(start)
	end, err := strconv.Atoi(queryForm.Get("end"))
	if err != nil {
		//	end = 0
	}
	input.Para.PField.End = int64(end)
	ugctype := queryForm.Get("type")

	input.Para.Type = ugctype
	keyword := queryForm.Get("keyword")
	if keyword == "" {
		keyword = "time"
	}
	input.Para.KeyWord = keyword
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

	psort := queryForm.Get("sort")
	if psort == "" {
		psort = "time"
	} else {
		input.Para.Sort = psort
	}

	unit := queryForm.Get("unit")
	input.Para.Unit = unit

	order := queryForm.Get("order")
	input.Para.Order = order
	//	ugctype, type_err := strconv.Atoi(queryForm.Get("type"))

	//get list
	//	err, list := new(models.TblUrgency).GetLast(input)
	err, list := new(urgency.TblUrgency).GetUrgencyC(&input.Para)
	if err != nil {
		panic(fmt.Sprintf("TblUrgency GetLast error:%s", err.Error()))
		return
	}
	//response
	public.Write(w, public.ErrOkErr, list)
}

type UrgencyGetInput struct {
	Para urgency.TblUrgencySearchPara
}
