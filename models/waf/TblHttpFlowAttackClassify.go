/********获取http流量攻击分类数(天)********/
package waf

import (
	"apt-web-server_v2/models/db"
	"apt-web-server_v2/models/modelsPublic"
	"fmt"
	"strings"
)

func (this *TblHFAC) TableName() string {
	return "waf_count_day"
}

func GetHFACMysqlCmd(tablename string, para *TblHFACSearchPara) []string {
	qslice, _ := modelsPublic.DefaultParaCmd("getlist", tablename, &para.PField)
	qslice = append(qslice, ";")
	fmt.Println(tablename, para)
	return qslice
}

func (this *TblHFAC) GetHFAClassify(para *TblHFACSearchPara) (error, *TblHFACData) {
	qslice := GetHFACMysqlCmd(this.TableName(), para)
	query := strings.Join(qslice, "")
	rows, err := db.DB.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	list := TblHFACData{}
	for rows.Next() {
		ugcCount := new(TblHFACList)
		ugc := new(TblHFAC)
		err = rows.Scan(
			&ugc.Id,
			&ugc.Time,
			&ugc.AttDisclosure,
			&ugc.AttDdos,
			&ugc.AttReputationIp,
			&ugc.AttLfi,
			&ugc.AttSqli,
			&ugc.AttXSS,
			&ugc.AttInjectionPHP,
			&ugc.AttGeneric,
			&ugc.AttRce,
			&ugc.AttProtocol,
			&ugc.AttRfi,
			&ugc.AttFixation,
			&ugc.Scaning,
			&ugc.Other)
		if err != nil {
			return err, nil
		}
		ugcCount.Time = ugc.Time
		ugcCount.HFACUniversal = ugc.HFACUniversal
		//ugcCount.ScanningPprobe = ugc.AttReputScanner + ugc.AttReputSripting + ugc.AttReputCrawler
		list.Classify = append(list.Classify, TblHFACList{ugcCount.HFACContent})
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	list.Counts = GetHFACCounts(this.TableName(), para)

	return nil, &list
}

func GetHFACCounts(tablename string, para *TblHFACSearchPara) int64 {
	qslice, _ := modelsPublic.DefaultParaCmd("getcounts", tablename, &para.PField)
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
	return int64(count)
}

func (this *TblHFAC) CreateSql() string {
	return fmt.Sprintf(
		`CREATE TABLE %s (
		id   integer unsigned  AUTO_INCREMENT NOT NULL,
		time   BIGINT NOT NULL DEFAULT 0,
		s0 BIGINT NOT NULL DEFAULT 0,
		s1 BIGINT NOT NULL DEFAULT 0,
		s2 BIGINT NOT NULL DEFAULT 0,
		s3 BIGINT NOT NULL DEFAULT 0,
		s4 BIGINT NOT NULL DEFAULT 0,
		s5 BIGINT NOT NULL DEFAULT 0,
		s6 BIGINT NOT NULL DEFAULT 0,
		s7 BIGINT NOT NULL DEFAULT 0,
		s8 BIGINT NOT NULL DEFAULT 0,
		s9 BIGINT NOT NULL DEFAULT 0,
		other integer NOT NULL DEFAULT 0,
		PRIMARY KEY (Id)
	)ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
		this.TableName())
}
