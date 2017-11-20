// sniffplcy_delete.go
//sniff 策略删除
package sniff

import (
	"apt-web-server_v2/controllers/public"
	"apt-web-server_v2/models/sniff"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	//"strconv"

	"github.com/julienschmidt/httprouter"
)

type SniffdeleteController struct{}

var SniffdeleteObj = new(SniffdeleteController)

func (this *SniffdeleteController) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//check ssn
	//TODO
	//parse parament
	//input := models.TblUrgencySearchPara{}

	input := SniffdeleteInput{}

	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		public.Write(w, public.ErrParamentErr, nil)
		return
	}

	fmt.Println(queryForm)
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
		rst := sniff.Sniff_delete(input.Para)
		if rst == false {
			public.Write(w, public.ErrOkErr, "err")
			return
		}
	case false:
		input.Para.Plcy_name[0] = Plcys_name
		rst := sniff.Sniff_delete(input.Para)
		if rst == false {
			public.Write(w, public.ErrOkErr, "err")
			return
		}
	default:
		public.Write(w, public.ErrOkErr, "err")
		return
	}
	fmt.Println(input)
	rst := sniff.Manage()
	if rst == false {
		fmt.Println("manage plcy err")
		public.Write(w, public.ErrOkErr, "err")
		return
	}
	public.Write(w, public.ErrOkErr, "ok")
}

type SniffdeleteInput struct {
	Para sniff.Sniffissued_st
}
