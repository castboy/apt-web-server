/********获取http流量攻击数(天)********/
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

type HFATController struct{}

var HFATObj = new(HFATController)

func (this *HFATController) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//check ssn
	//TODO
	//parse parament
	input := HFATGetInput{}

	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		public.Write(w, public.ErrParamentErr, nil)
		return
	}

	start, err := strconv.Atoi(queryForm.Get("start"))
	if err != nil {
		fmt.Println(err)
	}
	input.Para.PField.Start = int64(start)
	end, err := strconv.Atoi(queryForm.Get("end"))
	if err != nil {
		fmt.Println(err)
	}
	input.Para.PField.End = int64(end)
	//get list
	//	err, list := new(models.TblUrgency).GetLast(input)
	err, list := new(waf.TblHFAT).GetHFATrend(&input.Para)
	if err != nil {
		panic(fmt.Sprintf("TblUrgency GetLast error:%s", err.Error()))
		return
	}
	//response
	public.Write(w, public.ErrOkErr, list)
}

type HFATGetInput struct {
	Para waf.TblHFATSearchPara
}
