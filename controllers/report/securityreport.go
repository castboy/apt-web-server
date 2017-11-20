// securityreport
// monitorattack.go
package report

import (
	"apt-web-server_v2/controllers/public"
	"apt-web-server_v2/models/report"
	//"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type SecurityController struct{}

var SecurityObj = new(SecurityController)

func (this *SecurityController) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var output report.Security_st
	input := SecurityInput{}

	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		public.Write(w, public.ErrParamentErr, nil)
		return
	}
	if len(queryForm) != 2 {
		public.Write(w, public.ErrOkErr, "接口填写有误不能识别")
		return
	}
	Method := queryForm.Get("method")
	switch Method {
	case "all":
		//select all
		output = report.Security_Reportall()
		//fmt.Println(output)
	case "condition":
		Attribute := queryForm.Get("attribute")
		switch strings.Contains(Attribute, ",") {
		case true:
			spl := strings.Split(Attribute, ",")
			for i, _ := range spl {
				if spl[i] != "" {
					input.Para.Info = append(input.Para.Info, spl[i])
					//input.Para.Info[i] = spl[i]
				}
			}
		case false:
			input.Para.Info = append(input.Para.Info, Attribute)
			//input.Para.Info[0] = Attribute
		default:
			public.Write(w, public.ErrOkErr, "err")
			return
		}
		//fmt.Println(input.Para.Info)
		//select condition
		output = report.Security_Reportcondition(input.Para)
	default:
		public.Write(w, public.ErrOkErr, nil)
		return
	}
	if output.Score == -1 {
		public.Write(w, public.ErrOkErr, nil)
		return
	}
	//fmt.Println(output)

	public.Write(w, public.ErrOkErr, output)
}

type SecurityInput struct {
	Para report.Attribute_st
}
