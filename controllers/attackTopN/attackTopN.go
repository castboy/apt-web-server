/********攻击数topN统计********/
package attackTopN

import (
	"apt-web-server_v2/controllers/public"
	"apt-web-server_v2/models/attackTopN"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

type AttackController struct{}

var AttackObj = new(AttackController)

func (this *AttackController) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//check ssn
	//TODO
	//parse parament
	input := AttackGetInput{}

	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		public.Write(w, public.ErrParamentErr, nil)
		return
	}

	/*add all parameter*/
	paratype := queryForm.Get("type")
	input.Para.Type = paratype

	unit := queryForm.Get("unit")
	if unit == "" {
		unit = "day"
	}
	input.Para.Unit = unit

	var times int64
	switch unit {
	case "week":
		times = 60 * 60 * 24 * 7
	case "month":
		times = 60 * 60 * 24 * 30
	default:
		times = 0
	}
	times = 60 * 60 * 24 * 30
	start, err := strconv.Atoi(queryForm.Get("start"))
	if err != nil {
		//if paratype != "" {
		input.Para.Start = time.Now().Unix() - times
		//}
	} else {
		input.Para.Start = int64(start)
	}

	end, err := strconv.Atoi(queryForm.Get("end"))
	if err != nil {
		input.Para.End = time.Now().Unix()
	} else {
		input.Para.End = int64(end)
	}

	paraCount, err := strconv.Atoi(queryForm.Get("count"))
	if err != nil {
		input.Para.Count = 3
	} else {
		input.Para.Count = int(paraCount)
	}

	//order := queryForm.Get("order")
	//input.Para.Order = order
	/*add end*/
	//input.Para.LastCount = int32(lastCount)
	//get list
	//	err, list := new(models.TblUrgency).GetLast(input)

	err, list := new(attackTopN.AttackCount).GetAttackTopN(&input.Para)
	if err != nil {
		panic(fmt.Sprintf("TblUrgency GetLast error:%s", err.Error()))
		return
	}

	//response
	public.Write(w, public.ErrOkErr, list)
}

type AttackGetInput struct {
	Para attackTopN.AttackSearchPara
}
