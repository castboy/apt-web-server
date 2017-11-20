/********SSLCert统计********/
package whiteList

import (
	"apt-web-server_v2/models/db"
	"apt-web-server_v2/models/modelsPublic"
	"fmt"
	//"strings"
	//"time"
)

func WL_wlTableName() string {
	return "abnconn_whitelist"
}

func WL_add(wLpara *TblWLOperateParaIn) (err error /*num_added int32,*/, num_repeated int32) {
	var db_cmd string
	var query_cmd string
	var count int
	tx, erR := db.DB.Begin()
	if erR != nil {
		return erR, num_repeated
	}
	//defer tx.Rollback()
	for _, ele := range wLpara.WLOpElement {
		query_cmd = fmt.Sprintf(`select count(id) from %s where 
			src_ip='%s' and src_port='%d' and dest_ip='%s' and dest_port='%d' and proto='%d'`,
			WL_wlTableName(), ele.Sip, ele.Sport, ele.Dip, ele.Dport, ele.Proto)
		//fmt.Println("WL_add query_cmd is ", query_cmd)

		rows := modelsPublic.Select_mysql(query_cmd)
		for rows.Next() {
			if rows != nil {
				if err = rows.Scan(&count); err != nil {
					//fmt.Println("err is  ", err.Error())
					tx.Rollback()
					rows.Close()
					return err, num_repeated
				}
				//fmt.Println("WL_add query_cmd count is ", count)
				if count != 0 {
					rows.Close()
					break
				}
			}
		}
		if count > 0 {
			num_repeated++ //record the repeated wl number
			continue       //white list already exists, handle next one
		}

		// insert white list into db table
		db_cmd = fmt.Sprintf(`INSERT INTO %s (src_ip,src_port,dest_ip,dest_port,proto) 
			VALUES('%s',%d,'%s',%d,%d)`,
			WL_wlTableName(), ele.Sip, ele.Sport, ele.Dip, ele.Dport, ele.Proto)
		//fmt.Println("WL_add db_cmd is ", db_cmd)

		_, err = tx.Exec(db_cmd)
		if err != nil {
			//fmt.Println(err.Error())
			tx.Rollback()
			return err, num_repeated
		}
	}
	tx.Commit()

	return err, num_repeated
}

func WL_delete(wLpara *TblWLOperateParaIn) (err error, notFound int32) {
	var db_cmd string
	var query_cmd string
	var count int
	tx, erR := db.DB.Begin()
	if erR != nil {
		return erR, notFound
	}
	//defer tx.Rollback()
	for _, ele := range wLpara.WLOpElement {
		query_cmd = fmt.Sprintf(`select count(id) from %s where
			src_ip='%s' and src_port='%d' and dest_ip='%s' and dest_port='%d' and proto='%d'`,
			WL_wlTableName(), ele.Sip, ele.Sport, ele.Dip, ele.Dport, ele.Proto)
		//query_cmd = fmt.Sprintf(`select count(id) from %s where
		//	src_ip='%s' and src_port='%s' and dest_ip='%s' and dest_port='%s' and proto='%s'`,
		//	WL_wlTableName(), ele.Sip, ele.Sport, ele.Dip, ele.Dport, ele.Proto)
		//fmt.Println("WL_del query_cmd is ", query_cmd)

		rows := modelsPublic.Select_mysql(query_cmd)
		for rows.Next() {
			if rows != nil {
				if err = rows.Scan(&count); err != nil {
					//fmt.Println("err is  ", err.Error())
					tx.Rollback()
					return err, notFound
				}
				//fmt.Println("WL_del query_cmd count is ", count)
				if count != 0 {
					break
				}
			}
		}
		if count == 0 {
			notFound++
			continue //white list doesn't exists, handle next one
		}

		// delete white list into db table
		db_cmd = fmt.Sprintf(`DELETE FROM %s  
			WHERE src_ip='%s' and src_port='%d' and dest_ip='%s' and dest_port='%d' and proto='%d'`,
			WL_wlTableName(), ele.Sip, ele.Sport, ele.Dip, ele.Dport, ele.Proto)
		//fmt.Println("WL_del db_cmd is ", db_cmd)

		_, err = tx.Exec(db_cmd)
		if err != nil {
			//fmt.Println(err.Error())
			tx.Rollback()
			return err, notFound
		}
	}
	tx.Commit()

	return err, notFound
}

func WL_clear(wLpara *TblWLOperateParaIn) (err error) {
	var db_cmd string

	db_cmd = fmt.Sprintf(`TRUNCATE TABLE %s`,
		WL_wlTableName())

	//fmt.Println("WL_clear db_cmd is ", db_cmd)

	_, err1 := db.DB.Exec(db_cmd)
	if err1 != nil {
		fmt.Println("WL_clear,", err1.Error())
		return err1
	}

	return err1
}
