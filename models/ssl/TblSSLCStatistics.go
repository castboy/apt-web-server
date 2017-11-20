/********SSLCert统计********/
package ssl

import (
	"apt-web-server_v2/models/db"
	"fmt"
	"strings"
	//"time"
)

func (this *TblSSLCS) TableName() string {
	return "ssl_perday"
}

func GetSSLCSMysqlCmd(datemold, tablename string, para *TblSSLCSSearchPara) []string {
	var flag int
	var qslice_count, groupType string
	qslice := make([]string, 0)
	qslice_tbl := fmt.Sprintf(`(SELECT DATE_FORMAT(time,'%s') AS time,cli_ip as ip,
		s_cert_comnname,s_cert_unitname,s_cert_serialnum,s_signsummary,count 
		FROM %s WHERE time BETWEEN FROM_UNIXTIME(%d,'%s') 
		AND FROM_UNIXTIME(%d,'%s') AND s_verify=%d) AS a`,
		datemold, tablename, para.PField.Start, datemold,
		para.PField.End, datemold, para.Verify)
	qslice_ipnum := fmt.Sprintf(`(SELECT time,ip,s_cert_comnname,s_cert_unitname,
		s_cert_serialnum,s_signsummary FROM %s GROUP BY s_signsummary,ip) AS aa`,
		qslice_tbl)
	switch para.Type {
	case "ip":
		groupType = "s_signsummary"
		qslice_count = fmt.Sprintf(`SELECT time,s_cert_comnname,s_cert_unitname,
		s_cert_serialnum,count(ip) AS count FROM %s`,
			qslice_ipnum)
	case "cert":
		groupType = "s_signsummary"
		qslice_count = fmt.Sprintf(`SELECT time,s_cert_comnname,s_cert_unitname,
		s_cert_serialnum,sum(count) AS count FROM %s`,
			qslice_tbl)
	default:
		qslice_count = fmt.Sprintf(`SELECT DATE_FORMAT(time,'%s') AS time,cli_ip,
		s_cert_comnname,s_cert_unitname,s_cert_serialnum,count FROM %s`,
			datemold, tablename)
	}
	qslice = append(qslice, qslice_count)
	if para.PField.Start != 0 {
		qslice_time := fmt.Sprintf(` WHERE time BETWEEN FROM_UNIXTIME(%d,'%s') 
		    AND FROM_UNIXTIME(%d,'%s')`,
			para.PField.Start, datemold, para.PField.End, datemold)
		if para.Type != "ip" {
			qslice = append(qslice, qslice_time)
			flag = 1
		}
	}
	if para.Ip != "" {
		var qslice_ip string
		if flag != 0 {
			qslice_ip = fmt.Sprintf(" AND ip IN ('%s')", para.Ip)
		} else {
			qslice_ip = fmt.Sprintf(" WHERE ip IN ('%s')", para.Ip)
			flag = 1
		}
		qslice = append(qslice, qslice_ip)
	}
	if para.ComnName != "" {
		var qslice_cn string
		if flag != 0 {
			qslice_cn = fmt.Sprintf(" AND s_cert_comnname IN ('%s')", para.ComnName)
		} else {
			qslice_cn = fmt.Sprintf(" WHERE s_cert_comnname IN ('%s')", para.ComnName)
			flag = 1
		}
		qslice = append(qslice, qslice_cn)
	}
	if para.UnitName != "" {
		var qslice_un string
		if flag != 0 {
			qslice_un = fmt.Sprintf(" AND s_cert_unitname IN ('%s')", para.UnitName)
		} else {
			qslice_un = fmt.Sprintf(" WHERE s_cert_unitname IN ('%s')", para.UnitName)
			flag = 1
		}
		qslice = append(qslice, qslice_un)
	}
	if para.SerialNum != "" {
		var qslice_sn string
		if flag != 0 {
			qslice_sn = fmt.Sprintf(" AND s_cert_serialnum IN ('%s')", para.SerialNum)
		} else {
			qslice_sn = fmt.Sprintf(" WHERE s_cert_serialnum IN ('%s')", para.SerialNum)
			flag = 1
		}
		qslice = append(qslice, qslice_sn)
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

func (this *TblSSLCS) GetSSLCS(para *TblSSLCSSearchPara) (error, *TblSSLCSData) {
	//datemold := GetDateMold(para.Unit)
	datemold := "%Y-%m-%d"
	list := TblSSLCSData{}
	qslice := GetSSLCSMysqlCmd(datemold, this.TableName(), para)
	query := strings.Join(qslice, "")
	fmt.Println(qslice)
	rows, err := db.DB.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	ugc := new(TblSSLCS)
	for rows.Next() {
		err = rows.Scan(
			&ugc.Time,
			&ugc.Ip,
			&ugc.SSLCertComnName,
			&ugc.SSLCertUnitName,
			&ugc.SSLCertSerialNum,
			&ugc.Count)
		if err != nil {
			return err, nil
		}
		list.Elements = append(list.Elements, TblSSLCSList{ugc.TblSSLCSContent})
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	list.Counts = GetSSLCount("", para)
	list.Totality = GetSSLCount(this.TableName(), para)

	return nil, &list
}
func (this *TblSSLCS) GetSSLCSIp(para *TblSSLCSSearchPara) (error, *TblSSLCSIpData) {
	//datemold := GetDateMold(para.Unit)
	datemold := "%Y-%m-%d"
	list := TblSSLCSIpData{}
	qslice := GetSSLCSMysqlCmd(datemold, this.TableName(), para)
	query := strings.Join(qslice, "")
	fmt.Println(qslice)

	rows, err := db.DB.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	ugc := new(TblSSLCSIp)
	for rows.Next() {
		err = rows.Scan(
			&ugc.Time,
			&ugc.SSLCertComnName,
			&ugc.SSLCertUnitName,
			&ugc.SSLCertSerialNum,
			&ugc.Count)
		if err != nil {
			return err, nil
		}
		list.Elements = append(list.Elements, TblSSLCSIpList{ugc.TblSSLCSIpContent})
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	list.Counts = GetSSLCount("", para)
	list.Totality = GetSSLCount(this.TableName(), para)
	return nil, &list
}
func (this *TblSSLCS) GetSSLCSCert(para *TblSSLCSSearchPara) (error, *TblSSLCSCertData) {
	//datemold := GetDateMold(para.Unit)
	datemold := "%Y-%m-%d"
	list := TblSSLCSCertData{}
	qslice := GetSSLCSMysqlCmd(datemold, this.TableName(), para)
	query := strings.Join(qslice, "")
	fmt.Println(qslice)

	rows, err := db.DB.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	ugc := new(TblSSLCSCert)
	for rows.Next() {
		err = rows.Scan(
			&ugc.Time,
			&ugc.SSLCertComnName,
			&ugc.SSLCertUnitName,
			&ugc.SSLCertSerialNum,
			&ugc.Count)
		if err != nil {
			return err, nil
		}
		list.Elements = append(list.Elements, TblSSLCSCertList{ugc.TblSSLCSCertContent})
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	list.Counts = GetSSLCount("", para)
	list.Totality = GetSSLCount(this.TableName(), para)
	return nil, &list
}
func GetSSLCount(tablename string, para *TblSSLCSSearchPara) int64 {
	var groupType string
	qslice := make([]string, 0)
	//datemold := GetDateMold(para.Unit)
	datemold := "%Y-%m-%d"
	if para.Type == "ip" || para.Type == "cert" {
		groupType = "s_signsummary"
	} else {
		groupType = para.Type
	}
	qslice_tbl := fmt.Sprintf(`(SELECT DATE_FORMAT(time,'%s') AS time,cli_ip as ip,
	s_cert_comnname,s_cert_unitname,s_cert_serialnum,s_signsummary,count 
	FROM %s WHERE s_verify=%d) AS a`,
		datemold, tablename, para.Verify)
	qslice_count := fmt.Sprintf(`SELECT FOUND_ROWS() AS count`)
	qslice_total_type := fmt.Sprintf(`SELECT %s FROM %s`, groupType, qslice_tbl)
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
			qslice_time := fmt.Sprintf(" WHERE time BETWEEN FROM_UNIXTIME(%d,'%s') AND FROM_UNIXTIME(%d,'%s') ",
				para.PField.Start, datemold, para.PField.End, datemold)
			qslice = append(qslice, qslice_time)
			//flag = 1
		}
		if para.Type != "" {
			qslice_tmp := fmt.Sprintf(" GROUP BY %s", groupType)
			qslice = append(qslice, qslice_tmp)
			qslice_tc := fmt.Sprintf("SELECT COUNT(%s) FROM (%s) AS aa", groupType, strings.Join(qslice, ""))
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
	fmt.Println("count mysql cmd:", query)
	return int64(count)

}
