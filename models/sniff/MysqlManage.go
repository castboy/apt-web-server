// SniffMysqlManage.go
package sniff

import (
	"apt-web-server_v2/models/db"
	"apt-web-server_v2/models/modelsPublic"
	//"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func pcap_attack_ip(srcip string, str string) string {
	var ipstring = ""
	if srcip != "" {
		if strings.Contains(srcip, "/") == true {
			if str != "" {
				ipstring = fmt.Sprintf(`%s and src net %s`, str, srcip)
			} else {
				ipstring = fmt.Sprintf(`src net %s`, srcip)
			}
		} else {
			if str != "" {
				ipstring = fmt.Sprintf(`%s and src host %s`, str, srcip)
			} else {
				ipstring = fmt.Sprintf(`src host %s`, srcip)
			}
		}
	} else {
		if str != "" {
			ipstring = fmt.Sprintf(`%s %s`, str, srcip)
		}
	}
	return ipstring
}
func pcap_victim_ip(dstip string, str string) string {
	var ipstring = ""
	if dstip != "" {
		if strings.Contains(dstip, "/") == true {
			if str != "" {
				ipstring = fmt.Sprintf(`%s and dst net %s`, str, dstip)
			} else {
				ipstring = fmt.Sprintf(`dst net %s`, dstip)
			}
		} else {
			if str != "" {
				ipstring = fmt.Sprintf(`%s and dst host %s`, str, dstip)
			} else {
				ipstring = fmt.Sprintf(`dst host %s`, dstip)
			}
		}
	} else {
		if str != "" {
			ipstring = fmt.Sprintf(`%s %s`, str, dstip)
		}
	}
	return ipstring
}
func pacp_port(port string, str string) string {
	var portstring = ""
	if port != "" {
		if strings.Contains(port, "-") == true {
			if str != "" {
				portstring = fmt.Sprintf(`%s and portrange %s`, str, port)
			} else {
				portstring = fmt.Sprintf(`portrange %s`, port)
			}
		} else {
			if str != "" {
				portstring = fmt.Sprintf(`%s and port %s`, str, port)
			} else {
				portstring = fmt.Sprintf(`port %s`, port)
			}
		}
	} else {
		if str != "" {
			portstring = fmt.Sprintf(`%s %s`, str, port)
		}
	}
	return portstring
}
func pcap_proto(proto string, str string) string {
	var protostring = ""
	proto1 := proto
	fmt.Println(proto)
	switch proto {
	case "tcp":
		proto1 = "TCP"
	case "udp":
		proto1 = "UDP"
	default:
	}
	if proto1 != "" {
		if str != "" {
			protostring = fmt.Sprintf(`%s and proto %s`, str, proto1)
		} else {
			protostring = fmt.Sprintf(`proto %s`, proto1)
		}
	} else {
		if str != "" {
			protostring = fmt.Sprintf(`%s %s`, str, proto1)
		}
	}
	return protostring
}

/*
func Insert_to_mysql(insert_cmd string) bool {
	/*
		db := Init_mysql(Sqlcfg)
		defer db.Close()
	*
	//result, _ := db.Exec("insert into user values(?,?,?)", "test", 2, "test")
	//c, _ := result.RowsAffected()
	//fmt.Println("add affected rows:", c)
	_, err := db.DB.Exec(insert_cmd)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}
*/
func Insert_mysql(insert_cmd string, plcys Sniffplcy_st) bool {
	/*
		db := Init_mysql(Sqlcfg)
		defer db.Close()
	*/
	//result, _ := db.Exec("insert into user values(?,?,?)", "test", 2, "test")
	//c, _ := result.RowsAffected()
	//fmt.Println("add affected rows:", c)
	_, err := db.DB.Exec(insert_cmd)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else {
		//sniff 策略下发专用
		Mysql := fmt.Sprintf(`insert into sniff_plcy_dpi
		(plcy_name,src_ip,dst_ip,src_port,dst_port,proto,affect_time_start,affect_time_end)
		values('%s','%s','%s','%s','%s','%s','%d','%d')`, plcys.Plcy_name, plcys.Attack_ip, plcys.Victim_ip,
			"", plcys.Dst_port, plcys.Proto, plcys.Affect_time_start, plcys.Affect_time_end)
		_, err = db.DB.Exec(Mysql)
		if err != nil {
			fmt.Println(err.Error())
			Mysql = fmt.Sprintf("delete from sniff_plcy_ui where plcy_name = '%s'", plcys.Plcy_name)
			db.DB.Exec(Mysql)
			return false
		}
		Mysql = fmt.Sprintf(`insert into sniff_plcy_dpi
		(plcy_name,src_ip,dst_ip,src_port,dst_port,proto,affect_time_start,affect_time_end)
		values('%s','%s','%s','%s','%s','%s','%d','%d')`, plcys.Plcy_name, plcys.Victim_ip, plcys.Attack_ip,
			plcys.Dst_port, "", plcys.Proto, plcys.Affect_time_start, plcys.Affect_time_end)
		_, err = db.DB.Exec(Mysql)
		if err != nil {
			fmt.Println(err.Error())
			Mysql = fmt.Sprintf("delete from sniff_plcy_ui where plcy_name = '%s'", plcys.Plcy_name)
			db.DB.Exec(Mysql)
			return false
		}
	}
	return true
}

/*
func Select_mysql(select_cmd string) *sql.Rows {
	var rows *sql.Rows
	rows, err1 := db.DB.Query(select_cmd)
	if err1 != nil {
		fmt.Println(err1.Error())
		return rows
	}
	//fmt.Println(rows)
	//fmt.Println("rows type:", reflect.TypeOf(db))
	return rows
}
func Update_mysql(update_cmd string) bool {
	_, err := db.DB.Exec(update_cmd)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}
func Delete_mysql(delete_cmd string) bool {
	_, err := db.DB.Exec(delete_cmd)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}
*/
func Sniff_insert(plcys Sniffplcy_st) bool {
	var pcapplcy1 string
	var pcapplcy2 string
	plcys.Pcap_path = fmt.Sprintf("/pcap/%s/", plcys.Plcy_name)
	pcapplcy1 = pcap_attack_ip(plcys.Attack_ip, "")
	pcapplcy1 = pcap_victim_ip(plcys.Victim_ip, pcapplcy1)
	pcapplcy1 = pacp_port(plcys.Dst_port, pcapplcy1)
	pcapplcy1 = pcap_proto(plcys.Proto, pcapplcy1)
	//fmt.Println(pcapplcy1)
	if pcapplcy1 == "" {
		return false
	} else {
		pcapplcy2 = pcap_attack_ip(plcys.Victim_ip, "")
		pcapplcy2 = pcap_victim_ip(plcys.Attack_ip, pcapplcy2)
		pcapplcy2 = pacp_port(plcys.Dst_port, pcapplcy2)
		pcapplcy2 = pcap_proto(plcys.Proto, pcapplcy2)
		//fmt.Println(pcapplcy2)
		plcys.Pcap_plcy = fmt.Sprintf(`"(%s) or (%s)"`, pcapplcy1, pcapplcy2)
		//fmt.Println(plcys.Pcap_plcy)
		//sniff_plcy_ui 插入语句
		insert_cmd := fmt.Sprintf(`insert into sniff_plcy_ui
		(plcy_name,attack_ip,victim_ip,dst_port,proto,affect_time_start,affect_time_end,
		pcap_path,plcy_date,pcap_plcy)values('%s','%s','%s','%s','%s','%d','%d',
		'%s','%d','%s')`, plcys.Plcy_name, plcys.Attack_ip, plcys.Victim_ip,
			plcys.Dst_port, plcys.Proto, plcys.Affect_time_start, plcys.Affect_time_end,
			plcys.Pcap_path, plcys.Plcy_date, plcys.Pcap_plcy)
		fmt.Println(insert_cmd)
		if strings.Contains(plcys.Victim_ip, "/") == true {
			minIp, maxIp := getCidrIpRange(plcys.Victim_ip)
			plcys.Victim_ip = fmt.Sprintf(`%s-%s`, minIp, maxIp)
		}
		//fmt.Println(plcys.Victim_ip)
		rst := Insert_mysql(insert_cmd, plcys)
		if rst == false {
			return false
		}
	}
	return true
}
func Sniff_select(select_type Sniffselect_st) Sniffshow_st {
	var select_cmd string
	var select_s Sniffshow_st
	var tmp_plcy Sniffplcy_st
	switch select_type.Type {
	case "all":
		select_cmd = fmt.Sprintf(`select plcy_name,attack_ip,victim_ip,dst_port,
		proto,affect_time_start,affect_time_end,plcy_status,issued_status,pcap_path,plcy_date
		 from sniff_plcy_ui`)
	case "limit":
		select_cmd = fmt.Sprintf(`select plcy_name,attack_ip,victim_ip,dst_port,
		proto,affect_time_start,affect_time_end,plcy_status,issued_status,pcap_path,plcy_date
		from sniff_plcy_ui order by %s %s limit %d,%d`, select_type.Lies, select_type.Orderby,
			select_type.Page, select_type.Count)
	default:
		return select_s
	}
	//fmt.Println(select_cmd)
	rows := modelsPublic.Select_mysql(`select COUNT(plcy_name) from sniff_plcy_ui`)
	if rows == nil {
		return select_s
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&select_s.Count); err != nil {
			return select_s
		}
	}

	rows = modelsPublic.Select_mysql(select_cmd)
	if rows == nil {
		return select_s
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&(tmp_plcy.Plcy_name), &(tmp_plcy.Attack_ip),
			&(tmp_plcy.Victim_ip), &(tmp_plcy.Dst_port),
			&(tmp_plcy.Proto), &(tmp_plcy.Affect_time_start),
			&(tmp_plcy.Affect_time_end), &(tmp_plcy.Plcy_status), &(tmp_plcy.Issued_status),
			&(tmp_plcy.Pcap_path), &(tmp_plcy.Plcy_date)); err == nil {
			select_s.Plcy_s = append(select_s.Plcy_s, tmp_plcy)
			//select_s.Count = select_s.Count + 1
		}
	}
	return select_s
}
func Sniff_issued(issued Sniffissued_st) bool {
	var update_cmd string
	var count int
	var rst bool
	rst = modelsPublic.Update_mysql("update sniff_plcy_ui set plcy_status=0,issued_status=0")
	rst = modelsPublic.Update_mysql(`update sniff_plcy_issued set status=2,err=""`)
	if rst == false {
		return false
	}
	/*
		rst = Update_mysql("update sniff_plcy_ui set plcy_status=0,issued_status=0")
		if rst == false {
			return false
		}
	*/
	for i, _ := range issued.Plcy_name {
		if issued.Plcy_name[i] != "" {
			select_cmd := fmt.Sprintf(`select COUNT(plcy_name) from sniff_plcy_ui where plcy_name = '%s'`, issued.Plcy_name[i])
			rows := modelsPublic.Select_mysql(select_cmd)
			if rows != nil {
				defer rows.Close()
				for rows.Next() {
					if err := rows.Scan(&count); err != nil {
						return false
					}
					if count == 0 {
						modelsPublic.Update_mysql(`update sniff_plcy_issued set status=0,err="要下发策略不存在"`)
						return false
					}
				}
			}
			update_cmd = fmt.Sprintf(`update sniff_plcy_ui set plcy_status =
			1 where plcy_name = '%s'`, issued.Plcy_name[i])
			rst = modelsPublic.Update_mysql(update_cmd)
			if rst == false {
				modelsPublic.Update_mysql(`update sniff_plcy_issued set status=0,err="生成策略条件失败"`)
				return false
			}
		}
	}
	return true
}
func Sniff_delete(delete_st Sniffissued_st) bool {
	var delete_cmd string
	var rst bool
	for i, _ := range delete_st.Plcy_name {
		if delete_st.Plcy_name[i] != "" {
			delete_cmd = fmt.Sprintf(`delete from sniff_plcy_ui where plcy_name = '%s'`, delete_st.Plcy_name[i])
			rst = modelsPublic.Delete_mysql(delete_cmd)
			if rst == false {
				return false
			}
		}
	}
	return true
}

func Sniff_status(status_s Sniffstatus_st) Sniffstatus_st {
	var select_cmd string
	switch status_s.Type {
	case "issued":
		select_cmd = fmt.Sprintf(`select status,err from sniff_plcy_issued`)
	default:
		return status_s
	}

	rows := modelsPublic.Select_mysql(select_cmd)
	if rows == nil {
		return status_s
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&(status_s.Msg.Status), &(status_s.Msg.Log)); err == nil {
			//fmt.Println(status_s.Msg.Status)
		}
	}
	return status_s
}
