// sniffplcy_issued.go
//sniff 策略下发
package sniff

import (
	"apt-web-server_v2/controllers/public"
	"apt-web-server_v2/models/modelsPublic"
	"apt-web-server_v2/models/sniff"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	//"strconv"

	"github.com/julienschmidt/httprouter"
)

type SniffissuedController struct{}

var SniffissuedObj = new(SniffissuedController)

func (this *SniffissuedController) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//check ssn
	//TODO
	//parse parament
	//input := models.TblUrgencySearchPara{}

	input := SniffissuedInput{}

	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		public.Write(w, public.ErrParamentErr, nil)
		return
	}

	//fmt.Println(queryForm)
	Plcys_name := queryForm.Get("plcy_name")
	if len(Plcys_name) == 0 {
		public.Write(w, public.ErrOkErr, "err")
		return
	}
	switch strings.Contains(Plcys_name, ",") {
	case true:
		plcy_name := strings.Split(Plcys_name, ",")
		for i, _ := range plcy_name {
			if plcy_name[i] != "" {
				input.Para.Plcy_name[i] = plcy_name[i]
			}
		}
		rst := sniff.Sniff_issued(input.Para)
		if rst == false {
			public.Write(w, public.ErrOkErr, "err")
			return
		}
	case false:
		input.Para.Plcy_name[0] = Plcys_name
		rst := sniff.Sniff_issued(input.Para)
		if rst == false {
			public.Write(w, public.ErrOkErr, "err")
			return
		}
	default:
		public.Write(w, public.ErrOkErr, "err")
		return
	}
	//fmt.Println(input)
	rst := sniff.Manage()
	if rst == false {
		fmt.Println("manage plcy err")
		public.Write(w, public.ErrOkErr, "err")
		return
	}
	rst = modelsPublic.Update_mysql("update sniff_plcy_ui set issued_status=1")
	rst = modelsPublic.Update_mysql("update sniff_plcy_issued set status=1")
	if rst == false {
		public.Write(w, public.ErrOkErr, "err")
		return
	}

	public.Write(w, public.ErrOkErr, "ok")
}

type SniffissuedInput struct {
	Para sniff.Sniffissued_st
}
