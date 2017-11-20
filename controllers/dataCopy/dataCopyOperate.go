/********数据迁移任务********/
package dataCopy

import (
	"apt-web-server_v2/controllers/public"
	"apt-web-server_v2/models/dataCopy"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type DCTController struct{}

var DCTObj = new(DCTController)

func (this *DCTController) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//check ssn
	//TODO
	//parse parament
	var list *dataCopy.DCTOperateResult
	input := DCTGetInput{}

	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		public.Write(w, public.ErrParamentErr, nil)
		return
	}
	location := queryForm.Get("location")
	input.Para.Location = location

	date, err := strconv.Atoi(queryForm.Get("date"))
	if err != nil {
		//time = 0
	}
	input.Para.Date = int64(date)

	start := queryForm.Get("start")
	if err != nil {
		//start = 0
	}
	input.Para.Start = start
	end := queryForm.Get("end")
	if err != nil {
		//end = 0
	}
	input.Para.End = end
	diskpath := queryForm.Get("diskpath")
	if err != nil {
		//start = 0
	}
	input.Para.DiskPath = diskpath
	command := queryForm.Get("cmd")
	fmt.Println(input)
	switch command {
	case "add":
		err, list = new(dataCopy.TblDCT).CreateDCT(&input.Para)
	case "delete":
		err, list = new(dataCopy.TblDCT).DeleteDCT(&input.Para)
	case "start":
		err, list = new(dataCopy.TblDCT).StartDCT(&input.Para)
	case "getdisk":
		err, list = new(dataCopy.TblDCT).GetDiskDCT(&input.Para)
	case "remove":
		err, list = new(dataCopy.TblDCT).RemoveDiskDCT(&input.Para)
	}

	if err != nil {
		panic(fmt.Sprintf("The operate do %s error:%s", command, err.Error()))
		return
	}
	//response
	//Write(w, ErrOkErr, queryForm)
	public.Write(w, public.ErrOkErr, list)
}

type DCTGetInput struct {
	Para dataCopy.DCTOperatePara
}
