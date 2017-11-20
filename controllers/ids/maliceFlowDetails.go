/********获取恶意流量详情********/
package ids

import (
	"apt-web-server_v2/controllers/public"
	"apt-web-server_v2/models/ids"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type MFDController struct{}

var MFDObj = new(MFDController)

func (this *MFDController) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//check ssn
	//TODO
	//parse parament
	input := MFDGetInput{}

	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		public.Write(w, public.ErrParamentErr, nil)
		return
	}

	lastCount, err := strconv.Atoi(queryForm.Get("lastCount"))
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
	//	ugctype, type_err := strconv.Atoi(queryForm.Get("type"))
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
	id, err := strconv.Atoi(queryForm.Get("id"))
	if err != nil {
		id = 0
	}
	input.Para.Id = int32(id)
	count, err := strconv.Atoi(queryForm.Get("count"))
	if err != nil {
		//count = 1000
	}
	input.Para.Count = int32(count)

	sort := queryForm.Get("sort")
	if sort == "" || sort == "Time" {
		sort = "time"
	}
	input.Para.Sort = sort

	order := queryForm.Get("order")
	input.Para.Order = order
	/*add end*/
	input.Para.LastCount = int32(lastCount)
	//get list
	//	err, list := new(models.TblUrgency).GetLast(input)
	err, list := new(ids.TblMFD).GetMFDetails(&input.Para)
	if err != nil {
		panic(fmt.Sprintf("TblMFD GetLast error:%s", err.Error()))
		return
	}
	//response
	public.Write(w, public.ErrOkErr, list)
}

type MFDGetInput struct {
	Para ids.TblMFDSearchPara
}
