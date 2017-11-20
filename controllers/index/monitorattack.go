// monitorattack.go
package index

import (
	"apt-web-server_v2/controllers/public"
	"apt-web-server_v2/models/index"
	//"fmt"
	"net/http"
	"net/url"
	//"strconv"

	"github.com/julienschmidt/httprouter"
)

type AttackMonitorController struct{}

var AttackMonitorObj = new(AttackMonitorController)

func (this *AttackMonitorController) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var output index.MonitorAttack_st
	input := AttackMonitorInput{}

	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		public.Write(w, public.ErrParamentErr, nil)
		return
	}
	if len(queryForm) != 2 {
		public.Write(w, public.ErrOkErr, "接口填写有误不能识别")
		return
	}
	//fmt.Println(queryForm)
	input.Para.Method = queryForm.Get("method")
	//fmt.Println(input.Para.Method)
	input.Para.Attribute = queryForm.Get("attribute")
	//fmt.Println(input.Para.Attribute)
	switch input.Para.Method {
	case "look_up":
		output = index.Monitor_Lookup(input.Para)
	case "count":
		output = index.Monitor_Count(input.Para)
	default:
		public.Write(w, public.ErrOkErr, nil)
		return
	}
	if output.Count == 0 {
		public.Write(w, public.ErrOkErr, nil)
		return
	}
	//fmt.Println(output)
	public.Write(w, public.ErrOkErr, output)
}

type AttackMonitorInput struct {
	Para index.StatisticsAttackIn_st
}
