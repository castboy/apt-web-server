// attackstatistics.go
package index

import (
	"apt-web-server_v2/controllers/public"
	"apt-web-server_v2/models/index"
	//"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type AttackStatisticsController struct{}

var AttackStatisticsObj = new(AttackStatisticsController)

func (this *AttackStatisticsController) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var output index.StatisticsAttackOut_st
	input := AttackStatisticsInput{}

	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		public.Write(w, public.ErrParamentErr, nil)
		return
	}
	if len(queryForm) != 3 {
		public.Write(w, public.ErrOkErr, "接口填写有误不能识别")
		return
	}
	//fmt.Println(queryForm)
	input.Para.Method = queryForm.Get("method")
	//fmt.Println(input.Para.Method)
	input.Para.Attribute = queryForm.Get("attribute")
	//fmt.Println(input.Para.Attribute)
	strength, err := strconv.Atoi(queryForm.Get("strength"))
	if err != nil {
		public.Write(w, public.ErrOkErr, nil)
		return
	}
	input.Para.Strength = strength
	//fmt.Println(input.Para.Strength)
	switch input.Para.Method {
	case "percent":
		output = index.Statistics_Percent(input.Para)
	case "count":
		output = index.Statistics_Count(input.Para)
	case "top5":
		output = index.Statistics_Top5(input.Para)
	case "top10":
		output = index.Statistics_Top10(input.Para)
	default:
		public.Write(w, public.ErrOkErr, nil)
		return
	}
	if output.Count == 0 {
		public.Write(w, public.ErrOkErr, nil)
		return
	}
	public.Write(w, public.ErrOkErr, output)
}

type AttackStatisticsInput struct {
	Para index.StatisticsAttackIn_st
}
