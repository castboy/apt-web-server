/********获取http流量攻击分类数(天)********/
package controllers

import (
	"apt-web-server/models"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type HFACController struct{}

var HFACObj = new(HFACController)

func (this *HFACController) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//check ssn
	//TODO
	//parse parament
	input := HFACGetInput{}

	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		Write(w, ErrParamentErr, nil)
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
	err, list := new(models.TblHFAC).GetHFAClassify(&input.Para)
	if err != nil {
		panic(fmt.Sprintf("TblUrgency GetLast error:%s", err.Error()))
		return
	}
	//response
	Write(w, ErrOkErr, list)
}

type HFACGetInput struct {
	Para models.TblHFACSearchPara
}

