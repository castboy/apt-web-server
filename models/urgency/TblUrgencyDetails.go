/********获取紧急事件详情********/
package urgency

import (
	"apt-web-server_v2/models/db"
	"apt-web-server_v2/models/modelsPublic"
	"apt-web-server_v2/models/whiteList"
	"apt-web-server_v2/modules/mlog"
	"fmt"
	"strings"
)

func (this *TblUgcD) TableName() string {
	return "urgencymold"
}

func GetUgcDMysqlCmd(tablename string, para *TblUgcDSearchPara) []string {
	var whereflag int
	qslice := make([]string, 0)
	qslice_tmp := fmt.Sprintf(`SELECT time,src_ip,src_port,dest_ip,dest_port,
	    proto,servername,attack_type,severity,attackeros,attackedos,details 
		FROM %s`, tablename)
	qslice = append(qslice, qslice_tmp)
	if para.PField.Start != 0 && para.PField.End != 0 {
		qslice_time := fmt.Sprintf(" WHERE time BETWEEN '%d' AND '%d'",
			para.PField.Start, para.PField.End)
		qslice = append(qslice, qslice_time)
		whereflag = 1
	}
	if para.Type != "" {
		if whereflag == 1 {
			temp_t := fmt.Sprintf(" AND attack_type IN ('%s')", para.Type)
			qslice = append(qslice, temp_t)
		} else {
			temp_t := fmt.Sprintf(" WHERE attack_type IN ('%s')", para.Type)
			qslice = append(qslice, temp_t)
			whereflag = 1
		}
	}
	if para.Sort != "" {
		temp_s := fmt.Sprintf(" ORDER BY %s %s", para.Sort, para.Order)
		qslice = append(qslice, temp_s)
	}
	if para.LastCount != 0 {
		temp_LC := fmt.Sprintf(" LIMIT %d", para.LastCount)
		qslice = append(qslice, temp_LC)
	} else if para.Count != 0 {
		para.Page = para.Page - 1
		temp_PC := fmt.Sprintf(" LIMIT %d,%d", para.Page*para.Count, para.Count)
		qslice = append(qslice, temp_PC)
	}
	qslice = append(qslice, ";")
	fmt.Println(tablename, para)
	return qslice
}

func (this *TblUgcD) GetUrgencyDetails(para *TblUgcDSearchPara) (error, *TblUgcDData) {
	qslice := GetUgcDMysqlCmd(this.TableName(), para)
	query := strings.Join(qslice, "")

	rows, err := db.DB.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	list := TblUgcDData{}
	for rows.Next() {
		ugc := new(TblUgcD)
		err = rows.Scan(
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
	qslice, whereflag := modelsPublic.DefaultParaCmd("getcounts", tablename, &para.PField)
	if tablename != "" {
		if para.Type != "" {
			if whereflag != 0 {
				temp_t := fmt.Sprintf(" AND attack_type IN ('%s')", para.Type)
				qslice = append(qslice, temp_t)
			} else {
				temp_t := fmt.Sprintf(" WHERE attack_type IN ('%s')", para.Type)
				qslice = append(qslice, temp_t)
			}

		}
	}
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

//check if the abnormal connection alert is in WL
func UgcAbnConnCheckIfInWL(inputList *TblUgcDData) (err error) {
	var query_cmd string
	var pro int32
	for i, _ := range inputList.Elements {
		if inputList.Elements[i].Proto == "TCP" {
			pro = 6
		} else {
			pro = 17
		}
		query_cmd = fmt.Sprintf(`select count(id) from %s where
			src_ip='%s' and src_port='%d' and dest_ip='%s' and dest_port='%d' and proto='%d'`,
			whiteList.WL_wlTableName(), inputList.Elements[i].SrcIp, inputList.Elements[i].SrcPort, inputList.Elements[i].DestIp, inputList.Elements[i].DestPort, pro)
		//fmt.Println("UgcAbnConnCheckIfInWL query_cmd is ", query_cmd)

		rows := modelsPublic.Select_mysql(query_cmd)
		for rows.Next() {
			//fmt.Println("UgcAbnConnCheckIfInWL enter for rows.next() . ")
			var count int
			if rows != nil {
				if err = rows.Scan(&count); err != nil {
					fmt.Println("UgcAbnConnCheckIfInWL err is  ", err.Error())
					rows.Close()
					return err
				}
				//fmt.Println("UgcAbnConnCheckIfInWL query_cmd count is ", count)
				if count != 0 {
					inputList.Elements[i].Details = "Y" // Y means info is already in white list
					rows.Close()
					break
				}
			}
		}
	}

	return err
}
