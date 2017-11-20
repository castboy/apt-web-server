/********SSLCert统计********/
package ssl

import (
	"apt-web-server_v2/models/db"
	"fmt"
	"strings"
	//"time"
)

func (this *TblSSLCLst) GetTableName() string {
	return "cert_data"
	//return "cert_data_tmp"
}

func GetSSLCLstMysqlCmd(datemold, tablename string, para *TblSSLCLstSearchPara) []string {
	//var flag int
	var qslice_info string
	qslice := make([]string, 0)

	switch para.Type {
	case "serialnum":
		sType := "serial_num"
		qslice_info = fmt.Sprintf(`SELECT CommonName, org_name, OUName, serial_num,
			FROM_UNIXTIME(not_before,'%s'), FROM_UNIXTIME(not_after,'%s'), cert_ver 
			FROM %s WHERE %s='%s'`,
			datemold, datemold, tablename, sType, para.SerialNum)
	case "comnname":
		sType := "CommonName"
		qslice_info = fmt.Sprintf(`SELECT CommonName, org_name, OUName, serial_num, 
			FROM_UNIXTIME(not_before,'%s'), FROM_UNIXTIME(not_after,'%s'), cert_ver 
			FROM %s WHERE %s LIKE '%s%s%s'`,
			datemold, datemold, tablename, sType, "%", para.ComnName, "%")
	default:
		return qslice
	}

	qslice = append(qslice, qslice_info)

	if para.LastCount != 0 {
		temp_LC := fmt.Sprintf(" limit %d", para.LastCount)
		qslice = append(qslice, temp_LC)
	} else if para.Count != 0 {
		para.Page = para.Page - 1
		temp_PC := fmt.Sprintf(" limit %d,%d", para.Page*para.Count, para.Count)
		qslice = append(qslice, temp_PC)
	}
	qslice = append(qslice, ";")

	//fmt.Println("DBY the mysql cmd is ", qslice)
	return qslice
}

func (this *TblSSLCLst) GetSSLCLst(para *TblSSLCLstSearchPara) (error, *TblSSLCLstData) {
	//datemold := GetDateMold(para.Unit)  datemold = "%Y-%m-%d"
	datemold := "%Y-%m-%d %H:%i:%s"
	var cnt int64 = 0
	list := TblSSLCLstData{}

	qslice := GetSSLCLstMysqlCmd(datemold, this.GetTableName(), para)
	query := strings.Join(qslice, "")
	//fmt.Println("GetSSLCLstMysqlCmd get mysql cmd qslice is ",qslice)
	//fmt.Println("GetSSLCLstMysqlCmd get mysql cmd query is ",query)

	rows, err := db.DB.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	ugc := new(TblSSLCLst)
	for rows.Next() {
		err = rows.Scan(
			&ugc.SSLCertComnName,
			&ugc.SSLCertOrigName,
			&ugc.SSLCertUnitName,
			&ugc.SSLCertSerialNum,
			&ugc.SSLCertNotBefore,
			&ugc.SSLCertNotAfter,
			&ugc.SSLCertVersion)
		if err != nil {
			return err, nil
		}
		cnt++
		list.Elements = append(list.Elements, ugc.TblSSLCLstContent)
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	list.Counts = cnt //GetSSLCount("", para)
	if cnt != 0 {
		//list.Totality = GetSSLLstCount(this.GetTableName(), para)
		list.Totality = GetSSLLstCount(this.GetTableName(), para)
	}

	return nil, &list
}

func GetSSLLstCount(tablename string, para *TblSSLCLstSearchPara) int64 {
	qslice := make([]string, 0)
	var totalCnt int64
	var qslice_info string

	switch para.Type {
	case "serialnum":
		sType := "serial_num"
		qslice_info = fmt.Sprintf(`SELECT count(*) 
			FROM %s WHERE %s='%s'`,
			tablename, sType, para.SerialNum)
	case "comnname":
		sType := "CommonName"
		qslice_info = fmt.Sprintf(`SELECT count(*) 
			FROM %s WHERE %s LIKE '%s%s%s'`,
			tablename, sType, "%", para.ComnName, "%")
	default:
		fmt.Println("GetSSLLstCount , input para cmd error : ", para.Type)
		return 0
	}

	qslice = append(qslice, qslice_info)
	qslice = append(qslice, ";")
	query := strings.Join(qslice, "")

	rows, err := db.DB.Query(query)
	if err != nil {
		return 0
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&totalCnt)
		if err != nil {
			panic(err)
		}
	}
	//fmt.Println("GetSSLLstCount mysql cmd:", query)
	return int64(totalCnt)
}

func GetSSLLstCountOnCondition(tablename string, para *TblSSLCLstSearchPara) int64 {
	qslice := make([]string, 0)
	var totalCnt int64

	qslice_count := fmt.Sprintf(`select count(*) as totalnum from %s`,
		tablename)
	qslice = append(qslice, qslice_count)
	qslice = append(qslice, ";")
	query := strings.Join(qslice, "")

	rows, err := db.DB.Query(query)
	if err != nil {
		return 0
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&totalCnt)
		if err != nil {
			panic(err)
		}
	}
	//fmt.Println("GetSSLLstCountOnCondition mysql cmd:", query)
	return int64(totalCnt)
}
