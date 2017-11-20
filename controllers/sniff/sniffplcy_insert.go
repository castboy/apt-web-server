// sniffplcy_insert.go
//sniff 策略插入
package sniff

import (
	"apt-web-server_v2/controllers/public"
	"apt-web-server_v2/models/modelsPublic"
	"apt-web-server_v2/models/sniff"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type SniffcreatController struct{}

var SniffcreatObj = new(SniffcreatController)

func (this *SniffcreatController) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//check ssn
	//TODO
	//parse parament
	//input := models.TblUrgencySearchPara{}
	var count int
	input := SniffcreatInput{}

	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		public.Write(w, public.ErrParamentErr, nil)
		return
	}

	fmt.Println(queryForm)
	input.Para.Plcy_name = queryForm.Get("plcy_name")
	input.Para.Attack_ip = queryForm.Get("attack_ip")
	input.Para.Victim_ip = queryForm.Get("victim_ip")
	input.Para.Dst_port = queryForm.Get("dst_port")
	input.Para.Proto = queryForm.Get("proto")
	start_time, err := strconv.Atoi(queryForm.Get("start_time"))
	if err != nil {
		start_time = 0
	}
	input.Para.Affect_time_start = start_time
	end_time, err := strconv.Atoi(queryForm.Get("end_time"))
	if err != nil {
		end_time = 0
	}
	if start_time > end_time {
		public.Write(w, public.ErrOkErr, "开始时间应小于结束时间")
		return
	}
	input.Para.Affect_time_end = end_time
	input.Para.Affect_time_start = start_time
	creatdate, err := strconv.Atoi(queryForm.Get("creatdate"))
	if (err != nil) || (creatdate == 0) {
		public.Write(w, public.ErrOkErr, "其他失败原因")
		return
	}
	input.Para.Plcy_date = creatdate
	fmt.Println(input)
	insert_cmd := fmt.Sprintf(`select COUNT(plcy_name) from sniff_plcy_ui where plcy_name = '%s'`, input.Para.Plcy_name)
	rows := modelsPublic.Select_mysql(insert_cmd)
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			if err := rows.Scan(&count); err != nil {
				public.Write(w, public.ErrOkErr, "其他失败原因")
				return
			}
			if count != 0 {
				public.Write(w, public.ErrOkErr, "策略已存在")
				return
			}
		}
	}
	rst := sniff.Sniff_insert(input.Para)
	if rst == false {
		public.Write(w, public.ErrOkErr, "其他失败原因")
		return
	}
	public.Write(w, public.ErrOkErr, "ok")
}

type SniffcreatInput struct {
	Para sniff.Sniffplcy_st
}
