/********获取HTTP流量攻击详情********/
package waf

import (
	"apt-web-server_v2/controllers/public"
	"apt-web-server_v2/models/waf"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type HFADController struct{}

var HFADObj = new(HFADController)

func (this *HFADController) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//check ssn
	//TODO
	//parse parament
	input := HFADGetInput{}

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
	id, err := strconv.Atoi(queryForm.Get("id"))
	if err != nil {
		id = 0
	}
	input.Para.Id = int32(id)
	ugctype := queryForm.Get("type")
	if ugctype == "" {
		input.Para.Type = ""
	} else {
		input.Para.Type = ugctype
	}
	keyword := queryForm.Get("keyword")
	if keyword == "" {
		keyword = ""
	}
	input.Para.KeyWord = keyword
	tage := queryForm.Get("tage")
	if tage == "" {
		input.Para.Tage = "online"
	} else {
		input.Para.Tage = tage
	}
	taskName := queryForm.Get("taskname")
	input.Para.TaskName = taskName
	createTime, err := strconv.Atoi(queryForm.Get("createtime"))
	if err != nil {
		fmt.Println("the task create time is null")
	}
	input.Para.CreateTime = int64(createTime)
	mergeTag := queryForm.Get("mergetag")
	input.Para.MergeTag = mergeTag
	dupTag := queryForm.Get("duptag")
	input.Para.DupTag = dupTag

	page, err := strconv.Atoi(queryForm.Get("page"))
	if err != nil {
		page = 1
	}
	input.Para.Page = int32(page)
	count, err := strconv.Atoi(queryForm.Get("count"))
	if err != nil {
		count = 10
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
	var list *waf.TblHFADData
	var dupresult string
	if input.Para.Tage == "offlinedup" || input.Para.Tage == "duplicate" {
		err, list = new(waf.TblHFAD).GetHFADetailsDUP(&input.Para)
		if err != nil {
			panic(fmt.Sprintf("TblHFAD GetLast error:%s", err.Error()))
			dupresult = "fail"
			return
		}
		if mergeTag == "merge" {
			dupresult = "ok"
		}
	} else {
		err, list = new(waf.TblHFAD).GetHFADetails(&input.Para)
		if err != nil {
			panic(fmt.Sprintf("TblHFAD GetLast error:%s", err.Error()))
			return
		}
	}
	//response
	if tage == "offlinedup" && mergeTag == "merge" {
		public.Write(w, public.ErrOkErr, dupresult)
	} else {
		public.Write(w, public.ErrOkErr, list)
	}
}

type HFADGetInput struct {
	Para waf.TblHFADSearchPara
}
