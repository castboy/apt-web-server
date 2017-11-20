/********获取HTTP流量攻击详情********/
package waf

import (
	"apt-web-server_v2/models/db"
	"apt-web-server_v2/models/modelsPublic"
	"apt-web-server_v2/modules/mlog"
	"encoding/base64"
	"fmt"
	"strings"
)

func (this *TblHFAD) TableName(tage string) string {
	switch tage {
	case "offline":
		return "alert_waf_offline"
	case "duplicate":
		return "waf_dup"
	case "offlinedup":
		return "waf_ofl_dup"
	case "rule":
		return "alert_waf_offline_rule"
	default:
		return "alert_waf"
	}
}

func GetHFADMysqlCmd(tablename string, para *TblHFADSearchPara) []string {
	//qslice, whereflag := modelsPublic.DefaultParaCmd("getlist", tablename, &para.PField)
	qslice := make([]string, 0)
	whereflag := 0
	qslice = append(qslice, `SELECT`)
	column := fmt.Sprintf(` id,src_ip,src_port,dest_ip,dest_port,proto,src_country,
	    src_province,src_city,src_latitude,src_longitude,dest_country,
		dest_province,dest_city,dest_latitude,dest_longitude,operators,attack,
		time,client,rev,severity,maturity,accuracy,hostname,unique_id,ref,tags,
		rule_file,rule_line,rule_id,rule_data,rule_ver,version,request,response,msg,uri`)
	qslice = append(qslice, column)
	if tablename == "waf_dup" || tablename == "waf_ofl_dup" {
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
	var taskTable string
	switch para.Tage {
	case "offline", "offlinedup", "duplicate":
		taskTable = "ofl_task"
	case "rule":
		taskTable = "ofl_rule"
	}
	if para.Tage == "offline" ||
		para.Tage == "offlinedup" ||
		para.Tage == "duplicate" ||
		para.Tage == "rule" {
		err, taskId := modelsPublic.GetOfflineTaskID(taskTable, para.TaskName, para.CreateTime)
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
		if para.Type == "scaningprobe" {
			if whereflag == 1 {
				temp_t := fmt.Sprintf(` AND attack IN ('reputation_scanner',
				    'reputation_scripting','reputation_crawler')`)
				qslice = append(qslice, temp_t)
			} else {
				temp_t := fmt.Sprintf(` WHERE attack IN ('reputation_scanner',
				    'reputation_scripting','reputation_crawler')`)
				qslice = append(qslice, temp_t)
				whereflag = 1
			}
		} else {
			if whereflag == 1 {
				temp_t := fmt.Sprintf(" AND attack IN ('%s')", para.Type)
				qslice = append(qslice, temp_t)
			} else {
				temp_t := fmt.Sprintf(" WHERE attack IN ('%s')", para.Type)
				qslice = append(qslice, temp_t)
				whereflag = 1
			}
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

func (this *TblHFAD) GetHFADetails(para *TblHFADSearchPara) (error, *TblHFADData) {
	qslice := GetHFADMysqlCmd(this.TableName(para.Tage), para)
	query := strings.Join(qslice, "")
	mlog.Debug(string(query))
	rows, err := db.DB.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	list := TblHFADData{}
	for rows.Next() {
		ugc := new(TblHFAD)
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
			&ugc.Attack,
			&ugc.Time,
			&ugc.Client,
			&ugc.Rev,
			&ugc.Severity,
			&ugc.Maturity,
			&ugc.Accuracy,
			&ugc.HostName,
			&ugc.UniqueId,
			&ugc.Ref,
			&ugc.Tags,
			&ugc.Rule.File,
			&ugc.Rule.Line,
			&ugc.Rule.Id,
			&ugc.Rule.Data,
			&ugc.Rule.Version,
			&ugc.Version,
			&ugc.Request,
			&ugc.Response,
			&ugc.Message,
			&ugc.Uri)
		if err != nil {
			//return err, nil
		}
		uri, _ := base64.StdEncoding.DecodeString(ugc.Uri)
		ugc.Uri = string(uri)
		ruledata, _ := base64.StdEncoding.DecodeString(ugc.Rule.Data)
		ugc.Rule.Data = string(ruledata)
		list.Elements = append(list.Elements, TblHFADList{ugc.TblHFADContent})
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	list.Counts = GetHFADCounts("", para)
	list.Totality = GetHFADCounts(this.TableName(para.Tage), para)

	return nil, &list
}
func (this *TblHFAD) GetHFADetailsDUP(para *TblHFADSearchPara) (error, *TblHFADData) {
	qslice := GetHFADMysqlCmd(this.TableName(para.Tage), para)
	query := strings.Join(qslice, "")
	mlog.Debug(string(query))
	rows, err := db.DB.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	list := TblHFADData{}
	for rows.Next() {
		ugc := new(TblHFAD)
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
			&ugc.Attack,
			&ugc.Time,
			&ugc.Client,
			&ugc.Rev,
			&ugc.Severity,
			&ugc.Maturity,
			&ugc.Accuracy,
			&ugc.HostName,
			&ugc.UniqueId,
			&ugc.Ref,
			&ugc.Tags,
			&ugc.Rule.File,
			&ugc.Rule.Line,
			&ugc.Rule.Id,
			&ugc.Rule.Data,
			&ugc.Rule.Version,
			&ugc.Version,
			&ugc.Request,
			&ugc.Response,
			&ugc.Message,
			&ugc.Uri,
			&ugc.TblId,
			&ugc.TaskId)
		if err != nil {
			//return err, nil
		}
		if para.Tage == "offlinedup" && para.MergeTag == "merge" {
			mergecmd := fmt.Sprintf(`call checkwafofflinemerge('%s',%d,'%s',%d,'%s',
			'%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s',%d,'%s','%s',%d,
			%d,%d,'%s','%s','%s','%s','%s',%d,%d,'%s','%s','%s','%s','%s','%s','%s',%d,%d);`,
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
				ugc.Attack,
				ugc.Time,
				ugc.Client,
				ugc.Rev,
				ugc.Severity,
				ugc.Maturity,
				ugc.Accuracy,
				ugc.HostName,
				ugc.UniqueId,
				ugc.Ref,
				ugc.Tags,
				ugc.Rule.File,
				ugc.Rule.Line,
				ugc.Rule.Id,
				ugc.Rule.Data,
				ugc.Rule.Version,
				ugc.Version,
				ugc.Request,
				ugc.Response,
				ugc.Message,
				ugc.Uri,
				ugc.TblId,
				ugc.TaskId)
			rows, err := db.DB.Query(mergecmd)
			//fmt.Println(mergecmd)
			if err != nil {
				fmt.Println("merge data error!")
				ugc.Message = "merge data error!"
				//return err, nil
			}
			defer rows.Close()

		}
		uri, _ := base64.StdEncoding.DecodeString(ugc.Uri)
		ugc.Uri = string(uri)
		ruledata, _ := base64.StdEncoding.DecodeString(ugc.Rule.Data)
		ugc.Rule.Data = string(ruledata)
		list.Elements = append(list.Elements, TblHFADList{ugc.TblHFADContent})
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	if para.Tage == "offlinedup" && para.MergeTag == "merge" {
		return nil, &list
	}
	list.Counts = GetHFADCounts("", para)
	list.Totality = GetHFADCounts(this.TableName(para.Tage), para)

	return nil, &list
}

func GetHFADCounts(tablename string, para *TblHFADSearchPara) int64 {
	qslice, whereflag := modelsPublic.DefaultParaCmd("getcounts", tablename, &para.PField)
	var taskTbl string
	if tablename != "" {
		fmt.Println("para.Tage=", para.Tage)
		if para.Tage != "" && para.Tage != "online" {
			switch para.Tage {
			case "offline", "offlinedup", "duplicate":
				taskTbl = "ofl_task"
			case "rule":
				taskTbl = "ofl_rule"
			}
			err, taskId := modelsPublic.GetOfflineTaskID(taskTbl, para.TaskName, para.CreateTime)
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
			if para.Type == "scaningprobe" {
				if whereflag == 1 {
					temp_t := fmt.Sprintf(` AND attack IN ('reputation_scanner',
					    'reputation_scripting','reputation_crawler')`)
					qslice = append(qslice, temp_t)
				} else {
					temp_t := fmt.Sprintf(` WHERE attack IN ('reputation_scanner',
					    'reputation_scripting','reputation_crawler')`)
					qslice = append(qslice, temp_t)
					whereflag = 1
				}
			} else {
				if whereflag != 0 {
					temp_t := fmt.Sprintf(" AND attack IN ('%s')", para.Type)
					qslice = append(qslice, temp_t)
				} else {
					temp_t := fmt.Sprintf(" WHERE attack IN ('%s')", para.Type)
					qslice = append(qslice, temp_t)
				}
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
