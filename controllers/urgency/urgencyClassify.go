/********获取紧急事件分类数********/
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

type UgcClassifyController struct{}

var UgcClassifyObj = new(UgcClassifyController)

func (this *UgcClassifyController) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//check ssn
	//TODO
	//parse parament
	input := UgcClassifyGetInput{}

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

	unit := queryForm.Get("unit")
	input.Para.Unit = unit
	order := queryForm.Get("order")
	if order == "" {
		order = "desc"
	}
	input.Para.Order = order
	/*add end*/
	input.Para.LastCount = int32(lastCount)
	//get list
	//	err, list := new(models.TblUrgency).GetLast(input)
	var list *urgency.TblUgcClassifyData
	var list_c *urgency.TblUgcCCountData
	if input.Para.Unit != "" {
		err, list_c = new(urgency.TblUgcClassify).GetUgcCCount(&input.Para)
		if err != nil {
			panic(fmt.Sprintf("TblUrgency GetLast error:%s", err.Error()))
			return
		}
	} else {
		err, list = new(urgency.TblUgcClassify).GetUgcClassify(&input.Para)
		if err != nil {
			panic(fmt.Sprintf("TblUrgency GetLast error:%s", err.Error()))
			return
		}
	}

	//response
	if input.Para.Unit != "" {
		public.Write(w, public.ErrOkErr, list_c)
	} else {
		public.Write(w, public.ErrOkErr, list)
	}

}

type UgcClassifyGetInput struct {
	Para urgency.TblUgcClassifySearchPara
}
