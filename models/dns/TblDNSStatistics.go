/********DNS统计********/
package dns

import (
	"apt-web-server_v2/models/db"
	"fmt"
	"strings"
	//"time"
)

func (this *TblDNSS) TableName() string {
	return "dns_day"
}

func GetDNSSMysqlCmd(datemold, tablename string, para *TblDNSSSearchPara) []string {
	var flag int
	var qslice_count, groupType string
	qslice := make([]string, 0)
	qslice_tbl := fmt.Sprintf(`(SELECT DATE_FORMAT(time,'%s') AS time,ip,
		domain,count FROM %s WHERE time BETWEEN FROM_UNIXTIME(%d,'%s') 
		AND FROM_UNIXTIME(%d,'%s')) AS a`,
		datemold, tablename, para.PField.Start, datemold, para.PField.End, datemold)
	qslice_ipnum := fmt.Sprintf(`(SELECT time,domain,ip FROM %s 
		group by domain,ip) AS aa`,
		qslice_tbl)
	switch para.Type {
	case "ip":
		groupType = "domain"
		qslice_count = fmt.Sprintf(`SELECT time,domain,COUNT(ip) AS count FROM %s`,
			qslice_ipnum)
	case "domain":
		groupType = "domain"
		qslice_count = fmt.Sprintf(`SELECT time,domain,SUM(count) AS count FROM %s`,
			qslice_tbl)
	default:
		qslice_count = fmt.Sprintf(`SELECT DATE_FORMAT(time,'%s') AS time,ip,
		domain,count FROM %s`,
			datemold, tablename)
	}

	qslice = append(qslice, qslice_count)
	if para.PField.Start != 0 {
		qslice_time := fmt.Sprintf(" WHERE time BETWEEN FROM_UNIXTIME(%d,'%s') AND FROM_UNIXTIME(%d,'%s')",
			para.PField.Start, datemold, para.PField.End, datemold)
		if para.Type == "" {
			qslice = append(qslice, qslice_time)
			flag = 1
		}
	}
	if para.Ip != "" {
		if flag != 0 {
			qslice_ip := fmt.Sprintf(" AND ip IN ('%s')", para.Ip)
			qslice = append(qslice, qslice_ip)
		} else {
			qslice_ip := fmt.Sprintf(" WHERE ip IN ('%s')", para.Ip)
			qslice = append(qslice, qslice_ip)
			flag = 1
		}
	}
	if para.Domain != "" {
		if flag != 0 {
			qslice_domain := fmt.Sprintf(" AND domain IN ('%s')", para.Domain)
			qslice = append(qslice, qslice_domain)
		} else {
			qslice_domain := fmt.Sprintf(" WHERE domain IN ('%s')", para.Domain)
			qslice = append(qslice, qslice_domain)
			flag = 1
		}
	}
	if para.Type != "" {
		qslice_type := fmt.Sprintf(` GROUP BY %s`, groupType)
		qslice = append(qslice, qslice_type)
	}
	if para.Sort != "" {
		temp_s := fmt.Sprintf(" ORDER BY %s %s", para.Sort, para.Order)
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
	return qslice
}

func (this *TblDNSS) GetDNSS(para *TblDNSSSearchPara) (error, *TblDNSSData) {
	//datemold := GetDateMold(para.Unit)
	datemold := "%Y-%m-%d"
	list := TblDNSSData{}
	qslice := GetDNSSMysqlCmd(datemold, this.TableName(), para)
	query := strings.Join(qslice, "")
	fmt.Println(qslice)
	rows, err := db.DB.Query(query)
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
	//datemold := GetDateMold(para.Unit)
	datemold := "%Y-%m-%d"
	list := TblDNSSIpData{}
	qslice := GetDNSSMysqlCmd(datemold, this.TableName(), para)
	query := strings.Join(qslice, "")
	fmt.Println(qslice)

	rows, err := db.DB.Query(query)
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
	//datemold := GetDateMold(para.Unit)
	datemold := "%Y-%m-%d"
	list := TblDNSSDomainData{}
	qslice := GetDNSSMysqlCmd(datemold, this.TableName(), para)
	query := strings.Join(qslice, "")
	fmt.Println(qslice)

	rows, err := db.DB.Query(query)
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
	//datemold := GetDateMold(para.Unit)
	datemold := "%Y-%m-%d"
	if para.Type == "ip" {
		groupType = "domain"
	} else {
		groupType = para.Type
	}
	qslice_tbl := fmt.Sprintf(`(SELECT DATE_FORMAT(time,'%s') AS time,ip,
	    domain,count FROM %s) AS a`, datemold, tablename)
	qslice_count := fmt.Sprintf(`SELECT FOUND_ROWS() as count`)
	qslice_total_type := fmt.Sprintf(`SELECT %s,ip FROM %s`, groupType, qslice_tbl)
	qslice_total := fmt.Sprintf(`SELECT COUNT(time) FROM %s`, tablename)
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
			qslice_time := fmt.Sprintf(` WHERE time BETWEEN FROM_UNIXTIME(%d,'%s') 
			    AND FROM_UNIXTIME(%d,'%s')`,
				para.PField.Start, datemold, para.PField.End, datemold)
			qslice = append(qslice, qslice_time)
			//flag = 1
		}
		if para.Type != "" {
			qslice_tmp := fmt.Sprintf(` GROUP BY %s`, groupType)
			qslice = append(qslice, qslice_tmp)
			qslice_tc := fmt.Sprintf(`SELECT COUNT(%s) FROM (%s) AS aa`,
				groupType, strings.Join(qslice, ""))
			qslice_catch := make([]string, 0)
			qslice_catch = append(qslice_catch, qslice_tc)
			qslice = qslice_catch
		}

	} else {
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
	fmt.Println("count cmd:", qslice)
	return int64(count)

}
