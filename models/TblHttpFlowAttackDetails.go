/********获取HTTP流量攻击详情********/
package models

import (
	"apt-web-server/modules/mlog"
	"encoding/base64"
	"fmt"
	"strings"
)

func (this *TblHFAD) TableName(tage string) string {
	switch tage {
	case "offline":
		return "alert_waf_offline"
	case "duplicate":
		return "alert_waf_duplicate"
	case "offlinedup":
		return "alert_waf_offline_duplicate"
	default:
		return "alert_waf"
	}
}

func GetHFADMysqlCmd(tablename string, para *TblHFADSearchPara) []string {
	qslice, whereflag := DefaultParaCmd("getlist", tablename, &para.PField)
	if para.Tage == "offline" {
		err, taskId := GetOfflineTaskID("offline_assignment", para.TaskName, para.CreateTime)
		if err != nil {
			fmt.Println("get taskid error")
		}
		if whereflag == 1 {
			temp_id := fmt.Sprintf(" and Task_Id=%d", taskId)
			qslice = append(qslice, temp_id)
		} else {
			temp_id := fmt.Sprintf(" where Task_Id=%d", taskId)
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
		if para.Type == "scaningprobe" {
			if whereflag == 1 {
				temp_t := fmt.Sprintf(" and (attack='reputation_scanner' or attack='reputation_scripting' or attack='reputation_crawler')")
				qslice = append(qslice, temp_t)
			} else {
				temp_t := fmt.Sprintf(" where (attack='reputation_scanner' or attack='reputation_scripting' or attack='reputation_crawler')")
				qslice = append(qslice, temp_t)
				whereflag = 1
			}
		} else {
			if whereflag == 1 {
				temp_t := fmt.Sprintf(" and attack='%s'", para.Type)
				qslice = append(qslice, temp_t)
			} else {
				temp_t := fmt.Sprintf(" where attack='%s'", para.Type)
				qslice = append(qslice, temp_t)
				whereflag = 1
			}
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

func (this *TblHFAD) GetHFADetails(para *TblHFADSearchPara) (error, *TblHFADData) {
	qslice := GetHFADMysqlCmd(this.TableName(para.Tage), para)
	query := strings.Join(qslice, "")
	mlog.Debug(string(query))
	rows, err := db.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	list := TblHFADData{}
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
			return err, nil
		}
		uri, _ := base64.StdEncoding.DecodeString(ugc.Uri)
		ugc.Uri = string(uri)
		ruledata, _ := base64.StdEncoding.DecodeString(ugc.RuleData)
		ugc.RuleData = string(ruledata)
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
	rows, err := db.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	list := TblHFADData{}
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
			&ugc.Version,
			&ugc.TblId,
			&ugc.TaskId)
		if err != nil {
			return err, nil
		}
		if para.Tage == "offlinedup" && para.MergeTag == "merge" {
			mergecmd := fmt.Sprintf(`call checkwafofflinemerge(%d,'%s','%s','%s','%s',%d,%d,%d,'%s','%s','%s','%s',
							'%s','%s',%d,%d,'%s','%s','%s',%d,%d);`,
				/*			mergecmd := fmt.Sprintf(`insert into alert_waf(time,client,rev,msg,
								attack,severity,maturity,accuracy,hostname,uri,
								unique_id,ref,tags,rule_file,rule_line,rule_id,
								rule_data,rule_ver,version)
							value(%d,'%s','%s','%s','%s',%d,%d,%d,'%s','%s','%s','%s',
								'%s','%s',%d,%d,'%s','%s','%s');`,
				*/
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
				ugc.TblId,
				ugc.TaskId)
			rows, err := db.Query(mergecmd)
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
		ruledata, _ := base64.StdEncoding.DecodeString(ugc.RuleData)
		ugc.RuleData = string(ruledata)
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
	qslice, whereflag := DefaultParaCmd("getcounts", tablename, &para.PField)
	if tablename != "" {
		if para.Tage == "offline" {
			err, taskId := GetOfflineTaskID("offline_assignment", para.TaskName, para.CreateTime)
			if err != nil {
				fmt.Println("get taskid error")
			}
			if whereflag == 1 {
				temp_id := fmt.Sprintf(" and Task_Id=%d", taskId)
				qslice = append(qslice, temp_id)
			} else {
				temp_id := fmt.Sprintf(" where Task_Id=%d", taskId)
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
			if para.Type == "scaningprobe" {
				if whereflag == 1 {
					temp_t := fmt.Sprintf(" and (attack='reputation_scanner' or attack='reputation_scripting' or attack='reputation_crawler')")
					qslice = append(qslice, temp_t)
				} else {
					temp_t := fmt.Sprintf(" where (attack='reputation_scanner' or attack='reputation_scripting' or attack='reputation_crawler')")
					qslice = append(qslice, temp_t)
					whereflag = 1
				}
			} else {
				if whereflag != 0 {
					temp_t := fmt.Sprintf(" and attack='%s'", para.Type)
					qslice = append(qslice, temp_t)
				} else {
					temp_t := fmt.Sprintf(" where attack='%s'", para.Type)
					qslice = append(qslice, temp_t)
				}
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

func (this *TblHFAD) CreateSql() string {
	return fmt.Sprintf(
		`CREATE TABLE %s (
		id   integer unsigned  AUTO_INCREMENT NOT NULL,
		time   BIGINT NOT NULL DEFAULT 0,
		src_ip varchar(20) NOT NULL DEFAULT '',
		src_port integer NOT NULL ,
		dest_ip varchar(20) NOT NULL DEFAULT '',
		dest_port integer NOT NULL ,
		matched_data text NULL,
		severity integer NULL,
		tags text NULL,
		attack varchar(20) NOT NULL DEFAULT '',
		msg text NULL,
		host varchar(20) NULL ,
		url varchar(20) NULL ,
		rule_id integer NULL,
		rule_file varchar(20) NULL ,
		rule_line integer NULL,
		rule_matched_arg text NULL,
		rulle_version varchar(20) NULL ,
		version varchar(20) NOT NULL DEFAULT '',
		PRIMARY KEY (id),
		FULLTEXT(msg),
		FULLTEXT(matched_data),
		FULLTEXT(rule_matched_arg),
		FULLTEXT(tags)
		)ENGINE=MyISAM DEFAULT CHARSET=utf8;`,
		this.TableName(""))
}
