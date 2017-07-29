/********获取攻击数********/
package models

import (
	"fmt"
	"strings"
)

func (this *TblAttackCount) TableName() string {
	return "urgencymold"
}

func GetAttackCountMysqlCmd(datemold string, para *TblAttackCountSearchPara) []string {
	qslice := make([]string, 0)
	qslice_count := fmt.Sprintf(`call allattackcountandip('%s',%d,%d);`, datemold, para.PField.Start, para.PField.End)
	qslice = append(qslice, qslice_count)
	//qslice = append(qslice, ";")
	return qslice
}

func (this *TblAttackCount) GetAttackCount(para *TblAttackCountSearchPara) (error, *TblAttackCountData) {
	var datemold string
	var attackCounts, ipCounts int64
	switch para.Unit {
	case "day":
		datemold = "%Y-%m-%d"
	case "month":
		datemold = "%Y-%m"
	case "hour":
		datemold = "%Y-%m-%d %H"
	case "minute":
		datemold = "%Y-%m-%d %H-%i"
	default:
		datemold = "%Y-%m-%d"
	}
	list := TblAttackCountData{}
	qslice := GetAttackCountMysqlCmd(datemold, para)
	query := strings.Join(qslice, "")
	fmt.Println(qslice)

	rows, err := db.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	ugc := new(TblAttackCount)
	for rows.Next() {
		err = rows.Scan(
			&ugc.Time,
			&ugc.AttackCount,
			&ugc.IpCount)
		if err != nil {
			return err, nil
		}
		attackCounts += ugc.AttackCount
		ipCounts += ugc.IpCount
		list.Elements = append(list.Elements, TblAttackCountList{ugc.TblAttackCountContent})
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}

	list.Counts = attackCounts
	list.Totality = ipCounts
	return nil, &list
}
