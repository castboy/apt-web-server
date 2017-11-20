// EquipmentMysql.go
package equipment

import (
	"apt-web-server_v2/models/db"
	"apt-web-server_v2/models/modelsPublic"
	//"database/sql"
	"apt-web-server_v2/modules/mlog"
	"fmt"
	//"reflect"
	//"strings"
	//"bytes"
	"os/exec"

	_ "github.com/go-sql-driver/mysql"
)

func Equipment_scan() bool {
	_, err := db.DB.Exec("update equipment_scanflag set flag=1")
	if err != nil {
		return false
	}
	_, err = db.DB.Exec("delete from equipment where authority = 2")
	if err != nil {
		return false
	}
	_, err = db.DB.Exec("delete from equipment_scan")
	if err != nil {
		return false
	}
	cmd := exec.Command("/bin/sh", "/root/equipment_scan.sh")
	err = cmd.Start()
	if err != nil {
		fmt.Println("run err:", err.Error())
	}
	fmt.Println(cmd)
	return true
}
func Import_mysql(importdates Importdate_st) bool {
	var insert_cmd string
	tx, err := db.DB.Begin()
	if err != nil {
		return false
	}
	defer tx.Rollback()
	tx.Exec("delete from equipment_import")
	for i, _ := range importdates.Equipment {
		insert_cmd = fmt.Sprintf(`insert into equipment_import (ip,os_type,alias)
		values('%s','%s','%s')`, importdates.Equipment[i].Ip,
			importdates.Equipment[i].Os_type, importdates.Equipment[i].Alias)
		//fmt.Println(insert_cmd)
		_, err := tx.Exec(insert_cmd)
		if err != nil {
			fmt.Println(err.Error())
			tx.Rollback()
			return false
		}
	}
	tx.Commit()
	return true
}
func Equipment_create(input Equipmentinfo_st) bool {
	insert_cmd := fmt.Sprintf(`replace into equipment (ip,os_type,alias,authority,
	data_source,departmentId,time)values('%s','%s','%s',1,2,'%d',current_timestamp())`, input.Ip,
		input.Os_type, input.Alias, input.DepartmentId)
	//fmt.Println(insert_cmd)
	rst := modelsPublic.Insert_to_mysql(insert_cmd)
	if rst == false {
		mlog.Debug("##equipment_create:", insert_cmd)
		return false
	}
	return true
}
func Equipment_delete(ip EquipmentIP_st) bool {
	var delete_cmd string
	var rst bool
	//var count int
	for i, _ := range ip.Ip {
		if ip.Ip[i] != "" {
			delete_cmd = fmt.Sprintf(`update equipment set delete_flag = 1,
			authority =1 where ip = '%s'`, ip.Ip[i])
			rst = modelsPublic.Delete_mysql(delete_cmd)
			if rst == false {
				mlog.Debug("##equipment_delete:", delete_cmd)
				//fmt.Println(delete_cmd)
				return false
			}
		}
	}
	return true
}
func Equipment_update(input Equipmentinfo_st) bool {
	update_cmd := fmt.Sprintf(`update equipment set os_type='%s',alias='%s',
	authority=1,data_source=2,departmentId=%d where ip = '%s'`, input.Os_type,
		input.Alias, input.DepartmentId, input.Ip)
	rst := modelsPublic.Update_mysql(update_cmd)
	if rst == false {
		mlog.Debug("##equipment_update:", update_cmd)
		return false
	}
	return true
}
func Equipment_select(input Limt_st) Equipmentshow_st {
	var mysql_cmd string
	var count_cmd string
	var select_s Equipmentshow_st
	var tmp Equipmentshowinfo_st
	mysql_cmd = fmt.Sprintf(`select ip,os_type,alias,attack_count,
		IFNULL(department_name,"") from equipment left join department on 
		departmentId = id where %s = 0`, "delete_flag")
	count_cmd = fmt.Sprintf(`select count(ip) from equipment left join department on 
		departmentId = id where %s = 0`, "delete_flag")
	if input.Ip != "" {
		mysql_cmd = fmt.Sprintf(`%s and ip='%s'`, mysql_cmd, input.Ip)
		count_cmd = fmt.Sprintf(`%s and ip='%s'`, count_cmd, input.Ip)
	}
	if input.Os_type != "" {
		mysql_cmd = fmt.Sprintf(`%s and os_type='%s'`, mysql_cmd, input.Os_type)
		count_cmd = fmt.Sprintf(`%s and os_type='%s'`, count_cmd, input.Os_type)
	}
	switch input.DepartmentId {
	case -1:
		//查询全部
	default:
		//查询特定部门 （包含未定义、指定部门id的）
		mysql_cmd = fmt.Sprintf(`%s and departmentId=%d`, mysql_cmd, input.DepartmentId)
		count_cmd = fmt.Sprintf(`%s and departmentId=%d`, count_cmd, input.DepartmentId)
	}
	if input.Lies != "" && input.Orderby != "" {
		mysql_cmd = fmt.Sprintf(`%s order by %s %s`, mysql_cmd, input.Lies, input.Orderby)
	}
	mysql_cmd = fmt.Sprintf(`%s limit %d,%d`, mysql_cmd, input.Page, input.Count)
	//fmt.Println(mysql_cmd)
	//fmt.Println(count_cmd)
	rows := modelsPublic.Select_mysql(count_cmd)
	if rows == nil {
		return select_s
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&(select_s.Count)); err != nil {
			return select_s
		}
	}
	if select_s.Count == 0 {
		return select_s
	}
	rows2 := modelsPublic.Select_mysql(mysql_cmd)
	if rows2 == nil {
		return select_s
	}
	defer rows2.Close()
	for rows2.Next() {
		if err := rows2.Scan(&(tmp.Ip), &(tmp.Os_type),
			&(tmp.Alias), &(tmp.Attack_count),
			&(tmp.Departmentname)); err == nil {
			tmp.Service_port = ""
			mysql_cmd = fmt.Sprintf(`select GROUP_CONCAT(service_port)
				from equipment_service where service_ip='%s'`, tmp.Ip)
			//fmt.Println(mysql_cmd)
			rows3 := modelsPublic.Select_mysql(mysql_cmd)
			if rows3 == nil {
				return select_s
			}
			defer rows3.Close()
			for rows3.Next() {
				if err := rows3.Scan(&(tmp.Service_port)); err == nil {
				}
			}
			select_s.Info = append(select_s.Info, tmp)
		}
	}
	return select_s
}
func Equipment_scanflag() bool {
	var flag int
	rows := modelsPublic.Select_mysql("select flag from equipment_scanflag")
	if rows == nil {
		return false
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&(flag)); err == nil {
			//fmt.Println(err.Error())
			//fmt.Println("tmp", tmp)
			//select_s.Info = append(select_s.Info, tmp)
		}
	}
	if flag == 0 {
		return false
	} else {
		return true
	}

}
func Detailinfo_select(input Limt_st) Showdetailinfo_st {
	var select_s Showdetailinfo_st
	var tmp Detailinfo_st
	var mysql_cmd string
	var count_cmd string
	count_cmd = fmt.Sprintf(`select count(service_ip) from equipment_service 
	where service_ip='%s'`, input.Ip)
	//fmt.Println(count_cmd)
	mysql_cmd = fmt.Sprintf(`select service_ip,service_name,IFNULL(service_type,""),
	IFNULL(service_version,""),IFNULL(service_platform,""),IFNULL(service_port,""),
	IFNULL(service_banner,"") from equipment_service 
	where service_ip='%s' limit %d,%d`, input.Ip, input.Page, input.Count)
	//fmt.Println(mysql_cmd)

	rows := modelsPublic.Select_mysql(count_cmd)
	if rows == nil {
		return select_s
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&(select_s.Counts)); err != nil {
			return select_s
		}
	}
	if select_s.Counts == 0 {
		return select_s
	}

	rows1 := modelsPublic.Select_mysql(mysql_cmd)
	if rows1 == nil {
		return select_s
	}
	defer rows1.Close()
	for rows1.Next() {
		if err := rows1.Scan(&(tmp.Service_ip), &(tmp.Service_name), &(tmp.Service_type),
			&(tmp.Service_version), &(tmp.Service_platform), &(tmp.Service_port),
			&(tmp.Service_banner)); err == nil {
			select_s.Detailinfo = append(select_s.Detailinfo, tmp)
		}
	}
	return select_s
}

func Createdepartment(input Departmentip_st) string {
	var mysql_cmd string
	var count int
	//检查资产是否存在
	mysql_cmd = fmt.Sprintf(`select count(id) from department where 
	department_name='%s'`, input.Departmentname)
	rows := modelsPublic.Select_mysql(mysql_cmd)
	if rows == nil {
		return "其他失败原因"
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return "其他失败原因"
		}
		if count != 0 {
			return "部门已存在"
		}
	}
	//添加部门
	mysql_cmd = fmt.Sprintf(`insert into department (department_name) values('%s')`,
		input.Departmentname)
	rst := modelsPublic.Insert_to_mysql(mysql_cmd)
	if rst == false {
		return "部门添加失败"
	}
	//查看已添加部门 id
	mysql_cmd = fmt.Sprintf(`select id from department where department_name='%s'`,
		input.Departmentname)
	//fmt.Println(mysql_cmd)
	rows1 := modelsPublic.Select_mysql(mysql_cmd)
	if rows1 == nil {
		return "其他失败原因"
	}
	defer rows1.Close()
	for rows1.Next() {
		if err := rows1.Scan(&(input.DepartmentId)); err != nil {
			return "其他失败原因"
		}
		//fmt.Println(input.DepartmentId)
		if input.DepartmentId == 0 {
			return "其他失败原因"
		}
	}

	//给部门添加资产
	tx, err := db.DB.Begin()
	if err != nil {
		return "部门添加成功,资产添加失败"
	}
	defer tx.Rollback()
	for i, _ := range input.Ip {
		mysql_cmd = fmt.Sprintf(`update equipment set departmentId=%d where 
		ip='%s'`, input.DepartmentId, input.Ip[i])
		_, err := tx.Exec(mysql_cmd)
		if err != nil {
			fmt.Println(err.Error())
			tx.Rollback()
			return "部门添加成功,资产添加失败"
		}
	}
	tx.Commit()
	return "ok"
}
func Updatedepartment(input Departmentip_st) string {
	var mysql_cmd string
	var count int
	var id int
	//检查部门是否存在
	mysql_cmd = fmt.Sprintf(`select id from department where 
	department_name='%s'`, input.Departmentname)
	//fmt.Println(mysql_cmd)
	rows := modelsPublic.Select_mysql(mysql_cmd)
	if rows == nil {
		return "其他失败原因"
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return "其他失败原因"
		}
		count = count + 1
	}
	switch count {
	case 0:
		//修改部门名称
		mysql_cmd = fmt.Sprintf(`update department set department_name='%s' 
			where id=%d`, input.Departmentname, input.DepartmentId)
		rst := modelsPublic.Update_mysql(mysql_cmd)
		if rst == false {
			return "其他原因失败"
		}
	case 1:
		//不做部门名称修改
		if id != input.DepartmentId {
			return "部门名称重复"
		}
	default:
		//已存在部门名称
		return "部门名称重复"
	}

	//修改部门资产
	tx, err := db.DB.Begin()
	if err != nil {
		return "修改部门信息成功，修改资产信息失败"
	}
	defer tx.Rollback()
	mysql_cmd = fmt.Sprintf(`update equipment set departmentId=0 where 
		departmentId=%d`, input.DepartmentId)
	_, err = tx.Exec(mysql_cmd)
	if err != nil {
		tx.Rollback()
		return "修改部门信息成功，修改资产信息失败"
	}
	for i, _ := range input.Ip {
		mysql_cmd = fmt.Sprintf(`update equipment set departmentId=%d where 
		ip='%s'`, input.DepartmentId, input.Ip[i])
		_, err := tx.Exec(mysql_cmd)
		if err != nil {
			fmt.Println(err.Error())
			tx.Rollback()
			return "修改部门信息成功，修改资产信息失败"
		}
	}
	tx.Commit()
	return "ok"
}
func Ip_select(input int) Departmentip_st {
	var departmentip Departmentip_st
	var ip string
	var mysql_cmd string
	mysql_cmd = fmt.Sprintf(`select ip from equipment where departmentId=%d`, input)
	//fmt.Println(mysql_cmd)
	rows := modelsPublic.Select_mysql(mysql_cmd)
	if rows == nil {
		return departmentip
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&ip); err == nil {
			departmentip.Ip = append(departmentip.Ip, ip)
		}
	}
	return departmentip
}

func Department_create(input string) string {
	var mysql_cmd string
	var count int
	//检查资产是否存在
	mysql_cmd = fmt.Sprintf(`select count(id) from department where 
	department_name='%s'`, input)
	rows := modelsPublic.Select_mysql(mysql_cmd)
	if rows == nil {
		return "其他失败原因"
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return "其他失败原因"
		}
		if count != 0 {
			return "部门已存在"
		}
	}
	//添加部门
	mysql_cmd = fmt.Sprintf(`insert into department (department_name) values('%s')`,
		input)
	rst := modelsPublic.Insert_to_mysql(mysql_cmd)
	if rst == false {
		return "部门添加失败"
	}
	return "ok"
}
func Department_delete(input int) string {
	var mysql_cmd string
	tx, err := db.DB.Begin()
	if err != nil {
		return "err"
	}
	defer tx.Rollback()
	mysql_cmd = fmt.Sprintf(`delete from department where id=%d`, input)
	//fmt.Println(mysql_cmd)
	_, err = tx.Exec(mysql_cmd)
	if err != nil {
		tx.Rollback()
		return "err"
	}
	mysql_cmd = fmt.Sprintf(`update equipment set departmentId=0 where 
		departmentId=%d`, input)
	//fmt.Println(mysql_cmd)
	_, err = tx.Exec(mysql_cmd)
	if err != nil {
		tx.Rollback()
		return "err"
	}
	tx.Commit()
	return "ok"
}
func Managedepartment(input Departmentip_st) string {
	var mysql_cmd string
	//给部门添加资产
	tx, err := db.DB.Begin()
	if err != nil {
		return "资产添加失败"
	}
	defer tx.Rollback()
	for i, _ := range input.Ip {
		mysql_cmd = fmt.Sprintf(`update equipment set departmentId=%d where 
		ip='%s'`, input.DepartmentId, input.Ip[i])
		_, err := tx.Exec(mysql_cmd)
		if err != nil {
			fmt.Println(err.Error())
			tx.Rollback()
			return "资产添加失败"
		}
	}
	if input.DepartmentId != 0 {
		count := 0
		mysql_cmd = fmt.Sprintf(`select count(id) from department where 
			id='%d'`, input.DepartmentId)
		fmt.Println(mysql_cmd)
		rows1 := modelsPublic.Select_mysql(mysql_cmd)
		if rows1 == nil {
			tx.Rollback()
			return "其他失败原因"
		}
		defer rows1.Close()
		for rows1.Next() {
			if err := rows1.Scan(&count); err != nil {
				tx.Rollback()
				return "其他失败原因"
			}
			if count == 0 {
				tx.Rollback()
				return "部门已经不存在"
			}
		}
	}
	tx.Commit()
	return "ok"
}
func Department_update(input Departmentip_st) string {
	var mysql_cmd string
	var count int
	var id int
	//检查部门是否存在
	mysql_cmd = fmt.Sprintf(`select id from department where 
	department_name='%s'`, input.Departmentname)
	//fmt.Println(mysql_cmd)
	rows := modelsPublic.Select_mysql(mysql_cmd)
	if rows == nil {
		return "其他失败原因"
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return "其他失败原因"
		}
		count = count + 1
	}
	switch count {
	case 0:
		//修改部门名称
		mysql_cmd = fmt.Sprintf(`update department set department_name='%s' 
			where id=%d`, input.Departmentname, input.DepartmentId)
		rst := modelsPublic.Update_mysql(mysql_cmd)
		if rst == false {
			return "其他原因失败"
		}
	case 1:
		//不做部门名称修改
		if id != input.DepartmentId {
			return "部门名称重复"
		}
	default:
		//已存在部门名称
		return "部门名称重复"
	}
	return "ok"
}
func Department_select(input Limt_st) Departmentshow_st {
	var select_s Departmentshow_st
	var tmp Departmentshowinfo_st
	var mysql_cmd string
	var deptcount_cmd string
	var ipcount_cmd string
	if input.DepartmentId == 0 && input.Count == 0 {
		mysql_cmd = fmt.Sprintf(`select id,department_name from %s`, "department")
	} else {
		mysql_cmd = fmt.Sprintf(`select id,department_name from department limit %d,%d`,
			input.Page, input.Count)
	}
	deptcount_cmd = fmt.Sprintf(`select count(id) from %s`, "department")
	rows := modelsPublic.Select_mysql(deptcount_cmd)
	if rows == nil {
		return select_s
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&(select_s.Count)); err == nil {
		}
	}

	rows1 := modelsPublic.Select_mysql(mysql_cmd)
	if rows1 == nil {
		return select_s
	}
	defer rows1.Close()
	for rows1.Next() {
		if err := rows1.Scan(&(tmp.DepartmentId),
			&(tmp.Departmentname)); err == nil {
			ipcount_cmd = fmt.Sprintf(`select count(ip) from equipment where 
			departmentId=%d`, tmp.DepartmentId)
			rows2 := modelsPublic.Select_mysql(ipcount_cmd)
			if rows2 == nil {
				return select_s
			}
			defer rows2.Close()
			for rows2.Next() {
				if err := rows2.Scan(&(tmp.Equipmentcount)); err == nil {
				}
			}
			select_s.Info = append(select_s.Info, tmp)
		}
	}
	return select_s
}
