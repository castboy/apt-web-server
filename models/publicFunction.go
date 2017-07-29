package models

import (
	"fmt"
)

func DefaultParaCmd(cmdType string, tablename string, para *TblPublicPara) ([]string, int) {
	qslice := make([]string, 0)
	flag := 0
	var qslice_tmp, xdr_tblname, tbl_tag string
	if tablename == "alert_vds_offline" ||
		tablename == "alert_waf_offline" ||
		tablename == "alert_vds_offline_duplicate" ||
		tablename == "alert_waf_offline_duplicate" {
		xdr_tblname = "xdr_offline"
	} else {
		xdr_tblname = "xdr"
	}
	if tablename == "alert_vds" ||
		tablename == "alert_vds_offline" ||
		tablename == "alert_vds_offline_duplicate" ||
		tablename == "alert_vds_duplicate" {
		tbl_tag = "vds"
	} else if tablename == "alert_waf" ||
		tablename == "alert_waf_offline" ||
		tablename == "alert_waf_offline_duplicate" ||
		tablename == "alert_waf_duplicate" {
		tbl_tag = "waf"
	}
	switch cmdType {
	case "getcounts":
		if tablename != "" {
			if tablename == "alert_vds" ||
				tablename == "alert_vds_offline" ||
				tablename == "alert_waf" ||
				tablename == "alert_waf_offline" {
				qslice_tmp = fmt.Sprintf(`select count(%s.XDR_Id) from %s,%s where Alert_Id=%s.id and Alert_Type='%s'`,
					xdr_tblname, tablename, xdr_tblname, tablename, tbl_tag)
				flag = 1
			} else {
				qslice_tmp = fmt.Sprintf(`select count(*) from %s`, tablename)
			}

		} else {
			qslice_tmp = fmt.Sprintf(`select FOUND_ROWS() as count`)
			qslice = append(qslice, qslice_tmp)
			return qslice, 0
		}
	case "getlist":
		if tablename == "alert_vds" ||
			tablename == "alert_vds_offline" ||
			tablename == "alert_waf" ||
			tablename == "alert_waf_offline" {
			qslice_tmp = fmt.Sprintf(`select %s.* from %s,%s where Alert_Id = %s.id and Alert_Type = '%s'`,
				tablename, tablename, xdr_tblname, tablename, tbl_tag)
			flag = 1
		} else if tablename == "alert_vds_offline_duplicate" ||
			tablename == "alert_vds_duplicate" ||
			tablename == "alert_waf_offline_duplicate" ||
			tablename == "alert_waf_duplicate" {
			qslice_tmp = fmt.Sprintf(`select %s.* from %s,%s where Alert_Id = %s.tblid and Alert_Type = '%s'`,
				tablename, tablename, xdr_tblname, tablename, tbl_tag)
			flag = 1
		} else {
			qslice_tmp = fmt.Sprintf(`select * from %s`, tablename)
		}
	}
	qslice = append(qslice, qslice_tmp)
	if para.Start != 0 && para.End != 0 {
		var time_tag, temp_se string
		if tablename == "alert_vds" {
			time_tag = "log_time"
		} else if tablename == "alert_waf" {
			time_tag = "alert_waf.time"
		} else {
			time_tag = "time"
		}
		if flag == 1 {
			temp_se = fmt.Sprintf(" and %s>=%d and %s<=%d ",
				time_tag, para.Start, time_tag, para.End)
		} else {
			temp_se = fmt.Sprintf(" where %s>=%d and %s<=%d",
				time_tag, para.Start, time_tag, para.End)
			flag = 1
		}
		qslice = append(qslice, temp_se)
	}
	return qslice, flag
}
func GetOfflineTaskID(tblname, taskname string, time int64) (error, int) {
	var taskID int
	query := fmt.Sprintf(`select id from %s where name='%s' and time=%d;`,
		tblname,
		taskname,
		time)
	fmt.Println(query)
	rows, err := db.Query(query)
	if err != nil {
		//return err, nil
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(
			&taskID)
		if err != nil {
			return err, 0
		}
	}
	return nil, taskID
}
func GetTimes(taskId int) (int64, int64) {
	var maxtime, mintime int64
	query := fmt.Sprintf(`select max(Time),min(Time) from %s where taskid=%d;`,
		"xdr_offline",
		taskId)
	rows, err := db.Query(query)
	if err != nil {
		return 0, 0
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(
			&maxtime,
			&mintime)
		if err != nil {
			return 0, 0
		}
	}
	return maxtime, mintime
}
func DuplicateVds(tblName, taskName string, time int64) error {
	var query, sqlCmd string
	err, taskid := GetOfflineTaskID("offline_assignment", taskName, time)
	if err != nil {
		fmt.Println("get taskid error!")
	}
	if tblName == "alert_vds" {
		sqlCmd = "checkvdsduplicate"
		maxTime, minTime := GetTimes(taskid)
		query = fmt.Sprintf("select * from %s where time>=%d and time<=%d ;",
			tblName, minTime, maxTime)
	}
	if tblName == "alert_vds_offline" {
		sqlCmd = "checkvdsofflinedup"
		query = fmt.Sprintf("select %s.* from %s,xdr_offline where xdr_offline.Task_Id=%d;",
			tblName, tblName, taskid)
	}
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		ugc := new(TblFTD)
		err = rows.Scan(
			&ugc.Id,
			&ugc.LogTime,
			&ugc.ThreatName,
			&ugc.SubFile,
			&ugc.LocalThreatName,
			&ugc.LocalVType,
			&ugc.LocalPlatfrom,
			&ugc.LocalVName,
			&ugc.LocalExtent,
			&ugc.LocalEngineType,
			&ugc.LocalLogType,
			&ugc.LocalEngineIP,
			&ugc.SrcIp,
			&ugc.DestIp,
			&ugc.SrcPort,
			&ugc.DestPort,
			&ugc.AppFile,
			&ugc.HttpUrl)
		if err != nil {
			return err
		}
		dupcmd := fmt.Sprintf(`call %s(%d,'%s','%s','%s','%s','%s','%s','%s',
			'%s','%s','%s','%s','%s',%d,%d,'%s','%s',%d,%d);`,
			sqlCmd,
			ugc.LogTime,
			ugc.ThreatName,
			ugc.SubFile,
			ugc.LocalThreatName,
			ugc.LocalVType,
			ugc.LocalPlatfrom,
			ugc.LocalVName,
			ugc.LocalExtent,
			ugc.LocalEngineType,
			ugc.LocalLogType,
			ugc.LocalEngineIP,
			ugc.SrcIp,
			ugc.DestIp,
			ugc.SrcPort,
			ugc.DestPort,
			ugc.AppFile,
			ugc.HttpUrl,
			ugc.Id,
			taskid)
		rowsline, err := db.Query(dupcmd)
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
func DuplicateWaf(tblName, taskName string, time int64) error {
	var query, sqlCmd string
	err, taskid := GetOfflineTaskID("offline_assignment", taskName, time)
	if err != nil {
		fmt.Println("get taskid error!")
	}

	if tblName == "alert_waf" {
		sqlCmd = "checkwafduplicate"
		maxTime, minTime := GetTimes(taskid)
		query = fmt.Sprintf("select * from %s where time>=%d and time<=%d;",
			tblName, minTime, maxTime)
	}
	if tblName == "alert_waf_offline" {
		sqlCmd = "checkwafofflinedup"
		query = fmt.Sprintf("select %s.* from %s,xdr_offline where xdr_offline.Task_Id=%d;",
			tblName, tblName, taskid)
	}
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		ugc := new(TblHFAD)
		err = rows.Scan(
			&ugc.Id,
			&ugc.Time,
			&ugc.Client,
			&ugc.Rev,
			&ugc.Message,
			&ugc.Attack,
			&ugc.Severity,
			&ugc.Maturity,
			&ugc.Accuracy,
			&ugc.HostName,
			&ugc.Uri,
			&ugc.UniqueId,
			&ugc.Ref,
			&ugc.Tags,
			&ugc.RuleFile,
			&ugc.RuleLine,
			&ugc.RuleId,
			&ugc.RuleData,
			&ugc.RuleVersion,
			&ugc.Version)
		if err != nil {
			return err
		}
		dupcmd := fmt.Sprintf(`call %s(%d,'%s','%s','%s','%s',%d,%d,%d,'%s',
			'%s','%s','%s','%s','%s',%d,%d,'%s','%s','%s',%d,%d);`,
			sqlCmd,
			ugc.Time,
			ugc.Client,
			ugc.Rev,
			ugc.Message,
			ugc.Attack,
			ugc.Severity,
			ugc.Maturity,
			ugc.Accuracy,
			ugc.HostName,
			ugc.Uri,
			ugc.UniqueId,
			ugc.Ref,
			ugc.Tags,
			ugc.RuleFile,
			ugc.RuleLine,
			ugc.RuleId,
			ugc.RuleData,
			ugc.RuleVersion,
			ugc.Version,
			ugc.Id,
			taskid)
		rowsline, err := db.Query(dupcmd)
		if err != nil {
			return err
		}
		fmt.Println(dupcmd)
		defer rowsline.Close()
	}
	if err := rows.Err(); err != nil {
		return err
	}
	return err
}
