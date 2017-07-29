/********获取暴力破解详情********/
package models

import (
	"fmt"
	"strings"
	"time"
)

func (this *TblBFD) TableName() string {
	return "bruteForceAction"
}

func GetBFDMysqlCmd(tablename string, para *TblBFDSearchPara) []string {
	//qslice, whereflag := DefaultParaCmd("getlist", tablename, &para.PField)
	var whereflag int
	qslice := make([]string, 0)
	qslice_tmp := fmt.Sprintf(`select * from %s `, tablename)
	qslice = append(qslice, qslice_tmp)
	if para.PField.Start != 0 && para.PField.End != 0 {
		var qslice_time string
		tmStart := time.Unix(para.PField.Start, 0)
		tmEnd := time.Unix(para.PField.End, 0)
		qslice_time = fmt.Sprintf("where time between '%s' and '%s' ", tmStart.Format("2006-01-02 03:04:05 PM"), tmEnd.Format("2006-01-02 03:04:05 PM"))
		qslice = append(qslice, qslice_time)
		whereflag = 1
	}
	if para.Type != "" {
		if whereflag == 1 {
			temp_t := fmt.Sprintf(" and name='%s'", para.Type)
			qslice = append(qslice, temp_t)
		} else {
			temp_t := fmt.Sprintf(" where name='%s'", para.Type)
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

func (this *TblBFD) GetBruteForceDetails(para *TblBFDSearchPara) (error, *TblBFDData) {
	qslice := GetBFDMysqlCmd(this.TableName(), para)
	query := strings.Join(qslice, "")
	rows, err := db.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	list := TblBFDData{}
	for rows.Next() {
		ugc := new(TblBFD)
		err = rows.Scan(
			//&ugc.Id,
			&ugc.Ip,
			&ugc.Port,
			&ugc.Time,
			&ugc.Count,
			&ugc.Name,
			&ugc.Level)
		if err != nil {
			return err, nil
		}
		list.Elements = append(list.Elements, TblBFDList{ugc.TblBFDContent})
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	list.Counts = GetBruteForceCounts("", para)
	list.Totality = GetBruteForceCounts(this.TableName(), para)

	return nil, &list
}

func GetBruteForceCounts(tablename string, para *TblBFDSearchPara) int64 {
	//qslice, whereflag := DefaultParaCmd("getcounts", tablename, &para.PField)
	var whereflag int
	qslice := make([]string, 0)
	qslice_count := fmt.Sprintf(`select count(*) from %s `, tablename)
	qslice_total := fmt.Sprintf(`select FOUND_ROWS() as count`)
	if tablename != "" {
		qslice = append(qslice, qslice_count)
		if para.PField.Start != 0 && para.PField.End != 0 {
			var qslice_time string
			tmStart := time.Unix(para.PField.Start, 0)
			tmEnd := time.Unix(para.PField.End, 0)
			qslice_time = fmt.Sprintf("where time between '%s' and '%s' ", tmStart.Format("2006-01-02 03:04:05 PM"), tmEnd.Format("2006-01-02 03:04:05 PM"))
			qslice = append(qslice, qslice_time)
			whereflag = 1
		}
		if tablename != "" {
			if para.Type != "" {
				if whereflag != 0 {
					temp_t := fmt.Sprintf(" and name='%s'", para.Type)
					qslice = append(qslice, temp_t)
				} else {
					temp_t := fmt.Sprintf(" where name='%s'", para.Type)
					qslice = append(qslice, temp_t)
				}

			}
		}
	} else {
		qslice = append(qslice, qslice_total)
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

func (this *TblBFD) CreateSql() string {
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
