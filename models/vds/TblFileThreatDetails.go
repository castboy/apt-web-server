/********获取文件威胁详情********/
package vds

import (
	"apt-web-server_v2/models/db"
	"apt-web-server_v2/models/modelsPublic"
	"apt-web-server_v2/modules/mlog"
	"fmt"
	"strings"
)

func (this *TblFTD) TableName(tage string) string {
	switch tage {
	case "offline":
		return "alert_vds_offline"
	case "duplicate":
		return "vds_dup"
	case "offlinedup":
		return "vds_ofl_dup"
	default:
		return "alert_vds"
	}
}

func GetFTDMysqlCmd(tablename string, para *TblFTDSearchPara) []string {
	//qslice, whereflag := modelsPublic.DefaultParaCmd("getlist", tablename, &para.PField)
	qslice := make([]string, 0)
	whereflag := 0
	cmd := fmt.Sprintf(`SELECT`)
	qslice = append(qslice, cmd)
	column := fmt.Sprintf(` id,src_ip,src_port,dest_ip,dest_port,proto,src_country,
	    src_province,src_city,src_latitude,src_longitude,dest_country,
		dest_province,dest_city,dest_latitude,dest_longitude,operators,time,
		local_vtype,local_vname,local_logtype,threatname,subfile,local_threatname,
		local_platfrom,local_extent,local_enginetype,local_engineip,app_file,http_url`)
	qslice = append(qslice, column)
	if tablename == "vds_dup" || tablename == "vds_ofl_dup" {
		qslice = append(qslice, ",tbl_id,task_id")
	}
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
	if para.Tage == "offline" ||
		para.Tage == "offlinedup" ||
		para.Tage == "duplicate" {
		err, taskId := modelsPublic.GetOfflineTaskID("ofl_task", para.TaskName, para.CreateTime)
		if err != nil {
			fmt.Println("get taskid error")
		}
		if whereflag == 1 {
			temp_id := fmt.Sprintf(" AND task_id=%d", taskId)
			qslice = append(qslice, temp_id)
		} else {
			temp_id := fmt.Sprintf(" WHERE task_id=%d", taskId)
			qslice = append(qslice, temp_id)
			whereflag = 1
		}
	}
	if para.Type != "" {
		if whereflag == 1 {
			temp_t := fmt.Sprintf(" AND local_vtype IN ('%s')", para.Type)
			qslice = append(qslice, temp_t)
		} else {
			temp_t := fmt.Sprintf(" WHERE local_vtype IN ('%s')", para.Type)
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

func (this *TblFTD) GetFTDetails(para *TblFTDSearchPara) (error, *TblFTDData) {
	qslice := GetFTDMysqlCmd(this.TableName(para.Tage), para)
	query := strings.Join(qslice, "")
	mlog.Debug(string(query))
	rows, err := db.DB.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	list := TblFTDData{}
	for rows.Next() {
		ugc := new(TblFTD)
		err = rows.Scan(
			&ugc.Id,
			&ugc.SrcIp,
			&ugc.SrcPort,
			&ugc.DestIp,
			&ugc.DestPort,
			&ugc.Proto,
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
			&ugc.LocalVType,
			&ugc.LocalVName,
			&ugc.LocalLogType,
			&ugc.ThreatName,
			&ugc.SubFile,
			&ugc.LocalThreatName,
			&ugc.LocalPlatfrom,
			&ugc.LocalExtent,
			&ugc.LocalEngineType,
			&ugc.LocalEngineIP,
			&ugc.AppFile,
			&ugc.HttpUrl)
		if err != nil {
			return err, nil
		}
		list.Elements = append(list.Elements, TblFTDList{ugc.TblFTDContent})
		list.Counts++
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	//list.Counts = GetFTDCounts("", para)
	list.Totality = GetFTDCounts(this.TableName(para.Tage), para)

	return nil, &list
}
func (this *TblFTD) GetFTDetailsDUP(para *TblFTDSearchPara) (error, *TblFTDData) {
	qslice := GetFTDMysqlCmd(this.TableName(para.Tage), para)
	query := strings.Join(qslice, "")
	mlog.Debug(string(query))
	rows, err := db.DB.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	list := TblFTDData{}
	for rows.Next() {
		ugc := new(TblFTD)
		err = rows.Scan(
			&ugc.Id,
			&ugc.SrcIp,
			&ugc.SrcPort,
			&ugc.DestIp,
			&ugc.DestPort,
			&ugc.Proto,
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
			&ugc.LocalVType,
			&ugc.LocalVName,
			&ugc.LocalLogType,
			&ugc.ThreatName,
			&ugc.SubFile,
			&ugc.LocalThreatName,
			&ugc.LocalPlatfrom,
			&ugc.LocalExtent,
			&ugc.LocalEngineType,
			&ugc.LocalEngineIP,
			&ugc.AppFile,
			&ugc.HttpUrl,
			&ugc.TblId,
			&ugc.TaskId)
		if err != nil {
			return err, nil
		}
		if para.Tage == "offlinedup" && para.MergeTag == "merge" {
			mergecmd := fmt.Sprintf(`call checkvdsofflinemerge('%s',%d,'%s',%d,'%s',
			'%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s',%d,'%s','%s','%s',
			'%s','%s','%s','%s','%s','%s','%s','%s','%s',%d,%d);`,
				ugc.SrcIp,
				ugc.SrcPort,
				ugc.DestIp,
				ugc.DestPort,
				ugc.Proto,
				ugc.SrcIpInfo.Country,
				ugc.SrcIpInfo.Province,
				ugc.SrcIpInfo.City,
				ugc.SrcIpInfo.Lat,
				ugc.SrcIpInfo.Lng,
				ugc.DestIpInfo.Country,
				ugc.DestIpInfo.Province,
				ugc.DestIpInfo.City,
				ugc.DestIpInfo.Lat,
				ugc.DestIpInfo.Lng,
				ugc.Operators,
				ugc.Time,
				ugc.LocalVType,
				ugc.LocalVName,
				ugc.LocalLogType,
				ugc.ThreatName,
				ugc.SubFile,
				ugc.LocalThreatName,
				ugc.LocalPlatfrom,
				ugc.LocalExtent,
				ugc.LocalEngineType,
				ugc.LocalEngineIP,
				ugc.AppFile,
				ugc.HttpUrl,
				ugc.TblId,
				ugc.TaskId)
			rows, err := db.DB.Query(mergecmd)
			fmt.Println(mergecmd)
			if err != nil {
				fmt.Println("insert merge data error!")
				return err, nil
			}
			defer rows.Close()
		}
		list.Elements = append(list.Elements, TblFTDList{ugc.TblFTDContent})
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	if para.Tage == "offlinedup" && para.MergeTag == "merge" {
		return nil, &list
	}
	list.Counts = GetFTDCounts("", para)
	list.Totality = GetFTDCounts(this.TableName(para.Tage), para)

	return nil, &list
}

func GetFTDCounts(tablename string, para *TblFTDSearchPara) int64 {
	qslice, whereflag := modelsPublic.DefaultParaCmd("getcounts", tablename, &para.PField)
	if tablename != "" {
		if para.Tage != "" && para.Tage != "online" {
			if para.Tage == "offline" ||
				para.Tage == "offlinedup" ||
				para.Tage == "duplicate" {
				err, taskId := modelsPublic.GetOfflineTaskID("ofl_task", para.TaskName, para.CreateTime)
				if err != nil {
					fmt.Println("get taskid error")
				}
				if whereflag == 1 {
					temp_id := fmt.Sprintf(" AND task_id=%d", tablename, taskId)
					qslice = append(qslice, temp_id)
				} else {
					temp_id := fmt.Sprintf(" WHERE task_id=%d", tablename, taskId)
					qslice = append(qslice, temp_id)
					whereflag = 1
				}
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
		if para.Type != "" {
			if whereflag != 0 {
				temp_t := fmt.Sprintf(" AND local_vtype IN ('%s')", para.Type)
				qslice = append(qslice, temp_t)
			} else {
				temp_t := fmt.Sprintf(" WHERE local_vtype IN ('%s')", para.Type)
				qslice = append(qslice, temp_t)
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
