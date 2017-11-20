/********获取恶意流量详情********/
package ids

import (
	"apt-web-server_v2/models/db"
	"apt-web-server_v2/models/modelsPublic"
	"apt-web-server_v2/modules/mlog"
	"fmt"
	"strings"
)

func (this *TblMFD) TableName() string {
	return "alert_ids"
}

func GetMFDMysqlCmd(tablename string, para *TblMFDSearchPara) []string {
	//qslice, whereflag := modelsPublic.DefaultParaCmd("getlist", tablename, &para.PField)\
	qslice := make([]string, 0)
	whereflag := 0
	cmd := fmt.Sprintf(`SELECT`)
	qslice = append(qslice, cmd)
	column := fmt.Sprintf(` src_ip,src_port,dest_ip,dest_port,src_country,
	    src_province,src_city,src_latitude,src_longitude,dest_country,
		dest_province,dest_city,dest_latitude,dest_longitude,operators,time,
		byzoro_type,attack_type,proto,details,severity,engine`)
	qslice = append(qslice, column)
	fromTbl := fmt.Sprintf(` FROM %s`, tablename)
	qslice = append(qslice, fromTbl)
	if para.PField.Start != 0 && para.PField.End != 0 {
		var temp_se string
		temp_se = fmt.Sprintf(" WHERE time BETWEEN %d AND %d", para.PField.Start, para.PField.End)
		whereflag = 1
		qslice = append(qslice, temp_se)
	}
	if para.Id != 0 {
		if whereflag == 1 {
			temp_id := fmt.Sprintf(" AND id IN (%d)", para.Id)
			qslice = append(qslice, temp_id)
		} else {
			temp_id := fmt.Sprintf(" WHERE id IN (%d)", para.Id)
			qslice = append(qslice, temp_id)
			whereflag = 1
		}
	}
	if para.Type != "" {
		if whereflag == 1 {
			temp_t := fmt.Sprintf(" AND attack_type IN ('%s') OR byzoro_type IN ('%s')", para.Type, para.Type)
			qslice = append(qslice, temp_t)
		} else {
			temp_t := fmt.Sprintf(" WHERE attack_type IN ('%s') OR byzoro_type IN ('%s')", para.Type, para.Type)
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

func (this *TblMFD) GetMFDetails(para *TblMFDSearchPara) (error, *TblMFDData) {
	qslice := GetMFDMysqlCmd(this.TableName(), para)
	query := strings.Join(qslice, "")
	mlog.Debug(string(query))
	rows, err := db.DB.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	list := TblMFDData{}
	for rows.Next() {
		ugc := new(TblMFD)
		err = rows.Scan(
			&ugc.SrcIp,
			&ugc.SrcPort,
			&ugc.DestIp,
			&ugc.DestPort,
			&ugc.SrcIpInfo.Country,
			&ugc.SrcIpInfo.Province,
			&ugc.SrcIpInfo.City,
			&ugc.SrcIpInfo.Lat,
			&ugc.SrcIpInfo.Lng,
			&ugc.DestIpInfo.Country,
			&ugc.DestIpInfo.Province,
			&ugc.DestIpInfo.City,
			&ugc.DestIpInfo.Lat,
			&ugc.DestIpInfo.Lng,
			&ugc.Operators,
			&ugc.Time,
			&ugc.ByzoroType,
			&ugc.AttackType,
			&ugc.Proto,
			&ugc.Details,
			&ugc.Severity,
			&ugc.Engine)
		if err != nil {
			return err, nil
		}
		list.Elements = append(list.Elements, TblMFDList{ugc.TblMFDContent})
		list.Counts++
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	//list.Counts = GetMFDCounts("", para)
	list.Totality = GetMFDCounts(this.TableName(), para)
	return nil, &list
}

func GetMFDCounts(tablename string, para *TblMFDSearchPara) int64 {
	qslice, whereflag := modelsPublic.DefaultParaCmd("getcounts", tablename, &para.PField)
	if tablename != "" {
		if para.Type != "" {
			if whereflag != 0 {
				temp_t := fmt.Sprintf(" AND attack_type IN ('%s') OR byzoro_type IN ('%s')", para.Type, para.Type)
				qslice = append(qslice, temp_t)
			} else {
				temp_t := fmt.Sprintf(" WHERE attack_type IN ('%s') OR byzoro_type IN ('%s')", para.Type, para.Type)
				qslice = append(qslice, temp_t)
			}
		}
		if para.Id != 0 {
			if whereflag == 1 {
				temp_id := fmt.Sprintf(" AND id IN (%d)", para.Id)
				qslice = append(qslice, temp_id)
			} else {
				temp_id := fmt.Sprintf(" WHERE id IN (%d)", para.Id)
				qslice = append(qslice, temp_id)
				whereflag = 1
			}
		}
	}
	qslice = append(qslice, ";")
	query := strings.Join(qslice, "")
	mlog.Debug(string(query))
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
