/********获取紧急事件分类数********/
package models

import (
	"apt-web-server/modules/mlog"
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
		qslice_tmp := fmt.Sprintf(`select sum(case when attack_type='webshell' then 1 else 0 end) as webshell,
		sum(case when attack_type='exceptionalvisit' then 1 else 0 end) as exceptionalvisit,
		sum(case when attack_type='abnormal_connection' then 1 else 0 end) as abnormal_connection from urgencymold`)
		qslice = append(qslice, qslice_tmp)
	}
	if sqltype == "alert_waf" {
		qslice_tmp := fmt.Sprintf(`select sum(case when attack='sqli' then 1 else 0 end) as sqli,
		sum(case when attack='xss' then 1 else 0 end) as xss,
		sum(case when attack='injection_php' then 1 else 0 end) as injection_php,
		sum(case when attack='rfi' then 1 else 0 end) as rfi from alert_waf where severity in (0,1,2)`)
		qslice = append(qslice, qslice_tmp)
		flag = 1
	}
	if para.PField.Start != 0 && para.PField.End != 0 {
		if flag != 0 {
			qslice_tmp := fmt.Sprintf(` and time between %d and %d`, para.PField.Start, para.PField.End)
			qslice = append(qslice, qslice_tmp)
		} else {
			qslice_tmp := fmt.Sprintf(` where time between %d and %d`, para.PField.Start, para.PField.End)
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
			qslice_tmp := fmt.Sprintf(`select attack_type,severity,count(attack_type),time from urgencymold`)
			qslice = append(qslice, qslice_tmp)
		} else {
			qslice_tmp := fmt.Sprintf(`select time,src_ip,dest_ip,attack_type,severity,details from urgencymold`)
			qslice = append(qslice, qslice_tmp)
		}
		if para.Type != "" {
			qslice_tmp := fmt.Sprintf(` where attack_type='%s'`, para.Type)
			qslice = append(qslice, qslice_tmp)
			flag = 1
		}
		if para.PField.Start != 0 && para.PField.End != 0 {
			if flag == 1 {
				qslice_tmp := fmt.Sprintf(` and time between %d and %d`, para.PField.Start, para.PField.End)
				qslice = append(qslice, qslice_tmp)
			} else {
				qslice_tmp := fmt.Sprintf(` where time between %d and %d`, para.PField.Start, para.PField.End)
				qslice = append(qslice, qslice_tmp)
				flag = 1
			}
		}
	}

	if para.Type == "" {
		qslice_tmp := fmt.Sprintf(` union all`)
		qslice = append(qslice, qslice_tmp)
	}
	if para.Type == "xss" ||
		para.Type == "injection_php" ||
		para.Type == "rfi" ||
		para.Type == "sqli" ||
		para.Type == "" {
		if para.Unit != "" {
			qslice_tmp := fmt.Sprintf(` select attack,severity,count(attack),time from alert_waf`)
			qslice = append(qslice, qslice_tmp)
		} else {
			qslice_tmp := fmt.Sprintf(` select time,client,hostname,attack,severity,rule_data from alert_waf`)
			qslice = append(qslice, qslice_tmp)
		}
		if para.Type == "" {
			qslice_tmp := fmt.Sprintf(` where attack in ('sqli','xss','injection_php','rfi') and severity in (0,1,2)`)
			qslice = append(qslice, qslice_tmp)
		} else {
			qslice_tmp := fmt.Sprintf(` where attack='%s' and severity in (0,1,2)`, para.Type)
			qslice = append(qslice, qslice_tmp)
			flag = 1
		}
		if para.PField.Start != 0 && para.PField.End != 0 {
			if flag == 1 {
				qslice_tmp := fmt.Sprintf(` and time between %d and %d`, para.PField.Start, para.PField.End)
				qslice = append(qslice, qslice_tmp)
			} else {
				qslice_tmp := fmt.Sprintf(` where time between %d and %d`, para.PField.Start, para.PField.End)
				qslice = append(qslice, qslice_tmp)
				flag = 1
			}
		}
	}

	if para.Sort != "" {
		temp_s := fmt.Sprintf(" order by %s %s,severity desc", para.Sort, para.Order)
		qslice = append(qslice, temp_s)
	} else {
		temp_s := fmt.Sprintf(" order by severity")
		qslice = append(qslice, temp_s)
	}
	if para.LastCount != 0 {
		temp_LC := fmt.Sprintf(" limit %d", para.LastCount)
		qslice = append(qslice, temp_LC)
	} else if para.Count != 0 {
		para.Page = para.Page - 1
		temp_PC := fmt.Sprintf(" limit %d,%d", para.Page*para.Count, para.Count)
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
	rows, err := db.Query(query)
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
		rows_u, err := db.Query(query_u)
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
		rows_a, err := db.Query(query_a)
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
		qslice_tmp := fmt.Sprintf(`select count(id) from %s`, tablename)
		qslice = append(qslice, qslice_tmp)
		if tablename == "alert_waf" {
			if para.Type != "" {
				temp_t := fmt.Sprintf(" where attack='%s'", para.Type)
				qslice = append(qslice, temp_t)
			} else {
				qslice_waf := fmt.Sprintf(" where attack in ('sqli','xss','injection_php','rfi') and severity<=2")
				qslice = append(qslice, qslice_waf)
			}
			whereflag = 1
		}
		if tablename == "urgencymold" {
			if para.Type != "" {
				if whereflag != 0 {
					temp_t := fmt.Sprintf(" and attack_type='%s'", para.Type)
					qslice = append(qslice, temp_t)
				} else {
					temp_t := fmt.Sprintf(" where attack_type='%s'", para.Type)
					qslice = append(qslice, temp_t)
					whereflag = 1
				}
			}
		}

		if para.PField.Start != 0 && para.PField.End != 0 {
			if whereflag == 1 {
				qslice_tmp := fmt.Sprintf(` and time between %d and %d`, para.PField.Start, para.PField.End)
				qslice = append(qslice, qslice_tmp)
			} else {
				qslice_tmp := fmt.Sprintf(` where time between %d and %d`, para.PField.Start, para.PField.End)
				qslice = append(qslice, qslice_tmp)
				whereflag = 1
			}
		}
	} else {
		qslice_count := fmt.Sprintf(`select FOUND_ROWS() as count`)
		qslice = append(qslice, qslice_count)
	}
	qslice = append(qslice, ";")
	query := strings.Join(qslice, "")
	rows, err := db.Query(query)
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

func (this *TblUgcClassify) CreateSql() string {
	return fmt.Sprintf(
		`CREATE TABLE %s (
		id   integer unsigned  AUTO_INCREMENT NOT NULL,
		time   BIGINT NOT NULL DEFAULT 0,
		ugcType varchar(20) NOT NULL DEFAULT '',
		serverIp varchar(20) NOT NULL DEFAULT '',
		description varchar(50) NOT NULL DEFAULT '',
		PRIMARY KEY (Id)
	)ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
		this.TableName())
}
