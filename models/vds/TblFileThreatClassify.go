/********获取文件威胁分类数(天)********/
package vds

import (
	"apt-web-server_v2/models/db"
	"apt-web-server_v2/models/modelsPublic"
	"fmt"
	"strings"
)

func (this *TblFTC) TableName() string {
	return "vds_count_day"
}

func GetFTCMysqlCmd(tablename string, para *TblFTCSearchPara) []string {
	qslice, _ := modelsPublic.DefaultParaCmd("getlist", tablename, &para.PField)
	qslice = append(qslice, ";")
	fmt.Println(tablename, para)
	return qslice
}

func (this *TblFTC) GetUrgencyDetails(para *TblFTCSearchPara) (error, *TblFTCData) {
	qslice := GetFTCMysqlCmd(this.TableName(), para)
	query := strings.Join(qslice, "")
	rows, err := db.DB.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	list := TblFTCData{}
	for rows.Next() {
		ugcCount := new(TblFTCList)
		ugc := new(TblFTC)
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
		ugcCount.Time = ugc.Time
		ugcCount.FTCUniversal = ugc.FTCUniversal
		ugcCount.Malware = ugc.Malware + ugc.RiskTool + ugc.Joke + ugc.Other
		list.Classify = append(list.Classify, TblFTCList{ugcCount.FTCContent})
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	list.Counts = GetFTCCounts(this.TableName(), para)

	return nil, &list
}

func GetFTCCounts(tablename string, para *TblFTCSearchPara) int64 {
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

func (this *TblFTC) CreateSql() string {
	return fmt.Sprintf(
		`CREATE TABLE %s (
		id   integer unsigned  AUTO_INCREMENT NOT NULL,
		time   BIGINT NOT NULL DEFAULT 0,
		backdoor BIGINT NOT NULL DEFAULT 0,
		trojan BIGINT NOT NULL DEFAULT 0,
		risktool BIGINT NOT NULL DEFAULT 0,
		spyware BIGINT NOT NULL DEFAULT 0,
		malware BIGINT NOT NULL DEFAULT 0,
		virus BIGINT NOT NULL DEFAULT 0,
		worm BIGINT NOT NULL DEFAULT 0,
		joke BIGINT NOT NULL DEFAULT 0,
		adware BIGINT NOT NULL DEFAULT 0,
		hacktool BIGINT NOT NULL DEFAULT 0,
		exploit BIGINT NOT NULL DEFAULT 0,
		other BIGINT NOT NULL DEFAULT 0,
		PRIMARY KEY (Id)
	)ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
		this.TableName())
}
