/********获取文件威胁数(天)********/
package vds

import (
	"apt-web-server_v2/models/db"
	"apt-web-server_v2/models/modelsPublic"
	"fmt"
	"strings"
)

func (this *TblFTT) TableName() string {
	return "vds_count_day"
}

func GetFTTMysqlCmd(tablename string, para *TblFTTSearchPara) []string {
	qslice, _ := modelsPublic.DefaultParaCmd("getlist", tablename, &para.PField)
	qslice = append(qslice, ";")
	fmt.Println(tablename, para)
	return qslice
}

func (this *TblFTT) GetFTTrent(para *TblFTTSearchPara) (error, *TblFTTData) {
	qslice := GetFTTMysqlCmd(this.TableName(), para)
	query := strings.Join(qslice, "")
	rows, err := db.DB.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	list := TblFTTData{}
	ugcCount := new(TblFTTList)
	//ugcCount := new(TblFTT)
	for rows.Next() {
		ugc := new(TblFTT)
		err = rows.Scan(
			&ugc.Id,
			&ugc.Time,
			&ugc.BackDoor,
			&ugc.Trojan,
			&ugc.RiskTool,
			&ugc.Spyware,
			&ugc.Malware,
			&ugc.Virus,
			&ugc.Worm,
			&ugc.Joke,
			&ugc.Adware,
			&ugc.HackTool,
			&ugc.Exploit,
			&ugc.Other)
		if err != nil {
			return err, nil
		}
		//ugcCount.Id = ugc.Id
		//ugcCount.Time = ugc.Time
		ugcCount.BackDoor += ugc.BackDoor
		ugcCount.Trojan += ugc.Trojan
		//ugcCount.RiskTool += ugc.RiskTool
		ugcCount.Spyware += ugc.Spyware
		ugcCount.Malware += (ugc.Malware + ugc.RiskTool + ugc.Joke + ugc.Adware + ugc.Other)
		ugcCount.Virus += ugc.Virus
		ugcCount.Worm += ugc.Worm
		//ugcCount.Joke += ugc.Joke
		ugcCount.HackTool += ugc.HackTool
		ugcCount.Exploit += ugc.Exploit
		//ugcCount.Other += ugc.Other
		//fmt.Println(ugcCount)
	}
	list.Elements = append(list.Elements, TblFTTList{ugcCount.TblFTCRequest})
	if err := rows.Err(); err != nil {
		return err, nil
	}
	list.Counts = int64(ugcCount.BackDoor + ugcCount.Trojan + ugcCount.Spyware +
		ugcCount.Malware + ugcCount.Virus + ugcCount.Worm +
		ugcCount.HackTool + ugcCount.Exploit)
	//list.Counts = GetFTTCounts("", para)
	//list.Totality = GetFTTCounts(this.TableName(), para)

	return nil, &list
}

func GetFTTCounts(tablename string, para *TblFTTSearchPara) int64 {
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
