package offlineAssignment

import (
	"apt-web-server_v2/models/db"
	"apt-web-server_v2/models/modelsPublic"
	"apt-web-server_v2/modules/mlog"
	"errors"
	"fmt"
	"strings"
)

func GetTLMsqlCmd(tablename string, para *TblTaskSearchPara) []string {
	qslice, flag := modelsPublic.DefaultParaCmd("getlist", tablename, &para.PField)
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
	rows, err := db.DB.Query(query)
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
			&ugc.RuleSet,
			//&ugc.Rule,
			//&ugc.Rule2,
			//&ugc.Rule3,
			//&ugc.Rule4,
			//&ugc.Rule5,
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

		ruleInfoSingle := new(TblTaskList /*TaskContent2UI*/)

		ruleInfoSingle.Name = ugc.Name
		//ruleInfoSingle.RuleSet
		ruleInfoSingle.CreateTime = ugc.CreateTime
		ruleInfoSingle.Type = ugc.Type
		ruleInfoSingle.DataStart = ugc.DataStart
		ruleInfoSingle.DataEnd = ugc.DataEnd
		ruleInfoSingle.Weight = ugc.Weight
		ruleInfoSingle.Topic = ugc.Topic
		ruleInfoSingle.Status = ugc.Status
		ruleInfoSingle.Details = ugc.Details

		if para.OfflineTag == "rule" {

			//ruleAndIdSet := strings.Split(ugc.RuleSet, "|")
			//for idx, _ := range ruleAndIdSet {
			//	if len(ruleAndIdSet[idx]) > 0 {
			//		ruleAndId := strings.Split(ruleAndIdSet[idx], ":")
			//		IdTmp, err := strconv.Atoi(ruleAndId[0])
			//		if err != nil {
			//			rspMsg := fmt.Sprintf(`自定义檢測任務查找失敗，規則id -> %s 非法。`, ruleAndId[0])
			//			err1 := errors.New(rspMsg)
			//			return err1, &list
			//		}
			//		ruleTmp1 := new(TblRuleSdSet)
			//		ruleTmp1.Id = int64(IdTmp)
			//		ruleTmp1.Rule = ruleAndId[1]
			//		ruleInfoSingle.RuleSets = append(ruleInfoSingle.RuleSets, *ruleTmp1 /*TblRuleSdSet{ruleTmp1}*/)
			//	}
			//}

			/////////////////////////////////////////////////
			// 獲取rule信息，然後將rule的 id， 名字 填寫到返回消息list中
			searchPara := TblRuleSdSearchPara{}
			searchPara.RuleSet = ugc.RuleSet
			searchPara.Count = 100000
			searchPara.Page = 1
			searchPara.Type = "byidset" /*"byalias"*/

			errGetRule, listRule := GetRuleSdLst(&searchPara)
			if errGetRule != nil || listRule == nil {
				//這個異常沒有人 recover
				//panic(fmt.Sprintf("GetRuleSdLst error:%s", errGetRule.Error()))

				fmt.Println("GetSingleRuleConfText, call GetRuleSdLst failed, errGetRule is ", errGetRule.Error())
				rspMsg := fmt.Sprintf(`规则 %s 不存在。 %s 。`, ugc.RuleSet, errGetRule.Error())
				errGetRuleR := errors.New(rspMsg)
				return errGetRuleR, nil
			}

			ruleSetTmp := make([]string, 0)
			for idx, _ := range listRule.Elements {
				//ruleTmp1 := new(TblRuleSdSet)
				//ruleTmp1.Id = int64(listRule.Elements[idx].Id)
				//ruleTmp1.Rule = listRule.Elements[idx].Alias
				if idx > 0 {
					ruleSetTmp = append(ruleSetTmp, ",")
				}
				ruleSetTmp = append(ruleSetTmp, listRule.Elements[idx].Alias)
			}
			ruleInfoSingle.RuleSets = strings.Join(ruleSetTmp, "")
			/////////////////////////////////////////////////
		}

		list.Elements = append(list.Elements, *ruleInfoSingle /*TblTaskList{ugc.TaskContent}*/)
		//list.Elements = append(list.Elements, TblTaskList{ugc.TaskContent})
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	list.Counts = GetTLCounts("", para)
	list.Totality = GetTLCounts(this.TableName(para.OfflineTag), para)
	return nil, &list
}
func GetTLCounts(tablename string, para *TblTaskSearchPara) int64 {
	qslice, whereflag := modelsPublic.DefaultParaCmd("getcounts", tablename, &para.PField)
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
