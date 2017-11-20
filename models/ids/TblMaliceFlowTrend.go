/********获取恶意流量数(天)********/
package ids

import (
	"apt-web-server_v2/models/db"
	"apt-web-server_v2/models/modelsPublic"
	"fmt"
	"strings"
)

func (this *TblMFT) TableName() string {
	return "byzoro_ids_count"
	//return "ids_count_day"
}
func GetMFTMysqlCmd(tablename string, para *TblMFTSearchPara) []string {
	qslice, _ := modelsPublic.DefaultParaCmd("getlist", tablename, &para.PField)
	qslice = append(qslice, ";")
	fmt.Println(tablename, para)
	return qslice
}

func (this *TblMFT) GetMFTrend(para *TblMFTSearchPara) (error, *TblMFTData) {
	qslice := GetMFTMysqlCmd(this.TableName(), para)
	query := strings.Join(qslice, "")
	rows, err := db.DB.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	list := TblMFTData{}
	ugcCount := new(TblMFT)
	for rows.Next() {
		ugc := new(TblMFT)
		err = rows.Scan(
			&ugc.Id,
			&ugc.Time,
			&ugc.PrivilegeGain,
			&ugc.DDos,
			&ugc.InformationLeak,
			&ugc.WebAttack,
			&ugc.ApplicationAttack,
			&ugc.CandC,
			&ugc.Malware,
			&ugc.MiscAttack,
			&ugc.Other)
		if err != nil {
			return err, nil
		}
		ugcCount.Id = ugc.Id
		ugcCount.Time = ugc.Time
		ugcCount.PrivilegeGain += ugc.PrivilegeGain
		ugcCount.DDos += ugc.DDos
		ugcCount.InformationLeak += ugc.InformationLeak
		ugcCount.WebAttack += ugc.WebAttack
		ugcCount.ApplicationAttack += ugc.ApplicationAttack
		ugcCount.CandC += ugc.CandC
		ugcCount.Malware += ugc.Malware
		ugcCount.MiscAttack += ugc.MiscAttack
	}
	list.Elements = append(list.Elements, TblMFTList{ugcCount.TblMFL})
	if err := rows.Err(); err != nil {
		return err, nil
	}
	list.Counts = int64(ugcCount.PrivilegeGain + ugcCount.DDos +
		ugcCount.InformationLeak + ugcCount.WebAttack +
		ugcCount.ApplicationAttack + ugcCount.CandC +
		ugcCount.Malware + ugcCount.MiscAttack)
	//list.Counts = GetMFTCounts("", para)
	return nil, &list
}

func GetMFTCounts(tablename string, para *TblMFTSearchPara) int64 {
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
