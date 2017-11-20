// equipment_manage.go
package equipment

import (
	"apt-web-server_v2/controllers/public"
	"apt-web-server_v2/models/equipment"
	"apt-web-server_v2/models/modelsPublic"
	"apt-web-server_v2/modules/mlog"
	"fmt"
	"net/http"
	"net/url"
	//"reflect"
	"encoding/base64"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type EquipmentmanageController struct{}

var EquipmentmanageObj = new(EquipmentmanageController)

func (this *EquipmentmanageController) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		public.Write(w, public.ErrParamentErr, nil)
		return
	}
	operation := queryForm.Get("operation")
	if operation == "" {
		public.Write(w, public.ErrOkErr, "接口填写有误不能识别")
		return
	}
	switch operation {
	case "scan":
		if len(queryForm) != 1 {
			public.Write(w, public.ErrOkErr, "接口填写有误不能识别")
			return
		}
		rst := equipment.Equipment_scan()
		if rst == false {
			public.Write(w, public.ErrOkErr, "接口填写有误不能识别")
			return
		}
		public.Write(w, public.ErrOkErr, "ok")
		return
	case "create":
		if len(queryForm) != 5 {
			public.Write(w, public.ErrOkErr, "接口填写有误不能识别")
			return
		}
		equipmentcreate(w, queryForm)
		return
	case "delete":
		if len(queryForm) != 2 {
			public.Write(w, public.ErrOkErr, "接口填写有误不能识别")
			return
		}
		equipmentdelete(w, queryForm)
		return
	case "update":
		if len(queryForm) != 5 {
			public.Write(w, public.ErrOkErr, "接口填写有误不能识别")
			return
		}
		equipmentupdate(w, queryForm)
		return
	case "select":
		//fmt.Println(reflect.TypeOf(queryForm))
		if len(queryForm) != 8 {
			public.Write(w, public.ErrOkErr, "接口填写有误不能识别")
			return
		}
		equipmentselect(w, queryForm)
		return
	case "scanflag":
		if len(queryForm) != 1 {
			public.Write(w, public.ErrOkErr, "接口填写有误不能识别")
			return
		}
		rst := equipment.Equipment_scanflag()
		if rst == false {
			public.Write(w, public.ErrOkErr, 0)
			return
		}
		public.Write(w, public.ErrOkErr, 1)
		return
	case "detailinfo":
		if len(queryForm) != 4 {
			public.Write(w, public.ErrOkErr, "接口填写有误不能识别")
			return
		}
		detailinfoselect(w, queryForm)
		return
	default:
		public.Write(w, public.ErrOkErr, "接口填写有误不能识别")
		return
	}
}
func detailinfoselect(w http.ResponseWriter, queryForm url.Values) {
	input := equipmentselectInput{}
	input.Para.Ip = queryForm.Get("ip")
	page, err := strconv.Atoi(queryForm.Get("page"))
	if err != nil {
		page = 0
	}
	count, err := strconv.Atoi(queryForm.Get("count"))
	if err != nil {
		count = 0
	}
	input.Para.Page = (page - 1) * count
	input.Para.Count = count
	if count == 0 {
		public.Write(w, public.ErrOkErr, "err")
		return
	}
	select_s := equipment.Detailinfo_select(input.Para)
	if select_s.Detailinfo == nil {
		public.Write(w, public.ErrOkErr, nil)
		return
	}
	public.Write(w, public.ErrOkErr, select_s)
	return
}
func equipmentselect(w http.ResponseWriter, queryForm url.Values) {
	input := equipmentselectInput{}
	var select_s equipment.Equipmentshow_st
	input.Para.Ip = queryForm.Get("ip")
	input.Para.Os_type = queryForm.Get("os_type")
	input.Para.Lies = queryForm.Get("line")
	if input.Para.Lies == "ip" {
		input.Para.Lies = "time"
	}
	input.Para.Orderby = queryForm.Get("orderby")
	departmentId, err := strconv.Atoi(queryForm.Get("departmentId"))
	if err != nil {
		departmentId = 0
	}
	input.Para.DepartmentId = departmentId
	page, err := strconv.Atoi(queryForm.Get("page"))
	if err != nil {
		page = 0
	}
	count, err := strconv.Atoi(queryForm.Get("count"))
	if err != nil {
		count = 0
	}
	input.Para.Page = (page - 1) * count
	input.Para.Count = count
	if count == 0 {
		public.Write(w, public.ErrOkErr, "err")
		return
	}
	select_s = equipment.Equipment_select(input.Para)
	if select_s.Info == nil {
		public.Write(w, public.ErrOkErr, nil)
		return
	}
	public.Write(w, public.ErrOkErr, select_s)
	return
}
func equipmentcreate(w http.ResponseWriter, queryForm url.Values) {
	input := equipmentcreateInput{}
	var count int

	input.Para.Ip = queryForm.Get("ip")
	input.Para.Os_type = queryForm.Get("os_type")
	input.Para.Alias = queryForm.Get("alias")
	alias, err := base64.StdEncoding.DecodeString(input.Para.Alias)
	if nil != err {
		public.Write(w, public.ErrOkErr, "其他失败原因")
		mlog.Debug("##equipment_manage:", "alias nead base64")
		return
	}
	input.Para.Alias = string(alias)
	departmentId, err := strconv.Atoi(queryForm.Get("departmentId"))
	if err != nil {
		departmentId = 0
	}
	input.Para.DepartmentId = departmentId

	//fmt.Println(input)
	selete_cmd := fmt.Sprintf(`select COUNT(ip) from equipment where ip = '%s' and 
	delete_flag = 0`, input.Para.Ip)
	rows := modelsPublic.Select_mysql(selete_cmd)
	if rows == nil {
		public.Write(w, public.ErrOkErr, "其他失败原因")
		return
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			public.Write(w, public.ErrOkErr, "其他失败原因")
			return
		}
		if count != 0 {
			public.Write(w, public.ErrOkErr, "资产已存在")
			return
		}
	}
	rst := equipment.Equipment_create(input.Para)
	if rst == false {
		public.Write(w, public.ErrOkErr, "其他失败原因")
		return
	}

	public.Write(w, public.ErrOkErr, "ok")
}

func equipmentdelete(w http.ResponseWriter, queryForm url.Values) {
	input := equipmentdeleteInput{}
	Delete_ips := queryForm.Get("ip")
	if len(Delete_ips) == 0 {
		public.Write(w, public.ErrOkErr, "err")
		return
	}
	//fmt.Println(Delete_ips)
	switch strings.Contains(Delete_ips, ",") {
	case true:
		delete_ip := strings.Split(Delete_ips, ",")
		for i, _ := range delete_ip {
			if delete_ip[i] != "" {
				input.Para.Ip[i] = delete_ip[i]
			}
		}
		rst := equipment.Equipment_delete(input.Para)
		if rst == false {
			public.Write(w, public.ErrOkErr, "err")
			return
		}
	case false:
		input.Para.Ip[0] = Delete_ips
		rst := equipment.Equipment_delete(input.Para)
		if rst == false {
			public.Write(w, public.ErrOkErr, nil)
			return
		}
	default:
		public.Write(w, public.ErrOkErr, "err")
		return
	}
	//fmt.Println(input)
	public.Write(w, public.ErrOkErr, "ok")
}

func equipmentupdate(w http.ResponseWriter, queryForm url.Values) {
	input := equipmentupdateInput{}
	input.Para.Ip = queryForm.Get("ip")
	input.Para.Os_type = queryForm.Get("os_type")
	input.Para.Alias = queryForm.Get("alias")
	alias, err := base64.StdEncoding.DecodeString(input.Para.Alias)
	if nil != err {
		public.Write(w, public.ErrOkErr, "修改失败")
		return
	}
	input.Para.Alias = string(alias)
	departmentId, err := strconv.Atoi(queryForm.Get("departmentId"))
	if err != nil {
		departmentId = 0
	}
	input.Para.DepartmentId = departmentId
	//fmt.Println(input)
	rst := equipment.Equipment_update(input.Para)
	if rst == false {
		public.Write(w, public.ErrOkErr, "修改失败")
		return
	}
	public.Write(w, public.ErrOkErr, "ok")
}

type equipmentupdateInput struct {
	Para equipment.Equipmentinfo_st
}
type equipmentdeleteInput struct {
	Para equipment.EquipmentIP_st
}
type equipmentselectInput struct {
	Para equipment.Limt_st
}
type equipmentcreateInput struct {
	Para equipment.Equipmentinfo_st
}
