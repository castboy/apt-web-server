/********获取恶意流量详情********/
package models

import (
	"apt-web-server/modules/mlog"
	"fmt"
	"strings"
)

func (this *TblMFD) TableName() string {
	return "alert_ids"
}

func GetMFDMysqlCmd(tablename string, para *TblMFDSearchPara) []string {
	qslice, whereflag := DefaultParaCmd("getlist", tablename, &para.PField)
	if para.Type != "" {
		if whereflag == 1 {
			temp_t := fmt.Sprintf(" and attack_type='%s' or byzoro_type='%s'", para.Type, para.Type)
			qslice = append(qslice, temp_t)
		} else {
			temp_t := fmt.Sprintf(" where attack_type='%s' or byzoro_type='%s'", para.Type, para.Type)
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
	return qslice
}

func (this *TblMFD) GetMFDetails(para *TblMFDSearchPara) (error, *TblMFDData) {
	qslice := GetMFDMysqlCmd(this.TableName(), para)
	query := strings.Join(qslice, "")
	mlog.Debug(string(query))
	rows, err := db.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	list := TblMFDData{}
	for rows.Next() {
		ugc := new(TblMFD)
		err = rows.Scan(
			&ugc.Id,
			&ugc.Time,
			&ugc.SrcIp,
			&ugc.SrcPort,
			&ugc.DestIp,
			&ugc.DestPort,
			&ugc.Proto,
			&ugc.AttackType,
			&ugc.Details,
			&ugc.Severity,
			&ugc.Engine,
			&ugc.ByzoroType)
		if err != nil {
			return err, nil
		}
		list.Elements = append(list.Elements, TblMFDList{ugc.TblMFDContent})
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	list.Counts = GetMFDCounts("", para)
	list.Totality = GetMFDCounts(this.TableName(), para)
	return nil, &list
}

func GetMFDCounts(tablename string, para *TblMFDSearchPara) int64 {
	qslice, whereflag := DefaultParaCmd("getcounts", tablename, &para.PField)
	if tablename != "" {
		if para.Type != "" {
			if whereflag != 0 {
				temp_t := fmt.Sprintf(" and attack_type='%s' or byzoro_type='%s'", para.Type, para.Type)
				qslice = append(qslice, temp_t)
			} else {
				temp_t := fmt.Sprintf(" where attack_type='%s' or byzoro_type='%s'", para.Type, para.Type)
				qslice = append(qslice, temp_t)
			}

		}
	}
	qslice = append(qslice, ";")
	query := strings.Join(qslice, "")
	mlog.Debug(string(query))
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

func (this *TblMFD) CreateSql() string {
	return fmt.Sprintf(
		`CREATE TABLE %s (
		id   integer unsigned  AUTO_INCREMENT NOT NULL,
		time   BIGINT NOT NULL DEFAULT 0,
		src_ip varchar(20) NOT NULL DEFAULT '',
		src_port integer NOT NULL ,
		dest_ip varchar(20) NOT NULL DEFAULT '',
		dest_port integer NOT NULL ,
		proto integer NOT NULL ,
		attack_type varchar(20) NOT NULL DEFAULT '',
		details text NOT NULL ,
		severity integer NOT NULL ,
		engine varchar(20) NOT NULL DEFAULT '',
		PRIMARY KEY (id),
		FULLTEXT(details)
		)ENGINE=MyISAM DEFAULT CHARSET=utf8;`,
		this.TableName())
}

