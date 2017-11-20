/********获取恶意流量分类数(天)********/
package ids

import (
	"apt-web-server_v2/controllers/public"
	"apt-web-server_v2/models/ids"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type MFTController struct{}

var MFTObj = new(MFTController)

func (this *MFTController) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//check ssn
	//TODO
	//parse parament
	input := MFTGetInput{}

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
	err, list := new(ids.TblMFT).GetMFTrend(&input.Para)
	if err != nil {
		panic(fmt.Sprintf("TblUrgency GetLast error:%s", err.Error()))
		return
	}
	//response
	public.Write(w, public.ErrOkErr, list)
}

type MFTGetInput struct {
	Para ids.TblMFTSearchPara
}
