/********获取紧急事件详情********/
package models

import (
	"apt-web-server/modules/mlog"
	"fmt"
	"strings"
)

func (this *TblUgcD) TableName() string {
	return "urgencymold"
}

func GetUgcDMysqlCmd(tablename string, para *TblUgcDSearchPara) []string {
	qslice, whereflag := DefaultParaCmd("getlist", tablename, &para.PField)
	if para.Type != "" {
		if whereflag == 1 {
			temp_t := fmt.Sprintf(" and attack_type='%s'", para.Type)
			qslice = append(qslice, temp_t)
		} else {
			temp_t := fmt.Sprintf(" where attack_type='%s'", para.Type)
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
	} else if para.Count != 0 {
		para.Page = para.Page - 1
		temp_PC := fmt.Sprintf(" limit %d,%d", para.Page*para.Count, para.Count)
		qslice = append(qslice, temp_PC)
	}
	qslice = append(qslice, ";")
	fmt.Println(tablename, para)
	fmt.Println(qslice)
	return qslice
}

func (this *TblUgcD) GetUrgencyDetails(para *TblUgcDSearchPara) (error, *TblUgcDData) {
	qslice := GetUgcDMysqlCmd(this.TableName(), para)
	query := strings.Join(qslice, "")

	rows, err := db.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	list := TblUgcDData{}
	for rows.Next() {
		ugc := new(TblUgcD)
		err = rows.Scan(
			&ugc.Id,
			&ugc.Time,
			&ugc.SrcIp,
			&ugc.SrcPort,
			&ugc.DestIp,
			&ugc.DestPort,
			&ugc.Proto,
			&ugc.ServerName,
			&ugc.AttackType,
			&ugc.Serverity,
			&ugc.AttackerOS,
			&ugc.AttackedOS,
			&ugc.Details)
		if err != nil {
			mlog.Debug(query, err)
			//return err, nil
		}
		list.Elements = append(list.Elements, TblUgcDList{ugc.TblUgcDContent})
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	list.Counts = GetUrgencyDCounts("", para)
	list.Totality = GetUrgencyDCounts(this.TableName(), para)

	return nil, &list
}

func GetUrgencyDCounts(tablename string, para *TblUgcDSearchPara) int64 {
	qslice, whereflag := DefaultParaCmd("getcounts", tablename, &para.PField)
	if tablename != "" {
		if para.Type != "" {
			if whereflag != 0 {
				temp_t := fmt.Sprintf(" and attack_type='%s'", para.Type)
				qslice = append(qslice, temp_t)
			} else {
				temp_t := fmt.Sprintf(" where attack_type='%s'", para.Type)
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

func (this *TblUgcD) CreateSql() string {
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
