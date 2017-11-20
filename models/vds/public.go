package vds

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
	err, taskid := modelsPublic.GetOfflineTaskID("offline_assignment", taskName, time)
	if err != nil {
		fmt.Println("get taskid error!")
		return nil, "", 0, err
	}
	column := fmt.Sprintf(` id,src_ip,src_port,dest_ip,dest_port,proto,src_country,
	    src_province,src_city,src_latitude,src_longitude,dest_country,
		dest_province,dest_city,dest_latitude,dest_longitude,operators,time,
		local_vtype,local_vname,local_logtype,threatname,subfile,
		local_threatname,local_platfrom,local_extent,local_enginetype,
		local_engineip,app_file,http_url`)
	switch tblName {
	case "alert_vds":
		sqlCmd = "checkvdsduplicate"
		maxTime, minTime := modelsPublic.GetTimes("alert_vds_offline", taskid)
		qslice = append(qslice, column)
		ruleList = fmt.Sprintf(` WHERE time BETWEEN %d AND %d;`, maxTime, minTime)
	case "alert_vds_offline":
		sqlCmd = "checkvdsofflinedup"
		qslice = append(qslice, column)
		ruleList = fmt.Sprintf(` WHERE task_id=%d;`, tblName, taskid)
	}
	fromTbl := fmt.Sprintf(` FROM %s`, tblName)
	qslice = append(qslice, fromTbl)
	qslice = append(qslice, ruleList)
	return qslice, sqlCmd, taskid, nil
}
func DuplicateVds(tblName, taskName string, time int64) error {
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
			return err
		}
		dupcmd := fmt.Sprintf(`CALL %s('%s',%d,'%s',%d,'%s','%s','%s','%s','%s','%s',
		                    '%s','%s','%s','%s','%s','%s',%d,'%s','%s','%s','%s',
							'%s','%s','%s','%s','%s','%s','%s','%s',%d,%d);`,
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
