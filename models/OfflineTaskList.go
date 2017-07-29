package models

import (
	"apt-web-server/modules/mlog"
	"fmt"
	"strings"
)

func GetTLMsqlCmd(tablename string, para *TblTaskSearchPara) []string {
	qslice, flag := DefaultParaCmd("getlist", tablename, &para.PField)
	if para.Name != "" {
		var tmp_name string
		if flag == 0 {
			tmp_name = fmt.Sprintf(" where name='%s'", para.Name)
			flag = 1
		} else {
			tmp_name = fmt.Sprintf(" and name='%s'", para.Name)
		}
		qslice = append(qslice, tmp_name)
	}
	if para.Time != 0 {
		var tmp_time string
		if flag == 0 {
			tmp_time = fmt.Sprintf(" where time = %d", para.Time)
			flag = 1
		} else {
			tmp_time = fmt.Sprintf(" and time=%d", para.Time)
		}
		qslice = append(qslice, tmp_time)
	}
	if para.Type != "" {
		var tmp_type string
		if flag == 0 {
			tmp_type = fmt.Sprintf(" where type='%s'", para.Type)
			flag = 1
		} else {
			tmp_type = fmt.Sprintf(" and type='%s'", para.Type)
		}
		qslice = append(qslice, tmp_type)
	}
	if para.Status != "" {
		var tmp_status string
		if flag == 0 {
			tmp_status = fmt.Sprintf(" where status = '%s'", para.Status)
			flag = 1
		} else {
			tmp_status = fmt.Sprintf(" and status = '%s'", para.Status)
		}
		qslice = append(qslice, tmp_status)
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
	fmt.Println(para)
	return qslice
}

func (this *TblOLA) GetTaskDetails(para *TblTaskSearchPara) (error, *TblTaskData) {
	qslice := GetTLMsqlCmd(this.TableName(para.OfflineTag), para)
	query := strings.Join(qslice, "")
	mlog.Debug(string(query))
	rows, err := db.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()

	list := TblTaskData{}
	for rows.Next() {
		ugc := new(TblTaskDetails)
		err = rows.Scan(
			&ugc.Id,
			&ugc.Name,
			&ugc.CreateTime,
			&ugc.Type,
			&ugc.DataStart,
			&ugc.DataEnd,
			&ugc.Weight,
			&ugc.Topic,
			&ugc.Status,
			&ugc.Details)
		if err != nil {
			return err, nil
		}
		/*
			if ugc.Status == "wait" || ugc.Status == "running" {
				err, _ := this.GetStatus(ugc.Name, ugc.CreateTime)
				if err != nil {
					fmt.Println("get task ", ugc.Id, " status error!")
				}
			}
		*/
		list.Elements = append(list.Elements, TblTaskList{ugc.TaskContent})
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	list.Counts = GetTLCounts("", para)
	list.Totality = GetTLCounts(this.TableName(para.OfflineTag), para)
	return nil, &list
}
func GetTLCounts(tablename string, para *TblTaskSearchPara) int64 {
	qslice, whereflag := DefaultParaCmd("getcounts", tablename, &para.PField)
	if tablename != "" {
		if para.Type != "" {
			if whereflag != 0 {
				temp_t := fmt.Sprintf(" and type='%s'", para.Type)
				qslice = append(qslice, temp_t)
			} else {
				temp_t := fmt.Sprintf(" where type='%s'", para.Type)
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
