/********获取文件威胁数(天)********/
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

type FTTController struct{}

var FTTObj = new(FTTController)

func (this *FTTController) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//check ssn
	//TODO
	//parse parament
	input := FTTGetInput{}

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
		fmt.Print(err)
	}
	input.Para.PField.End = int64(end)

	//get list
	err, list := new(vds.TblFTT).GetFTTrent(&input.Para)
	if err != nil {
		panic(fmt.Sprintf("TblUrgency GetLast error:%s", err.Error()))
		return
	}
	//response
	public.Write(w, public.ErrOkErr, list)
}

type FTTGetInput struct {
	Para vds.TblFTTSearchPara
}
