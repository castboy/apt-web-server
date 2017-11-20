// sniffplcy_select.go
//sniff 策略查询
package sniff

import (
	"apt-web-server_v2/controllers/public"
	"apt-web-server_v2/models/sniff"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type SniffselectController struct{}

var SniffselectObj = new(SniffselectController)

func (this *SniffselectController) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var select_s sniff.Sniffshow_st
	input := SniffselectInput{}

	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		public.Write(w, public.ErrParamentErr, nil)
		return
	}

	fmt.Println(queryForm)
	input.Para.Type = queryForm.Get("type")
	//fmt.Println(input)
	switch input.Para.Type {
	case "all":
		select_s = sniff.Sniff_select(input.Para)
		if select_s.Plcy_s == nil {
			public.Write(w, public.ErrOkErr, nil)
			return
		}
	case "limit":
		page, err := strconv.Atoi(queryForm.Get("page"))
		if err != nil {
			page = 0
		}
		count, err := strconv.Atoi(queryForm.Get("count"))
		if err != nil {
			count = 0
		}
		input.Para.Page = (page - 1) * count
		input.Para.Count = count
		input.Para.Lies = queryForm.Get("lies")
		input.Para.Orderby = queryForm.Get("orderby")
		select_s = sniff.Sniff_select(input.Para)
		if select_s.Plcy_s == nil {
			public.Write(w, public.ErrOkErr, nil)
			return
		}
	default:
		public.Write(w, public.ErrOkErr, nil)
		return
	}
	public.Write(w, public.ErrOkErr, select_s)
}

type SniffselectInput struct {
	Para sniff.Sniffselect_st
}
