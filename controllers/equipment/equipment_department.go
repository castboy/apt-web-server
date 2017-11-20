// equipment_department.go
package equipment

import (
	"apt-web-server_v2/controllers/public"
	"apt-web-server_v2/models/equipment"
	//"fmt"
	"encoding/base64"
	"net/http"
	"net/url"
	//"reflect"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type EquipmentdepartmentController struct{}

var EquipmentdepartmentObj = new(EquipmentdepartmentController)

func (this *EquipmentdepartmentController) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		public.Write(w, public.ErrParamentErr, nil)
		return
	}
	operation := queryForm.Get("operation")
	if operation == "" {
		public.Write(w, public.ErrOkErr, "??????????")
		return
	}
	switch operation {
	case "create":
		if len(queryForm) != 2 {
			public.Write(w, public.ErrOkErr, "??????????")
			return
		}
		departmentcreate(w, queryForm)
		return
	case "delete":
		if len(queryForm) != 2 {
			public.Write(w, public.ErrOkErr, "??????????")
			return
		}
		departmentdelete(w, queryForm)
		return
	case "update":
		if len(queryForm) != 3 {
			public.Write(w, public.ErrOkErr, "??????????")
			return
		}
		departmentupdate(w, queryForm)
		return
	case "select":
		//fmt.Println(reflect.TypeOf(queryForm))
		if len(queryForm) != 4 {
			public.Write(w, public.ErrOkErr, "??????????")
			return
		}
		departmentselect(w, queryForm)
		return
	case "ipselect":
		if len(queryForm) != 2 {
			public.Write(w, public.ErrOkErr, "??????????")
			return
		}
		ipselect(w, queryForm)
		return
	default:
		public.Write(w, public.ErrOkErr, "??????????")
		return
	}
}
func departmentupdate(w http.ResponseWriter, queryForm url.Values) {
	var departmentname string
	input := createdepartmentInput{}
	departmentname = queryForm.Get("departmentname")
	name, err := base64.StdEncoding.DecodeString(departmentname)
	if nil != err {
		public.Write(w, public.ErrOkErr, "err")
		return
	}
	departmentId, err := strconv.Atoi(queryForm.Get("departmentId"))
	if err != nil {
		departmentId = 0
	}
	input.Para.DepartmentId = departmentId
	input.Para.Departmentname = string(name)
	rst := equipment.Department_update(input.Para)
	public.Write(w, public.ErrOkErr, rst)
	return
}
func departmentcreate(w http.ResponseWriter, queryForm url.Values) {
	var departmentname string
	departmentname = queryForm.Get("departmentname")
	name, err := base64.StdEncoding.DecodeString(departmentname)
	if nil != err {
		public.Write(w, public.ErrOkErr, "departmentname nead base64")
		return
	}
	departmentname = string(name)
	rst := equipment.Department_create(departmentname)
	public.Write(w, public.ErrOkErr, rst)
	return
}
func departmentdelete(w http.ResponseWriter, queryForm url.Values) {
	var departmentId int
	departmentId, err := strconv.Atoi(queryForm.Get("departmentId"))
	if err != nil {
		departmentId = 0
	}
	rst := equipment.Department_delete(departmentId)
	public.Write(w, public.ErrOkErr, rst)
	return
}
func ipselect(w http.ResponseWriter, queryForm url.Values) {
	var departmentId int
	departmentId, err := strconv.Atoi(queryForm.Get("departmentId"))
	if err != nil {
		departmentId = 0
	}
	select_s := equipment.Ip_select(departmentId)
	if select_s.Ip == nil {
		public.Write(w, public.ErrOkErr, nil)
		return
	}
	public.Write(w, public.ErrOkErr, select_s)
	return
}
func departmentselect(w http.ResponseWriter, queryForm url.Values) {
	var departmentId int
	input := equipmentselectInput{}
	departmentId, err := strconv.Atoi(queryForm.Get("departmentId"))
	if err != nil {
		departmentId = 0
	}
	page, err := strconv.Atoi(queryForm.Get("page"))
	if err != nil {
		page = 0
	}
	count, err := strconv.Atoi(queryForm.Get("count"))
	if err != nil {
		count = 0
	}
	input.Para.DepartmentId = departmentId
	input.Para.Page = (page - 1) * count
	input.Para.Count = count
	if count == 0 && departmentId != 0 {
		public.Write(w, public.ErrOkErr, "err")
		return
	}

	select_s := equipment.Department_select(input.Para)
	if select_s.Info == nil {
		public.Write(w, public.ErrOkErr, nil)
		return
	}
	public.Write(w, public.ErrOkErr, select_s)
	return
}

type createdepartmentInput struct {
	Para equipment.Departmentip_st
}
