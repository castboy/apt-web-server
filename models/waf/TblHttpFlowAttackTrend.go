/********获取http流量攻击分类数(天)********/
package waf

import (
	"apt-web-server_v2/models/db"
	"apt-web-server_v2/models/modelsPublic"
	"fmt"
	"strings"
)

func (this *TblHFAT) TableName() string {
	return "waf_count_day"
}

func GetHFATMysqlCmd(tablename string, para *TblHFATSearchPara) []string {
	qslice, _ := modelsPublic.DefaultParaCmd("getlist", tablename, &para.PField)
	qslice = append(qslice, ";")
	fmt.Println(tablename, para)
	return qslice
}

func (this *TblHFAT) GetHFATrend(para *TblHFATSearchPara) (error, *TblHFATData) {
	qslice := GetHFATMysqlCmd(this.TableName(), para)
	query := strings.Join(qslice, "")
	rows, err := db.DB.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	list := TblHFATData{}
	ugcCount := new(TblHFACList)
	for rows.Next() {
		ugc := new(TblHFAT)
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
		ugcCount.AttDisclosure += ugc.AttDisclosure
		ugcCount.AttDdos += ugc.AttDdos
		ugcCount.AttReputationIp += ugc.AttReputationIp
		ugcCount.AttLfi += ugc.AttLfi
		ugcCount.AttSqli += ugc.AttSqli
		ugcCount.AttXSS += ugc.AttXSS
		ugcCount.AttInjectionPHP += ugc.AttInjectionPHP
		ugcCount.AttGeneric += ugc.AttGeneric
		ugcCount.AttRce += ugc.AttRce
		ugcCount.AttProtocol += ugc.AttProtocol
		ugcCount.AttRfi += ugc.AttRfi
		ugcCount.AttFixation += ugc.AttFixation
		ugcCount.Scaning += ugc.Scaning
		ugcCount.Other += ugc.Other
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	list.Elements = append(list.Elements, TblHFATList{ugcCount.TblHFACRequest})
	list.Counts = int64(ugcCount.AttDisclosure + ugcCount.AttDdos +
		ugcCount.AttReputationIp + ugcCount.AttLfi + ugcCount.AttSqli +
		ugcCount.AttXSS + ugcCount.AttInjectionPHP + ugcCount.AttGeneric +
		ugcCount.AttRce + ugcCount.AttProtocol + ugcCount.AttRfi +
		ugcCount.AttFixation + ugcCount.Scaning + ugcCount.Other)

	return nil, &list
}

func GetHFATCounts(tablename string, para *TblHFATSearchPara) int64 {
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
