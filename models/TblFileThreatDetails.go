/********获取文件威胁详情********/
package models

import (
	"apt-web-server/modules/mlog"
	"fmt"
	"strings"
)

func (this *TblFTD) TableName(tage string) string {
	switch tage {
	case "offline":
		return "alert_vds_offline"
	case "duplicate":
		return "alert_vds_duplicate"
	case "offlinedup":
		return "alert_vds_offline_duplicate"
	default:
		return "alert_vds"
	}
}

func GetFTDMysqlCmd(tablename string, para *TblFTDSearchPara) []string {
	qslice, whereflag := DefaultParaCmd("getlist", tablename, &para.PField)
	if para.Tage == "offline" {
		err, taskId := GetOfflineTaskID("offline_assignment", para.TaskName, para.CreateTime)
		if err != nil {
			fmt.Println("get taskid error")
		}
		if whereflag == 1 {
			temp_id := fmt.Sprintf(" and xdr_offline.Task_Id=%d", taskId)
			qslice = append(qslice, temp_id)
		} else {
			temp_id := fmt.Sprintf(" where xdr_offline.Task_Id=%d", taskId)
			qslice = append(qslice, temp_id)
			whereflag = 1
		}
	}
	if para.Tage == "offlinedup" || para.Tage == "duplicate" {
		err, taskId := GetOfflineTaskID("offline_assignment", para.TaskName, para.CreateTime)
		if err != nil {
			fmt.Println("get taskid error")
		}
		if whereflag == 1 {
			temp_id := fmt.Sprintf(" and %s.taskid=%d", tablename, taskId)
			qslice = append(qslice, temp_id)
		} else {
			temp_id := fmt.Sprintf(" where %s.taskid=%d", tablename, taskId)
			qslice = append(qslice, temp_id)
			whereflag = 1
		}
	}
	if para.Type != "" {
		if whereflag == 1 {
			temp_t := fmt.Sprintf(" and local_vtype='%s'", para.Type)
			qslice = append(qslice, temp_t)
		} else {
			temp_t := fmt.Sprintf(" where local_vtype='%s'", para.Type)
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
	} else if para.Count != 0 {
		para.Page = para.Page - 1
		temp_PC := fmt.Sprintf(" limit %d,%d", para.Page*para.Count, para.Count)
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
	rows, err := db.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	list := TblFTDData{}
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
			return err, nil
		}
		list.Elements = append(list.Elements, TblFTDList{ugc.TblFTDContent})
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	list.Counts = GetFTDCounts("", para)
	list.Totality = GetFTDCounts(this.TableName(para.Tage), para)

	return nil, &list
}

func (this *TblFTD) GetFTDetailsDUP(para *TblFTDSearchPara) (error, *TblFTDData) {
	qslice := GetFTDMysqlCmd(this.TableName(para.Tage), para)
	query := strings.Join(qslice, "")
	mlog.Debug(string(query))
	rows, err := db.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	list := TblFTDData{}
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
			&ugc.HttpUrl,
			&ugc.TblId,
			&ugc.TaskId)
		if err != nil {
			return err, nil
		}
		if para.Tage == "offlinedup" && para.MergeTag == "merge" {
			/*			mergecmd := fmt.Sprintf(`insert into alert_vds(log_time,threatname,subfile,
							local_threatname,local_vtype,local_platfrom,
							local_vname,local_extent,local_enginetype,
							local_logtype,local_engineip,sourceip,
							destip,sourceport,destport,app_file,http_url)
						value(%d,'%s','%s','%s','%s','%s','%s','%s','%s',
							'%s','%s','%s','%s',%d,%d,'%s','%s');`,
			*/
			mergecmd := fmt.Sprintf(`call checkvdsofflinemerge(%d,'%s','%s','%s','%s','%s','%s','%s','%s',
							'%s','%s','%s','%s',%d,%d,'%s','%s',%d,%d);`,
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
				ugc.TblId,
				ugc.TaskId)
			rows, err := db.Query(mergecmd)
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
	qslice, whereflag := DefaultParaCmd("getcounts", tablename, &para.PField)
	if tablename != "" {
		if para.Tage == "offline" {
			err, taskId := GetOfflineTaskID("offline_assignment", para.TaskName, para.CreateTime)
			if err != nil {
				fmt.Println("get taskid error")
			}
			if whereflag == 1 {
				temp_id := fmt.Sprintf(" and xdr_offline.Task_Id=%d", taskId)
				qslice = append(qslice, temp_id)
			} else {
				temp_id := fmt.Sprintf(" where xdr_offline.Task_Id=%d", taskId)
				qslice = append(qslice, temp_id)
				whereflag = 1
			}
		}
		if para.Tage == "offlinedup" || para.Tage == "duplicate" {
			err, taskId := GetOfflineTaskID("offline_assignment", para.TaskName, para.CreateTime)
			if err != nil {
				fmt.Println("get taskid error")
			}
			if whereflag == 1 {
				temp_id := fmt.Sprintf(" and %s.taskid=%d", tablename, taskId)
				qslice = append(qslice, temp_id)
			} else {
				temp_id := fmt.Sprintf(" where %s.taskid=%d", tablename, taskId)
				qslice = append(qslice, temp_id)
				whereflag = 1
			}
		}
		if para.Type != "" {
			if whereflag != 0 {
				temp_t := fmt.Sprintf(" and local_vtype='%s'", para.Type)
				qslice = append(qslice, temp_t)
			} else {
				temp_t := fmt.Sprintf(" where local_vtype='%s'", para.Type)
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

func (this *TblFTD) CreateSql() string {
	return fmt.Sprintf(
		`CREATE TABLE %s (
		id   integer unsigned  AUTO_INCREMENT NOT NULL,
		log_time   BIGINT NOT NULL DEFAULT 0,
		threatname varchar(20) NOT NULL DEFAULT '',
		subfile varchar(20) NOT NULL DEFAULT '',
		local_threatname text NULL ,
		local_vtype varchar(20) NOT NULL DEFAULT '',
		local_platfrom varchar(20) NOT NULL DEFAULT '',
		local_vname varchar(20) NOT NULL DEFAULT '',
		local_extent varchar(20) NOT NULL DEFAULT '',
		local_enginetype varchar(20) NOT NULL DEFAULT '',
		local_logtype varchar(20) NOT NULL DEFAULT '',
		local_engineip varchar(20) NOT NULL DEFAULT '',
		PRIMARY KEY (id),
		FULLTEXT(local_threatname)
		)ENGINE=MyISAM DEFAULT CHARSET=utf8;`,
		this.TableName(""))
}
