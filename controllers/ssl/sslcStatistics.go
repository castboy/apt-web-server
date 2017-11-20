/********SSLCert统计********/
package ssl

import (
	"apt-web-server_v2/controllers/public"
	"apt-web-server_v2/models/ssl"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

type SSLCSController struct{}

var SSLCSObj = new(SSLCSController)

func (this *SSLCSController) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	input := SSLCSGetInput{}

	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		public.Write(w, public.ErrParamentErr, nil)
		return
	}
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

	start, err := strconv.Atoi(queryForm.Get("start"))
	if err != nil {
		if paratype != "" {
			input.Para.PField.Start = time.Now().Unix() - times
		}
	} else {
		input.Para.PField.Start = int64(start)
	}

	end, err := strconv.Atoi(queryForm.Get("end"))
	if err != nil {
		input.Para.PField.End = time.Now().Unix()
	} else {
		input.Para.PField.End = int64(end)
	}

	paraverify, err := strconv.Atoi(queryForm.Get("verify"))
	if err != nil {
		if paratype != "" {
			input.Para.Verify = 1
		}
	} else {
		input.Para.Verify = paraverify
	}

	paraip := queryForm.Get("ip")
	input.Para.Ip = paraip

	paraComnName := queryForm.Get("commnname")
	input.Para.ComnName = paraComnName

	paraUnitName := queryForm.Get("unitname")
	input.Para.UnitName = paraUnitName

	paraSerialNum := queryForm.Get("serialnum")
	input.Para.SerialNum = paraSerialNum

	keyword := queryForm.Get("keyword")
	if keyword == "" {
		keyword = "time"
	}
	input.Para.KeyWord = keyword

	page, err := strconv.Atoi(queryForm.Get("page"))
	if err != nil {
		page = 1
	}
	input.Para.Page = int32(page)
	count, err := strconv.Atoi(queryForm.Get("count"))
	if err != nil {
		//count = 1000
	}
	input.Para.Count = int32(count)

	sort := queryForm.Get("sort")
	if sort == "" {
		sort = "count"
	}
	input.Para.Sort = sort

	order := queryForm.Get("order")
	if order == "" {
		order = "ASC"
	}
	input.Para.Order = order

	lastCount, err := strconv.Atoi(queryForm.Get("lastcount"))
	if err != nil {
		if input.Para.Type != "" {
			lastCount = 5
		}
	}
	input.Para.LastCount = int32(lastCount)

	var list *ssl.TblSSLCSData
	var list_ip *ssl.TblSSLCSIpData
	var list_cert *ssl.TblSSLCSCertData
	switch input.Para.Type {
	case "ip":
		err, list_ip = new(ssl.TblSSLCS).GetSSLCSIp(&input.Para)
	case "cert":
		err, list_cert = new(ssl.TblSSLCS).GetSSLCSCert(&input.Para)
	default:
		err, list = new(ssl.TblSSLCS).GetSSLCS(&input.Para)
	}
	if err != nil {
		panic(fmt.Sprintf("TblDNS GetLast error:%s", err.Error()))
		return
	}
	//response
	switch input.Para.Type {
	case "ip":
		public.Write(w, public.ErrOkErr, list_ip)
	case "cert":
		public.Write(w, public.ErrOkErr, list_cert)
	default:
		public.Write(w, public.ErrOkErr, list)
	}
}

type SSLCSGetInput struct {
	Para ssl.TblSSLCSSearchPara
}
