/********获取紧急事件分类数********/
package models

import (
	"fmt"
	"strings"
	"time"
)

func (this *TblUgcMCount) TableName() string {
	return "urgencymold"
}

func GetUgcMCountMysqlCmd(sqltype string, para *TblUgcMCountSearchPara) []string {
	//var flag int
	qslice := make([]string, 0)
	if sqltype == "urgency" {
		qslice_tmp := fmt.Sprintf(`select sum(case when attack_type='webshell' then 1 else 0 end) as webshell,sum(case when attack_type='exceptionalvisit' then 1 else 0 end) as exceptionalvisit,sum(case when attack_type='abnormal_connection' then 1 else 0 end) as abnormal_connection from urgencymold `)
		qslice = append(qslice, qslice_tmp)
	}
	if sqltype == "bruteforce" {
		qslice_tmp := fmt.Sprintf(`select sum(count) from bruteForceAction `)
		qslice = append(qslice, qslice_tmp)
	}
	if sqltype == "portscan" {
		qslice_tmp := fmt.Sprintf(`select sum(count) from alert_portscan `)
		qslice = append(qslice, qslice_tmp)
	}
	if para.PField.Start != 0 && para.PField.End != 0 {
		var qslice_time string
		if sqltype == "bruteforce" {
			tmStart := time.Unix(para.PField.Start, 0)
			tmEnd := time.Unix(para.PField.End, 0)
			qslice_time = fmt.Sprintf("where time between '%s' and '%s' ", tmStart.Format("2006-01-02 03:04:05 PM"), tmEnd.Format("2006-01-02 03:04:05 PM"))
		} else {
			qslice_time = fmt.Sprintf(" where time between %d and %d ", para.PField.Start, para.PField.End)
		}
		qslice = append(qslice, qslice_time)
		//flag = 1
	}
	qslice = append(qslice, ";")
	//fmt.Println(tablename, para)
	fmt.Println(qslice)
	return qslice
}

func (this *TblUgcMCount) GetUgcMCount(para *TblUgcMCountSearchPara) (error, *TblUgcMCountData) {
	var i int64
	start := para.PField.Start
	end := para.PField.End
	seconds := GetUgcSeconds(para.Unit)
	list := TblUgcMCountData{}
	for i = 0; i <= ((end - start) / seconds); i++ {
		para.PField.Start = start + seconds*i
		para.PField.End = start + seconds*(i+1) - 1
		if para.PField.End >= end {
			para.PField.End = end
		}
		ugc := new(TblUgcMCount)
		qslice_u := GetUgcMCountMysqlCmd("urgency", para)
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
		qslice_b := GetUgcMCountMysqlCmd("bruteforce", para)
		query_b := strings.Join(qslice_b, "")
		fmt.Println(qslice_u)
		rows_b, err := db.Query(query_b)
		if err != nil {
			return err, nil
		}
		defer rows_b.Close()
		for rows_b.Next() {
			err = rows_b.Scan(
				&ugc.BruteForce)
			if err != nil {
				ugc.BruteForce = 0
				//return err, nil
			}
		}
		if err := rows_b.Err(); err != nil {
			return err, nil
		}
		qslice_p := GetUgcMCountMysqlCmd("portscan", para)
		query_p := strings.Join(qslice_p, "")
		fmt.Println(qslice_u)
		rows_p, err := db.Query(query_p)
		if err != nil {
			return err, nil
		}
		defer rows_p.Close()
		for rows_p.Next() {
			err = rows_p.Scan(
				&ugc.PortScan)
			if err != nil {
				ugc.PortScan = 0
				//return err, nil
			}
		}
		if err := rows_p.Err(); err != nil {
			return err, nil
		}
		if start == 0 && end == 0 {
			ugc.Time = int64(time.Now().Unix())
		} else {
			ugc.Time = para.PField.Start
		}

		list.Elements = append(list.Elements, TblUgcMCountList{ugc.TblUgcMCountContent})
	}
	list.Counts = i
	return nil, &list
}
