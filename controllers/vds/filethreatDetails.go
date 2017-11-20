/********获取文件威胁详情********/
package vds

import (
	"apt-web-server_v2/controllers/public"
	"apt-web-server_v2/models/vds"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type FTDController struct{}

var FTDObj = new(FTDController)

func (this *FTDController) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//check ssn
	//TODO
	//parse parament
	input := FTDGetInput{}

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
	id, err := strconv.Atoi(queryForm.Get("id"))
	if err != nil {
		id = 0
	}
	input.Para.Id = int32(id)
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
	if keyword == "" || keyword == "Time" {
		keyword = "time"
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
	var list *vds.TblFTDData
	var dupresult string
	if input.Para.Tage == "offlinedup" || input.Para.Tage == "duplicate" {
		err, list = new(vds.TblFTD).GetFTDetailsDUP(&input.Para)
		if err != nil {
			panic(fmt.Sprintf("TblFTD GetLast error:%s", err.Error()))
			dupresult = "fail"
			return
		}
		if mergeTag == "merge" {
			dupresult = "ok"
		}
	} else {
		err, list = new(vds.TblFTD).GetFTDetails(&input.Para)
		if err != nil {
			panic(fmt.Sprintf("TblFTD GetLast error:%s", err.Error()))
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

type FTDGetInput struct {
	Para vds.TblFTDSearchPara
}
