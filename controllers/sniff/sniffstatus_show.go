// sniffstatus_show.go
//sniff 查询下发状态
package sniff

import (
	"apt-web-server_v2/controllers/public"
	"apt-web-server_v2/models/sniff"
	"fmt"
	"net/http"
	"net/url"
	//"strconv"

	"github.com/julienschmidt/httprouter"
)

type SniffstatusController struct{}

var SniffstatusObj = new(SniffstatusController)

func (this *SniffstatusController) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	input := SniffstatusInput{}

	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		public.Write(w, public.ErrParamentErr, nil)
		return
	}

	fmt.Println(queryForm)
	input.Para.Type = queryForm.Get("type")
	//fmt.Println(input)
	switch input.Para.Type {
	case "issued":
		input.Para = sniff.Sniff_status(input.Para)
		public.Write(w, public.ErrOkErr, input.Para.Msg)
		return
	default:
		public.Write(w, public.ErrOkErr, "err")
		return
	}
	public.Write(w, public.ErrOkErr, "err")
}

type SniffstatusInput struct {
	Para sniff.Sniffstatus_st
}
