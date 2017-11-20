/********最新紧急事件********/
package urgency

import (
	"apt-web-server_v2/controllers/public"
	"apt-web-server_v2/models/urgency"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type UgcLController struct{}

var UgcLObj = new(UgcLController)

func (this *UgcLController) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//check ssn
	//TODO
	//parse parament
	input := UgcLGetInput{}

	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		public.Write(w, public.ErrParamentErr, nil)
		return
	}

	lastCount, err := strconv.Atoi(queryForm.Get("lastcount"))
	if err != nil {
		lastCount = 10
	}
	input.Para.LastCount = int32(lastCount)

	err, list := new(urgency.TblUgcL).GetUrgencyLatest(&input.Para)
	if err != nil {
		panic(fmt.Sprintf("TblUrgency GetLast error:%s", err.Error()))
		return
	}

	public.Write(w, public.ErrOkErr, list)
}

type UgcLGetInput struct {
	Para urgency.TblUgcLSearchPara
}
