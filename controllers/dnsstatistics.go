/********DNS统计********/
package controllers

import (
	"apt-web-server/models"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

type DNSSController struct{}

var DNSSObj = new(DNSSController)

func (this *DNSSController) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	input := DNSSGetInput{}

	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		Write(w, ErrParamentErr, nil)
		return
	}
	paratype := queryForm.Get("type")
	input.Para.Type = paratype

	start, err := strconv.Atoi(queryForm.Get("start"))
	if err != nil {
		if paratype != "" {
			input.Para.PField.Start = time.Now().Unix()
		}
	} else {
		input.Para.PField.Start = int64(start)
	}

	end, err := strconv.Atoi(queryForm.Get("end"))
	if err != nil {
		//	end = 0
	}
	input.Para.PField.End = int64(end)

	paraip := queryForm.Get("ip")
	input.Para.Ip = paraip

	paradomain := queryForm.Get("domain")
	input.Para.Domain = paradomain

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

	unit := queryForm.Get("unit")
	if unit == "" {
		unit = "day"
	}
	input.Para.Unit = unit

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

	var list *models.TblDNSSData
	var list_ip *models.TblDNSSIpData
	var list_domain *models.TblDNSSDomainData
	switch input.Para.Type {
	case "ip":
		err, list_ip = new(models.TblDNSS).GetDNSSIp(&input.Para)
	case "domain":
		err, list_domain = new(models.TblDNSS).GetDNSSDomain(&input.Para)
	default:
		err, list = new(models.TblDNSS).GetDNSS(&input.Para)
	}
	if err != nil {
		panic(fmt.Sprintf("TblDNS GetLast error:%s", err.Error()))
		return
	}
	//response
	switch input.Para.Type {
	case "ip":
		Write(w, ErrOkErr, list_ip)
	case "domain":
		Write(w, ErrOkErr, list_domain)
	default:
		Write(w, ErrOkErr, list)
	}
}

type DNSSGetInput struct {
	Para models.TblDNSSSearchPara
}
