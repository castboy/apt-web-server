package waf

import (
	"apt-web-server_v2/models/db"
	"apt-web-server_v2/models/modelsPublic"
	"fmt"
	"strings"
)

func GetDupSql(tblName, taskName string, time int64) ([]string, string, int, error) {
	var ruleList, sqlCmd string
	qslice := make([]string, 0)
	qslice = append(qslice, `SELECT`)
	err, taskid := modelsPublic.GetOfflineTaskID("ofl_task", taskName, time)
	if err != nil {
		fmt.Println("get taskid error!")
	}
	column := fmt.Sprintf(` id,src_ip,src_port,dest_ip,dest_port,proto,src_country,
	    src_province,src_city,src_latitude,src_longitude,dest_country,
		dest_province,dest_city,dest_latitude,dest_longitude,operators,attack,
		time,client,rev,severity,maturity,accuracy,hostname,unique_id,ref,tags,
		rule_file,rule_line,rule_id,rule_data,rule_ver,version,request,response,msg,uri`)
	switch tblName {
	case "alert_waf":
		sqlCmd = "checkwafduplicate"
		maxTime, minTime := modelsPublic.GetTimes("alert_vds_offline", taskid)
		qslice = append(qslice, column)
		ruleList = fmt.Sprintf(` WHERE time BETWEEN %d AND %d;`, minTime, maxTime)
	case "alert_waf_offline":
		sqlCmd = "checkwafofflinedup"
		qslice = append(qslice, column)
		ruleList = fmt.Sprintf(` WHERE task_id=%d;`, tblName, taskid)
	}
	fromTbl := fmt.Sprintf(` FROM %s`, tblName)
	qslice = append(qslice, fromTbl)
	qslice = append(qslice, ruleList)
	return qslice, sqlCmd, taskid, nil
}
func DuplicateWaf(tblName, taskName string, time int64) error {
	qslice, sqlCmd, taskid, err := GetDupSql(tblName, taskName, time)
	if err != nil {
		return err
	}
	query := strings.Join(qslice, "")
	rows, err := db.DB.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

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
			return err
		}
		dupcmd := fmt.Sprintf(`CALL %s('%s',%d,'%s',%d,'%s','%s','%s','%s','%s','%s',
		    '%s','%s','%s','%s','%s','%s','%s',%d,'%s','%s',%d,%d,%d,'%s','%s','%s',
			'%s','%s',%d,%d,'%s','%s','%s','%s','%s','%s','%s',%d,%d);`,
			sqlCmd,
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
			ugc.Id,
			taskid)
		rowsline, err := db.DB.Query(dupcmd)
		if err != nil {
			return err
		}
		defer rowsline.Close()
	}
	if err := rows.Err(); err != nil {
		return err
	}
	return err
}
