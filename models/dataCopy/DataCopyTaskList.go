package dataCopy

import (
	"apt-web-server_v2/models/db"
	"apt-web-server_v2/models/modelsPublic"
	"apt-web-server_v2/modules/mlog"
	"fmt"
	"strings"
)

func GetDCTLMsqlCmd(tablename string, para *TblDCTLSearchPara) []string {
	qslice := make([]string, 0)
	qslice_tmp := fmt.Sprintf(`SELECT location,startday,endday,time,status,
	    diskpath,details FROM %s `, tablename)
	qslice = append(qslice, qslice_tmp)
	flag := 0
	if para.PField.Start != 0 && para.PField.End != 0 {
		qslice_time := fmt.Sprintf(" WHERE time BETWEEN '%d' AND '%d' ",
			para.PField.Start, para.PField.End)
		qslice = append(qslice, qslice_time)
		flag = 1
	}
	if para.Location != "" {
		var tmp_location string
		if flag == 0 {
			tmp_location = fmt.Sprintf(" WHERE location IN ('%s')", para.Location)
			flag = 1
		} else {
			tmp_location = fmt.Sprintf(" AND location IN ('%s')", para.Location)
		}
		qslice = append(qslice, tmp_location)
	}
	if para.Date != 0 {
		var tmp_date string
		if flag == 0 {
			tmp_date = fmt.Sprintf(" WHERE date IN ('%d')", para.Date)
			flag = 1
		} else {
			tmp_date = fmt.Sprintf(" AND date IN ('%d')", para.Date)
		}
		qslice = append(qslice, tmp_date)
	}
	if para.Status != "" {
		var tmp_status string
		if flag == 0 {
			tmp_status = fmt.Sprintf(" WHERE status IN ('%s')", para.Status)
			flag = 1
		} else {
			tmp_status = fmt.Sprintf(" AND status IN ('%s')", para.Status)
		}
		qslice = append(qslice, tmp_status)
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
	fmt.Println(para)
	return qslice
}

func (this *TblDCT) GetDCTL(para *TblDCTLSearchPara) (error, *TblDCTLData) {
	qslice := GetDCTLMsqlCmd(this.TableName(), para)
	query := strings.Join(qslice, "")
	mlog.Debug(string(query))
	rows, err := db.DB.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()

	list := TblDCTLData{}
	for rows.Next() {
		ugc := new(TblDCT)
		err = rows.Scan(
			&ugc.Location,
			&ugc.DataStart,
			&ugc.DataEnd,
			&ugc.Date,
			&ugc.Status,
			&ugc.DiskPath,
			&ugc.Details)
		if err != nil {
			return err, nil
		}
		list.Elements = append(list.Elements, TblDCTList{ugc.DCTLContent})
		//list.Counts++
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	list.Totality = GetDCTLCounts(this.TableName(), para)
	return nil, &list
}
func GetDCTLCounts(tablename string, para *TblDCTLSearchPara) int64 {
	qslice, whereflag := modelsPublic.DefaultParaCmd("getcounts", tablename, &para.PField)
	if tablename != "" {
		if para.Status != "" {
			if whereflag != 0 {
				temp_s := fmt.Sprintf(" AND status IN ('%s')", para.Status)
				qslice = append(qslice, temp_s)
			} else {
				temp_s := fmt.Sprintf(" WHERE status IN ('%s')", para.Status)
				qslice = append(qslice, temp_s)
			}
		}
	}
	qslice = append(qslice, ";")
	query := strings.Join(qslice, "")
	mlog.Debug(string(query))
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
	return int64(count)
}
