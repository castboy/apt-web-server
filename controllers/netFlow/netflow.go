package netFlow

import (
	"apt-web-server_v2/controllers/public"
	"apt-web-server_v2/models/netFlow"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type NetFlowController struct{}

var NetFlowObj = new(NetFlowController)

func (this *NetFlowController) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//check ssn
	//TODO
	//parse parament
	var list *netFlow.TblNetFlowData
	ipTag, directionTag, protocolTag := 1, 10, 100

	input := NetFlowGetInput{}

	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		public.Write(w, public.ErrParamentErr, nil)
		return
	}
	//check parameter start,end,asset,direction,protocol
	start, err := strconv.Atoi(queryForm.Get("start"))
	if err != nil {
		//start = 0
	}
	input.Para.PField.Start = int64(start)

	end, err := strconv.Atoi(queryForm.Get("end"))
	if err != nil {
		//end = 0
	}
	input.Para.PField.End = int64(end)

	unit := queryForm.Get("unit")
	//	if unit == "" {
	//		input.Para.Unit = "minute"
	//	} else {
	input.Para.Unit = unit
	//	}

	asset := queryForm.Get("asset")
	if asset == "" {
		ipTag = 0
	}
	input.Para.AssetIP = asset

	direction := queryForm.Get("direction")
	if direction == "" {
		directionTag = 0
	}
	input.Para.Direction = direction

	protocol := queryForm.Get("protocol")
	if protocol == "" {
		protocolTag = 0
	}
	input.Para.Protocol = protocol
	//chose function to get list
	parameterType := ipTag + directionTag + protocolTag
	switch parameterType {
	case 1:
		err, list = new(netFlow.TblNetFlowIP).GetNetFlowIP(&input.Para)
	case 10:
		err, list = new(netFlow.TblNetFlowD).GetNetFlowD(&input.Para)
	case 100:
		err, list = new(netFlow.TblNetFlowP).GetNetFlowP(&input.Para)
	case 11:
		err, list = new(netFlow.TblNetFlowIPD).GetNetFlowIPD(&input.Para)
	case 101:
		err, list = new(netFlow.TblNetFlowIPP).GetNetFlowIPP(&input.Para)
	case 110:
		err, list = new(netFlow.TblNetFlowDP).GetNetFlowDP(&input.Para)
	default:
		err, list = new(netFlow.TblNetFlow).GetNetFlow(&input.Para)
	}

	if err != nil {
		panic(fmt.Sprintf("TblNetFlow GetLast error:%s", err.Error()))
		return
	}
	//response
	public.Write(w, public.ErrOkErr, list)
}

type NetFlowGetInput struct {
	Para netFlow.TblNetFlowSearchPara
}
