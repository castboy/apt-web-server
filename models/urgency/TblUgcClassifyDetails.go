/********获取紧急事件分类数********/
package urgency

import (
	"apt-web-server_v2/models/db"
	"apt-web-server_v2/modules/mlog"
	"fmt"
	"strings"
	"time"
)

func (this *TblUgcClassify) TableName() string {
	return "urgencymold"
}

func GetUgcCCountMysqlCmd(sqltype string, para *TblUgcClassifySearchPara) []string {
	var flag int
	qslice := make([]string, 0)
	if sqltype == "urgency" {
		qslice_tmp := fmt.Sprintf(`SELECT SUM(CASE WHEN attack_type='webshell' THEN 1 ELSE 0 END) AS webshell,
		SUM(CASE WHEN attack_type='exceptionalvisit' THEN 1 ELSE 0 END) AS exceptionalvisit,
		SUM(CASE WHEN attack_type='abnormal_connection' THEN 1 ELSE 0 END) AS abnormal_connection FROM urgencymold`)
		qslice = append(qslice, qslice_tmp)
	}
	if sqltype == "alert_waf" {
		qslice_tmp := fmt.Sprintf(`SELECT SUM(CASE WHEN attack='sqli' THEN 1 ELSE 0 END) AS sqli,
		SUM(CASE WHEN attack='xss' THEN 1 ELSE 0 END) AS xss,
		SUM(CASE WHEN attack='injection_php' THEN 1 ELSE 0 END) AS injection_php,
		SUM(CASE WHEN attack='rfi' THEN 1 ELSE 0 END) AS rfi FROM alert_waf WHERE severity IN (0,1,2)`)
		qslice = append(qslice, qslice_tmp)
		flag = 1
	}
	if para.PField.Start != 0 && para.PField.End != 0 {
		if flag != 0 {
			qslice_tmp := fmt.Sprintf(` AND time BETWEEN %d AND %d`, para.PField.Start, para.PField.End)
			qslice = append(qslice, qslice_tmp)
		} else {
			qslice_tmp := fmt.Sprintf(` WHERE time BETWEEN %d AND %d`, para.PField.Start, para.PField.End)
			qslice = append(qslice, qslice_tmp)
			flag = 1
		}
	}
	qslice = append(qslice, ";")
	//fmt.Println(tablename, para)
	fmt.Println(qslice)
	return qslice
}

func GetUgcClassifyMysqlCmd(tablename string, para *TblUgcClassifySearchPara) []string {
	var flag int
	qslice := make([]string, 0)
	if para.Type == "webshell" ||
		para.Type == "abnormal_connection" ||
		para.Type == "exceptionalvisit" ||
		para.Type == "" {
		if para.Unit != "" {
			qslice_tmp := fmt.Sprintf(`SELECT attack_type,severity,count(attack_type),time FROM urgencymold`)
			qslice = append(qslice, qslice_tmp)
		} else {
			qslice_tmp := fmt.Sprintf(`SELECT time,src_ip,dest_ip,attack_type,severity,details FROM urgencymold`)
			qslice = append(qslice, qslice_tmp)
		}
		if para.Type != "" {
			qslice_tmp := fmt.Sprintf(` WHERE attack_type IN ('%s')`, para.Type)
			qslice = append(qslice, qslice_tmp)
			flag = 1
		}
		if para.PField.Start != 0 && para.PField.End != 0 {
			if flag == 1 {
				qslice_tmp := fmt.Sprintf(` AND time BETWEEN %d AND %d`, para.PField.Start, para.PField.End)
				qslice = append(qslice, qslice_tmp)
			} else {
				qslice_tmp := fmt.Sprintf(` WHERE time BETWEEN %d AND %d`, para.PField.Start, para.PField.End)
				qslice = append(qslice, qslice_tmp)
				flag = 1
			}
		}
	}

	if para.Type == "" {
		qslice_tmp := fmt.Sprintf(` UNION ALL`)
		qslice = append(qslice, qslice_tmp)
	}
	if para.Type == "xss" ||
		para.Type == "injection_php" ||
		para.Type == "rfi" ||
		para.Type == "sqli" ||
		para.Type == "" {
		if para.Unit != "" {
			qslice_tmp := fmt.Sprintf(` SELECT attack,severity,COUNT(attack),time FROM alert_waf`)
			qslice = append(qslice, qslice_tmp)
		} else {
			qslice_tmp := fmt.Sprintf(` SELECT time,client,hostname,attack,severity,rule_data FROM alert_waf`)
			qslice = append(qslice, qslice_tmp)
		}
		if para.Type == "" {
			qslice_tmp := fmt.Sprintf(` WHERE attack IN ('sqli','xss','injection_php','rfi') AND severity IN (0,1,2)`)
			qslice = append(qslice, qslice_tmp)
		} else {
			qslice_tmp := fmt.Sprintf(` WHERE attack IN ('%s') AND severity IN (0,1,2)`, para.Type)
			qslice = append(qslice, qslice_tmp)
			flag = 1
		}
		if para.PField.Start != 0 && para.PField.End != 0 {
			if flag == 1 {
				qslice_tmp := fmt.Sprintf(` AND time BETWEEN %d AND %d`, para.PField.Start, para.PField.End)
				qslice = append(qslice, qslice_tmp)
			} else {
				qslice_tmp := fmt.Sprintf(` WHERE time BETWEEN %d AND %d`, para.PField.Start, para.PField.End)
				qslice = append(qslice, qslice_tmp)
				flag = 1
			}
		}
	}

	if para.Sort != "" {
		temp_s := fmt.Sprintf(" ORDER BY %s %s,severity DESC", para.Sort, para.Order)
		qslice = append(qslice, temp_s)
	} else {
		temp_s := fmt.Sprintf(" ORDER BY severity")
		qslice = append(qslice, temp_s)
	}
	if para.LastCount != 0 {
		temp_LC := fmt.Sprintf(" LIMIT %d", para.LastCount)
		qslice = append(qslice, temp_LC)
	} else if para.Count != 0 {
		para.Page = para.Page - 1
		temp_PC := fmt.Sprintf(" LIMIT %d,%d", para.Page*para.Count, para.Count)
		qslice = append(qslice, temp_PC)
	}
	qslice = append(qslice, ";")
	fmt.Println(tablename, para)
	return qslice
}

func (this *TblUgcClassify) GetUgcClassify(para *TblUgcClassifySearchPara) (error, *TblUgcClassifyData) {
	var tableName string
	qslice := GetUgcClassifyMysqlCmd(this.TableName(), para)
	query := strings.Join(qslice, "")
	fmt.Println(qslice)
	rows, err := db.DB.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	list := TblUgcClassifyData{}
	for rows.Next() {
		ugc := new(TblUgcClassify)
		err = rows.Scan(
			&ugc.Time,
			&ugc.Srcip,
			&ugc.Destip,
			&ugc.AttackType,
			&ugc.Severity,
			&ugc.Details)
		if err != nil {
			mlog.Debug(query, err)
		}
		list.Elements = append(list.Elements, TblUgcClassifyList{ugc.TblUgcClassifyContent})
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	if para.Type == "webshell" ||
		para.Type == "abnormal_connection" ||
		para.Type == "exceptionalvisit" {
		tableName = "urgencymold"
	} else if para.Type == "xss" ||
		para.Type == "injection_php" ||
		para.Type == "rfi" ||
		para.Type == "sqli" {
		tableName = "alert_waf"
	}
	list.Counts = GetUgcClassifyCounts("", para)
	if para.Type != "" {
		list.Totality = GetUgcClassifyCounts(tableName, para)
	} else {
		list.Totality = GetUgcClassifyCounts("urgencymold", para) + GetUgcClassifyCounts("alert_waf", para)
	}
	return nil, &list
}

func (this *TblUgcClassify) GetUgcCCount(para *TblUgcClassifySearchPara) (error, *TblUgcCCountData) {
	var i, ucount int64
	start := para.PField.Start
	end := para.PField.End
	seconds := GetUgcSeconds(para.Unit)
	list := TblUgcCCountData{}
	for i = 0; i <= ((end - start) / seconds); i++ {
		para.PField.Start = start + seconds*i
		para.PField.End = start + seconds*(i+1) - 1
		if para.PField.End > end {
			para.PField.End = end
		}
		ugc := new(TblUgcCCount)
		qslice_u := GetUgcCCountMysqlCmd("urgency", para)
		query_u := strings.Join(qslice_u, "")
		fmt.Println(qslice_u)
		rows_u, err := db.DB.Query(query_u)
		if err != nil {
			return err, nil
		}
		defer rows_u.Close()
		for rows_u.Next() {
			err = rows_u.Scan(
				&ugc.Webshell,
				&ugc.ExceptionalVisit,
				&ugc.AbnormalConnection)
			if err != nil {
				ugc.Webshell = 0
				ugc.ExceptionalVisit = 0
				ugc.AbnormalConnection = 0
				//return err, nil
			}
		}
		if err := rows_u.Err(); err != nil {
			return err, nil
		}
		qslice_a := GetUgcCCountMysqlCmd("alert_waf", para)
		query_a := strings.Join(qslice_a, "")
		fmt.Println(qslice_u)
		rows_a, err := db.DB.Query(query_a)
		if err != nil {
			return err, nil
		}
		defer rows_a.Close()
		for rows_a.Next() {
			err = rows_a.Scan(
				&ugc.Sqli,
				&ugc.Xss,
				&ugc.InjectionPHP,
				&ugc.Rfi)
			if err != nil {
				ugc.Sqli = 0
				ugc.Xss = 0
				ugc.InjectionPHP = 0
				//return err, nil
			}
		}
		if err := rows_a.Err(); err != nil {
			return err, nil
		}
		if start == 0 && end == 0 {
			ugc.Time = int64(time.Now().Unix())
		} else {
			ugc.Time = para.PField.Start
		}
		ucount += (ugc.Webshell + ugc.ExceptionalVisit + ugc.AbnormalConnection + ugc.Sqli + ugc.Xss + ugc.InjectionPHP + ugc.Rfi)
		list.Elements = append(list.Elements, TblUgcCCountList{ugc.TblUgcCCountContent})
	}
	list.Counts = ucount
	return nil, &list
}

func GetUgcClassifyCounts(tablename string, para *TblUgcClassifySearchPara) int64 {
	//qslice, whereflag := DefaultParaCmd("getcounts", tablename, &para.PField)
	var whereflag int
	qslice := make([]string, 0)
	if tablename != "" {
		qslice_tmp := fmt.Sprintf(`SELECT COUNT(id) FROM %s`, tablename)
		qslice = append(qslice, qslice_tmp)
		if tablename == "alert_waf" {
			if para.Type != "" {
				temp_t := fmt.Sprintf(" WHERE attack IN ('%s')", para.Type)
				qslice = append(qslice, temp_t)
			} else {
				qslice_waf := fmt.Sprintf(` WHERE attack IN ('sqli','xss','injection_php','rfi') AND severity IN (0,1,2)`)
				qslice = append(qslice, qslice_waf)
			}
			whereflag = 1
		}
		if tablename == "urgencymold" {
			if para.Type != "" {
				if whereflag != 0 {
					temp_t := fmt.Sprintf(" AND attack_type IN ('%s')", para.Type)
					qslice = append(qslice, temp_t)
				} else {
					temp_t := fmt.Sprintf(" WHERE attack_type IN ('%s')", para.Type)
					qslice = append(qslice, temp_t)
					whereflag = 1
				}
			}
		}

		if para.PField.Start != 0 && para.PField.End != 0 {
			if whereflag == 1 {
				qslice_tmp := fmt.Sprintf(` AND time BETWEEN %d AND %d`,
					para.PField.Start, para.PField.End)
				qslice = append(qslice, qslice_tmp)
			} else {
				qslice_tmp := fmt.Sprintf(` WHERE time BETWEEN %d AND %d`,
					para.PField.Start, para.PField.End)
				qslice = append(qslice, qslice_tmp)
				whereflag = 1
			}
		}
	} else {
		qslice_count := fmt.Sprintf(`SELECT FOUND_ROWS() as count`)
		qslice = append(qslice, qslice_count)
	}
	qslice = append(qslice, ";")
	query := strings.Join(qslice, "")
	rows, err := db.DB.Query(query)
	if err != nil {
		return 0
	}
	defer rows.Close()
	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println(query)
	return int64(count)
}
