/********获取扫描事件详情********/
package models

import (
	"fmt"
	"strings"
)

func (this *TblSED) TableName() string {
	return "alert_portscan"
}

func GetSEDMysqlCmd(tablename string, para *TblSEDSearchPara) []string {
	qslice, whereflag := DefaultParaCmd("getlist", tablename, &para.PField)
	if para.Type != "" {
		if whereflag == 1 {
			temp_t := fmt.Sprintf(" and conntype in ('%s')", para.Type)
			qslice = append(qslice, temp_t)
		} else {
			temp_t := fmt.Sprintf(" where conntype in ('%s')", para.Type)
			qslice = append(qslice, temp_t)
			whereflag = 1
		}
	}
	if para.Sort != "" {
		temp_s := fmt.Sprintf(" order by %s %s", para.Sort, para.Order)
		qslice = append(qslice, temp_s)
	}
	if para.LastCount != 0 {
		temp_LC := fmt.Sprintf(" limit %d", para.LastCount)
		qslice = append(qslice, temp_LC)
	} else {
		para.Page = para.Page - 1
		temp_PC := fmt.Sprintf(" limit %d,%d", para.Page*para.Count, para.Count)
		qslice = append(qslice, temp_PC)
	}
	qslice = append(qslice, ";")
	fmt.Println(tablename, para)
	fmt.Println(qslice)
	return qslice
}

func (this *TblSED) GetScanEventDetails(para *TblSEDSearchPara) (error, *TblSEDData) {
	qslice := GetSEDMysqlCmd(this.TableName(), para)
	query := strings.Join(qslice, "")
	rows, err := db.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	list := TblSEDData{}
	for rows.Next() {
		ugc := new(TblSED)
		err = rows.Scan(
			&ugc.Id,
			&ugc.Time,
			&ugc.Conntype,
			&ugc.Host,
			&ugc.AlertType,
			&ugc.Count)
		if err != nil {
			return err, nil
		}
		list.Elements = append(list.Elements, TblSEDList{ugc.TblSEDContent})
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	list.Counts = GetScanEventCounts("", para)
	list.Totality = GetScanEventCounts(this.TableName(), para)

	return nil, &list
}

func GetScanEventCounts(tablename string, para *TblSEDSearchPara) int64 {
	qslice, whereflag := DefaultParaCmd("getcounts", tablename, &para.PField)
	if tablename != "" {
		if para.Type != "" {
			if whereflag != 0 {
				temp_t := fmt.Sprintf(" and conntype in ('%s')", para.Type)
				qslice = append(qslice, temp_t)
			} else {
				temp_t := fmt.Sprintf(" where conntype in ('%s')", para.Type)
				qslice = append(qslice, temp_t)
			}

		}
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

func (this *TblSED) CreateSql() string {
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
