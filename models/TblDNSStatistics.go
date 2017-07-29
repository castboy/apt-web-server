/********DNS统计********/
package models

import (
	"fmt"
	"strings"
	"time"
)

func (this *TblDNSS) TableName() string {
	return "DNS_day"
}
func GetDateMold(unit string) string {
	var datemold string
	switch unit {
	case "day":
		datemold = "%Y-%m-%d"
	case "week":
		datemold = "%Y %u"
	case "month":
		datemold = "%Y-%m"
	default:
		datemold = "%Y-%m-%d"
	}
	return datemold
}
func GetDNSSMysqlCmd(datemold, tablename string, para *TblDNSSSearchPara) []string {
	var flag int
	var qslice_count, groupType string
	qslice := make([]string, 0)
	qslice_tbl := fmt.Sprintf(`(SELECT DATE_FORMAT(time,'%s') AS time,ip,
	domain,count FROM %s) AS a`,
		datemold, tablename)
	switch para.Type {
	case "ip":
		groupType = "domain"
		qslice_count = fmt.Sprintf(`SELECT time,domain,count(ip) AS count FROM %s`,
			qslice_tbl)
	case "domain":
		groupType = "domain"
		qslice_count = fmt.Sprintf(`SELECT time,domain,sum(count) AS count FROM %s`,
			qslice_tbl)
	default:
		qslice_count = fmt.Sprintf(`SELECT DATE_FORMAT(time,'%s') AS time,ip,
		domain,count FROM %s`,
			datemold, tablename)
	}

	qslice = append(qslice, qslice_count)
	if para.PField.Start != 0 {
		if para.PField.End == 0 {
			para.PField.End = time.Now().Unix()
		}
		qslice_time := fmt.Sprintf(" WHERE time BETWEEN FROM_UNIXTIME(%d,'%s') AND FROM_UNIXTIME(%d,'%s')",
			para.PField.Start, datemold, para.PField.End, datemold)
		qslice = append(qslice, qslice_time)
		flag = 1
	}
	if para.Ip != "" {
		if flag != 0 {
			qslice_ip := fmt.Sprintf(" AND ip in ('%s')", para.Ip)
			qslice = append(qslice, qslice_ip)
		} else {
			qslice_ip := fmt.Sprintf(" WHERE ip in ('%s')", para.Ip)
			qslice = append(qslice, qslice_ip)
			flag = 1
		}
	}
	if para.Domain != "" {
		if flag != 0 {
			qslice_domain := fmt.Sprintf(" AND domain in ('%s')", para.Domain)
			qslice = append(qslice, qslice_domain)
		} else {
			qslice_domain := fmt.Sprintf(" WHERE domain in ('%s')", para.Domain)
			qslice = append(qslice, qslice_domain)
			flag = 1
		}
	}
	if para.Type != "" {
		qslice_type := fmt.Sprintf(` group by %s`, groupType)
		qslice = append(qslice, qslice_type)
	}
	if para.Sort != "" {
		temp_s := fmt.Sprintf(" order by %s %s", para.Sort, para.Order)
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
	return qslice
}

func (this *TblDNSS) GetDNSS(para *TblDNSSSearchPara) (error, *TblDNSSData) {
	datemold := GetDateMold(para.Unit)
	list := TblDNSSData{}
	qslice := GetDNSSMysqlCmd(datemold, this.TableName(), para)
	query := strings.Join(qslice, "")
	fmt.Println(qslice)
	rows, err := db.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	ugc := new(TblDNSS)
	for rows.Next() {
		err = rows.Scan(
			&ugc.Time,
			&ugc.Ip,
			&ugc.Domain,
			&ugc.Count)
		if err != nil {
			return err, nil
		}
		list.Elements = append(list.Elements, TblDNSSList{ugc.TblDNSSContent})
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	list.Counts = GetDNSCount("", para)
	list.Totality = GetDNSCount(this.TableName(), para)

	return nil, &list
}
func (this *TblDNSS) GetDNSSIp(para *TblDNSSSearchPara) (error, *TblDNSSIpData) {
	datemold := GetDateMold(para.Unit)
	list := TblDNSSIpData{}
	qslice := GetDNSSMysqlCmd(datemold, this.TableName(), para)
	query := strings.Join(qslice, "")
	fmt.Println(qslice)

	rows, err := db.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	ugc := new(TblDNSSIp)
	for rows.Next() {
		err = rows.Scan(
			&ugc.Time,
			&ugc.Domain,
			&ugc.Count)
		if err != nil {
			return err, nil
		}
		list.Elements = append(list.Elements, TblDNSSIpList{ugc.TblDNSSIpContent})
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	list.Counts = GetDNSCount("", para)
	list.Totality = GetDNSCount(this.TableName(), para)
	return nil, &list
}
func (this *TblDNSS) GetDNSSDomain(para *TblDNSSSearchPara) (error, *TblDNSSDomainData) {
	datemold := GetDateMold(para.Unit)
	list := TblDNSSDomainData{}
	qslice := GetDNSSMysqlCmd(datemold, this.TableName(), para)
	query := strings.Join(qslice, "")
	fmt.Println(qslice)

	rows, err := db.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	ugc := new(TblDNSSDomain)
	for rows.Next() {
		err = rows.Scan(
			&ugc.Time,
			&ugc.Domain,
			&ugc.Count)
		if err != nil {
			return err, nil
		}
		list.Elements = append(list.Elements, TblDNSSDomainList{ugc.TblDNSSDomainContent})
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	list.Counts = GetDNSCount("", para)
	list.Totality = GetDNSCount(this.TableName(), para)
	return nil, &list
}
func GetDNSCount(tablename string, para *TblDNSSSearchPara) int64 {
	var groupType string
	qslice := make([]string, 0)
	datemold := GetDateMold(para.Unit)
	if para.Type == "ip" {
		groupType = "domain"
	} else {
		groupType = para.Type
	}
	qslice_tbl := fmt.Sprintf(`(SELECT DATE_FORMAT(time,'%s') AS time,ip,
	domain,count FROM %s) AS a`,
		datemold, tablename)
	qslice_count := fmt.Sprintf(`select FOUND_ROWS() as count`)
	qslice_total_type := fmt.Sprintf(`select %s from %s`, groupType, qslice_tbl)
	qslice_total := fmt.Sprintf(`SELECT count(time) from %s`, tablename)
	if tablename != "" {
		if para.Type != "" {
			qslice = append(qslice, qslice_total_type)
		} else {
			qslice = append(qslice, qslice_total)
		}

		if para.PField.Start != 0 {
			if para.PField.End == 0 {
				para.PField.End = para.PField.Start
			}
			qslice_time := fmt.Sprintf(" WHERE time BETWEEN FROM_UNIXTIME(%d,'%s') AND FROM_UNIXTIME(%d,'%s') ",
				para.PField.Start, datemold, para.PField.End, datemold)
			qslice = append(qslice, qslice_time)
			//flag = 1
		}
		if para.Type != "" {
			qslice_tmp := fmt.Sprintf(" group by %s", groupType)
			qslice = append(qslice, qslice_tmp)
			qslice_tc := fmt.Sprintf("select count(%s) from (%s) as aa", groupType, strings.Join(qslice, ""))
			qslice_catch := make([]string, 0)
			qslice_catch = append(qslice_catch, qslice_tc)
			qslice = qslice_catch
		}

	} else {
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
	return int64(count)

}
