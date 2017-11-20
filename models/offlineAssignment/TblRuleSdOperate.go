/********自定义规则 规则的处理相关的代码********/
//package models
package offlineAssignment

import (
	//"encoding/base64"
	"apt-web-server_v2/models/db"
	//"apt-web-server_v2/models/modelsPublic"
	//"apt-web-server_v2/modules/mlog"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	//"time"
)

func Rule_OL_TableName() string {
	return "assignment_rule"
}

func Rule_Defined_Add(rulePara *TblRuleOperPara) (err error) {
	var db_cmd string
	db_cmd_slice := make([]string, 0)
	var query_cmd string
	var count int
	//var VarCnt int
	//var VarCnt2 int

	query_cmd = fmt.Sprintf(`select count(id) from %s where alias='%s'`,
		Rule_OL_TableName(), rulePara.Alias)
	fmt.Println("Rule_Defined_Add query_cmd is ", query_cmd)
	//rows := Select_mysql(query_cmd)
	rows, err := db.DB.Query(query_cmd)
	if err != nil {
		rspMsg := fmt.Sprintf(`添加规则失败。 query error, %s 。`, err.Error())
		err1 := errors.New(rspMsg)
		return err1
	}

	defer rows.Close()
	//check if the rule has already existed.
	for rows.Next() {
		if rows != nil {
			if err = rows.Scan(&count); err != nil {
				fmt.Println("Rule_Defined_Add : rows.Scan failed, err is  ", err.Error())
				rspMsg := fmt.Sprintf(`Rule_Defined_Add : rows.Scan failed, err is %s 。`, err.Error())
				err = errors.New(rspMsg)
				return err
			}

			fmt.Println("Rule_Defined_Add query_cmd count is ", count)

			if count != 0 {
				//rows.Close()
				rspMsg := fmt.Sprintf(`该规则已经存在，规则别名为 %s 。`, rulePara.Alias)
				err = errors.New(rspMsg)
				return err
			}
		}
		break
	}

	//fmt.Println("Rule_Defined_Add out of for, query_cmd count is ", count)

	// insert rule into db table
	//db_cmd_p1 := fmt.Sprintf(`INSERT INTO %s `, Rule_OL_TableName())
	//db_cmd = append(db_cmd, db_cmd_p1)
	//db_cmd = fmt.Sprintf(`INSERT INTO %s (alias,varsnum,vars,varsinfo,vars2,varsinfo2,vars3,varsinfo3,vars4,varsinfo4,vars5,varsinfo5,oper,operinfo,phase,severity,accuracy,maturity,tag,details)
	//	VALUES('%s',%d,'%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s')`,
	//	Rule_OL_TableName(), rulePara.Alias, len(rulePara.VarSet),
	//	rulePara.VarSet[0].Var, rulePara.VarSet[0].VarInfo, rulePara.VarSet[1].Var, rulePara.VarSet[1].VarInfo, rulePara.VarSet[2].Var, rulePara.VarSet[2].VarInfo, rulePara.VarSet[3].Var, rulePara.VarSet[3].VarInfo, rulePara.VarSet[4].Var, rulePara.VarSet[4].VarInfo,
	//	rulePara.Oper, rulePara.OperInfo, rulePara.Phase, rulePara.Severity, rulePara.Accuracy, rulePara.Maturity, rulePara.Tag, rulePara.Details)

	db_cmd_01 := fmt.Sprintf(`INSERT INTO %s (alias,varset,`,
		Rule_OL_TableName() /*, rulePara.Alias, len(rulePara.VarSet)*/)
	db_cmd_slice = append(db_cmd_slice, db_cmd_01)

	//for idx, _ := range rulePara.VarSet {
	//	if len(rulePara.VarSet[idx].Var) > 0 {
	//		VarCnt++
	//		if VarCnt > 1 {
	//			db_cmd_tmp := fmt.Sprintf(`vars%d,varsinfo%d,`, idx+1, idx+1)
	//			db_cmd_slice = append(db_cmd_slice, db_cmd_tmp)
	//		}
	//	}
	//}

	db_cmd_03 := fmt.Sprintf(`oper,operinfo,tfunc,phase,severity,accuracy,maturity,tag,details)
		VALUES('%s','%s','%s', '%s','%s','%s','%s','%s','%s','%s','%s')`,
		rulePara.Alias, rulePara.VarSet, /*rulePara.VarSet[0].Var, rulePara.VarSet[0].VarInfo*/
		rulePara.Oper, rulePara.OperInfo, rulePara.TransFunc, rulePara.Phase, rulePara.Severity, rulePara.Accuracy, rulePara.Maturity, rulePara.Tag, rulePara.Details)
	db_cmd_slice = append(db_cmd_slice, db_cmd_03)

	//for idx2, _ := range rulePara.VarSet {
	//	VarCnt2++
	//	if VarCnt2 > 1 {
	//		db_cmd_tmp := fmt.Sprintf(`'%s','%s',`, rulePara.VarSet[idx2].Var, rulePara.VarSet[idx2].VarInfo)
	//		db_cmd_slice = append(db_cmd_slice, db_cmd_tmp)
	//	}
	//}

	//db_cmd_05 := fmt.Sprintf(`'%s','%s','%s','%s','%s','%s','%s','%s')`,
	//	rulePara.Oper, rulePara.OperInfo, rulePara.Phase, rulePara.Severity, rulePara.Accuracy, rulePara.Maturity, rulePara.Tag, rulePara.Details)
	//db_cmd_slice = append(db_cmd_slice, db_cmd_05)

	//db_cmd_slice = append(db_cmd_slice, ";")
	db_cmd = strings.Join(db_cmd_slice, "")

	fmt.Println("Rule_Defined_Add db_cmd for insert is :", db_cmd)

	rows, err1 := db.DB.Query(db_cmd)
	defer rows.Close()
	if err1 != nil {
		//mlog.Debug(query, "CreatAssignment error")
		rspMsg := fmt.Sprintf(`添加规则失败。 insert error, %s 。`, err1.Error())
		err = errors.New(rspMsg)
		return err
	}

	return err
}

func Rule_Defined_Mod(rulePara *TblRuleOperPara) (err error) {
	var db_cmd string
	db_cmd_slice := make([]string, 0)
	var query_cmd string
	var count int
	var aliasTmp string
	//var VarCnt int

	//query_cmd = fmt.Sprintf(`select count(id) from %s where alias='%s'`,
	//	Rule_OL_TableName(), rulePara.Alias)
	query_cmd = fmt.Sprintf(`select count(id),alias from %s where id=%d`,
		Rule_OL_TableName(), rulePara.Id)
	fmt.Println("Rule_Defined_Mod query_cmd is ", query_cmd)
	//rows := Select_mysql(query_cmd)
	rows, err := db.DB.Query(query_cmd)
	if err != nil {
		rspMsg := fmt.Sprintf(`修改规则失败。 query error, %s 。`, err.Error())
		err1 := errors.New(rspMsg)
		return err1
	}

	defer rows.Close()
	//check if the rule has already existed.
	for rows.Next() {
		if rows != nil {
			if err = rows.Scan(&count, &aliasTmp); err != nil {
				fmt.Println("Rule_Defined_Add : rows.Scan failed, err is  ", err.Error())
				rspMsg := fmt.Sprintf(`Rule_Defined_Mod : rows.Scan failed, err is %s 。`, err.Error())
				err = errors.New(rspMsg)
				return err
			}
			fmt.Println("Rule_Defined_Mod query_cmd count is ", count)
			if count == 0 {
				//rows.Close()
				rspMsg := fmt.Sprintf(`想要修改的规则并不存在，规则别名为 %s, id为 %d 。`, rulePara.Alias, rulePara.Id)
				//rspMsg := fmt.Sprintf(`想要修改的规则并不存在，规则 Id 为 %d 。`, rulePara.Id)
				err = errors.New(rspMsg)
				return err
			}
		}
		break
	}

	////////////////////////////////////////////////////////////////////////////
	if aliasTmp != rulePara.Alias {
		// check if the new rule name has already existed.
		alias_cmd := fmt.Sprintf(`select count(id) from %s where alias='%s'`,
			Rule_OL_TableName(), rulePara.Alias)
		fmt.Println("Rule_Defined_Mod cmd for query of new alias is :", alias_cmd)
		rowtmps, QAerr := db.DB.Query(alias_cmd)
		if QAerr != nil {
			rspMsg := fmt.Sprintf(`修改规则失败。 query for new Alias error, %s 。`, QAerr.Error())
			err1 := errors.New(rspMsg)
			return err1
		}

		defer rowtmps.Close()
		for rowtmps.Next() {
			if rowtmps != nil {
				if err = rowtmps.Scan(&count); err != nil {
					fmt.Println("Rule_Defined_Mod : rowtmps.Scan failed, err is  ", err.Error())
					rspMsg := fmt.Sprintf(`Rule_Defined_Mod : rowtmpsrowtmps.Scan failed, err is %s 。`, err.Error())
					err = errors.New(rspMsg)
					return err
				}
				//fmt.Println("Rule_Defined_Mod alias_cmd count is ", count)
				if count != 0 {
					//rowtmps.Close()
					rspMsg := fmt.Sprintf(`修改规则失败，新的规则别名已经存在，新的规则别名为 %s, id为 %d 。`, rulePara.Alias, rulePara.Id)
					//rspMsg := fmt.Sprintf(`想要修改的规则并不存在，规则 Id 为 %d 。`, rulePara.Id)
					err = errors.New(rspMsg)
					return err
				}
			}
			break
		}
	}

	////////////////////////////////////////////////////////////////////////////
	// insert rule into db table
	//db_cmd = fmt.Sprintf(`REPLACE INTO %s (alias,varsnum,vars,varsinfo,vars2,varsinfo2,vars3,varsinfo3,vars4,varsinfo4,vars5,varsinfo5,oper,operinfo,phase,severity,accuracy,maturity,tag,details)
	//	VALUES('%s',%d,'%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s')`,
	//	Rule_OL_TableName(), rulePara.Alias, len(rulePara.VarSet),
	//	rulePara.VarSet[0].Var, rulePara.VarSet[0].VarInfo, rulePara.VarSet[1].Var, rulePara.VarSet[1].VarInfo, rulePara.VarSet[2].Var, rulePara.VarSet[2].VarInfo, rulePara.VarSet[3].Var, rulePara.VarSet[3].VarInfo, rulePara.VarSet[4].Var, rulePara.VarSet[4].VarInfo,
	//	rulePara.Oper, rulePara.OperInfo, rulePara.Phase, rulePara.Severity, rulePara.Accuracy, rulePara.Maturity, rulePara.Tag, rulePara.Details)

	//db_cmd = fmt.Sprintf(`UPDATE %s SET alias='%s',varsnum=%d,vars='%s',varsinfo='%s',vars2='%s',varsinfo2'%s',vars3='%s',varsinfo3='%s',vars4='%s',varsinfo4='%s',vars5='%s',varsinfo5='%s',oper='%s',operinfo='%s',phase='%s',severity='%s',accuracy='%s',maturity='%s',tag='%s',details='%s' WHERE alias='%s'`,
	//	Rule_OL_TableName(), rulePara.Alias, len(rulePara.VarSet),
	//	rulePara.VarSet[0].Var, rulePara.VarSet[0].VarInfo, rulePara.VarSet[1].Var, rulePara.VarSet[1].VarInfo, rulePara.VarSet[2].Var, rulePara.VarSet[2].VarInfo, rulePara.VarSet[3].Var, rulePara.VarSet[3].VarInfo, rulePara.VarSet[4].Var, rulePara.VarSet[4].VarInfo,
	//	rulePara.Oper, rulePara.OperInfo, rulePara.Phase, rulePara.Severity, rulePara.Accuracy, rulePara.Maturity, rulePara.Tag, rulePara.Details, rulePara.Alias)

	db_cmd_01 := fmt.Sprintf(`UPDATE %s SET alias='%s',varset='%s',oper='%s',operinfo='%s',tfunc='%s',phase='%s',severity='%s',accuracy='%s',maturity='%s',tag='%s',details='%s' WHERE id=%d`,
		Rule_OL_TableName(), rulePara.Alias, rulePara.VarSet,
		rulePara.Oper, rulePara.OperInfo, rulePara.TransFunc, rulePara.Phase, rulePara.Severity, rulePara.Accuracy, rulePara.Maturity, rulePara.Tag, rulePara.Details, rulePara.Id)
	db_cmd_slice = append(db_cmd_slice, db_cmd_01)

	//for idx, _ := range rulePara.VarSet {
	//	if len(rulePara.VarSet[idx].Var) > 0 {
	//		VarCnt++
	//		if VarCnt > 1 {
	//			db_cmd_tmp := fmt.Sprintf(`vars%d='%s',varsinfo%d='%s',`,
	//				idx+1, rulePara.VarSet[idx].Var, idx+1, rulePara.VarSet[idx].VarInfo)
	//			db_cmd_slice = append(db_cmd_slice, db_cmd_tmp)
	//		}
	//	}
	//}

	//db_cmd_03 := fmt.Sprintf(`oper='%s',operinfo='%s',phase='%s',severity='%s',accuracy='%s',maturity='%s',tag='%s',details='%s' WHERE alias='%s'`,
	//	rulePara.Oper, rulePara.OperInfo, rulePara.Phase, rulePara.Severity, rulePara.Accuracy, rulePara.Maturity, rulePara.Tag, rulePara.Details, rulePara.Alias)
	//db_cmd_slice = append(db_cmd_slice, db_cmd_03)

	//db_cmd_slice = append(db_cmd_slice, ";")
	db_cmd = strings.Join(db_cmd_slice, "")

	fmt.Println("Rule_Defined_Mod db_cmd for mod is ", db_cmd)

	rows, err1 := db.DB.Query(db_cmd)
	defer rows.Close()
	if err1 != nil {
		//mlog.Debug(query, "Rule_Defined_Mod error")
		rspMsg := fmt.Sprintf(`修改规则失败。 proc error, %s 。`, err1.Error())
		err = errors.New(rspMsg)
		return err
	}

	return err
}

func Rule_Defined_Del(rulePara *TblRuleOperPara) (err error) {
	var db_cmd string
	//db_cmd := make([]string, 0)
	var query_cmd string
	//var count int

	//query_cmd = fmt.Sprintf(`select count(id) from %s where alias='%s'`,
	//	Rule_OL_TableName(), rulePara.Alias)
	query_cmd = fmt.Sprintf(`select count(id) from %s where id=%d`,
		Rule_OL_TableName(), rulePara.Id)
	fmt.Println("Rule_Defined_Del query_cmd is ", query_cmd)
	//rows := Select_mysql(query_cmd)
	rows, err := db.DB.Query(query_cmd)
	if err != nil {
		rspMsg := fmt.Sprintf(`删除规则失败。 query error, %s 。`, err.Error())
		err1 := errors.New(rspMsg)
		return err1
	}

	defer rows.Close()
	//check if the rule has already existed.
	//for rows.Next() {
	//	if rows != nil {
	//		if err = rows.Scan(&count); err != nil {
	//			fmt.Println("Rule_Defined_Del : rows.Scan failed, err is  ", err.Error())
	//			rspMsg := fmt.Sprintf(`Rule_Defined_Del : rows.Scan failed, err is %s 。`, err.Error())
	//			err = errors.New(rspMsg)
	//			return err
	//		}
	//		//fmt.Println("Rule_Defined_Del query_cmd count is ", count)
	//		if count == 0 {
	//			//rows.Close()
	//			rspMsg := fmt.Sprintf(`要删除的规则并不存在，规则别名为 %s 。`, rulePara.Alias)
	//			err = errors.New(rspMsg)
	//			return err
	//		}
	//	}
	//	break
	//}

	// delete rule from db table
	//db_cmd_p1 := fmt.Sprintf(`INSERT INTO %s `, Rule_OL_TableName())
	//db_cmd = append(db_cmd, db_cmd_p1)
	//db_cmd = fmt.Sprintf(`DELETE FROM %s WHERE alias='%s'`,
	//	Rule_OL_TableName(), rulePara.Alias)
	//db_cmd = fmt.Sprintf(`DELETE FROM %s WHERE id=%d`,
	//	Rule_OL_TableName(), rulePara.Id)
	db_cmd = fmt.Sprintf(`DELETE FROM %s WHERE id IN (%s)`,
		Rule_OL_TableName(), rulePara.IdSet)
	fmt.Println("Rule_Defined_Del db_cmd for del is ", db_cmd)

	rows, err1 := db.DB.Query(db_cmd)
	defer rows.Close()
	if err1 != nil {
		//mlog.Debug(query, "CreatAssignment error")
		rspMsg := fmt.Sprintf(`删除规则失败。 insert error, %s 。`, err1.Error())
		err = errors.New(rspMsg)
		return err
	}

	return err
}

func GetRuleSdLst(para *TblRuleSdSearchPara) (error, *TblRuleSdLstData) {
	//datemold := GetDateMold(para.Unit)  datemold = "%Y-%m-%d"
	datemold := "%Y-%m-%d %H:%i:%s"
	var cnt int64 = 0
	list := TblRuleSdLstData{}
	//VarEle := TblRuleSdLstVarSet{}

	qslice := GetRuleSdSearchMysqlCmd(datemold, Rule_OL_TableName(), para)
	query := strings.Join(qslice, "")

	fmt.Println(qslice)
	fmt.Println(query)

	rows, err := db.DB.Query(query)
	if err != nil {
		fmt.Println("GetRuleSdLst, call db.DB.Query failed, err is :", err.Error())
		return err, nil
	}
	defer rows.Close()
	//ugc := new(TblWLLst)
	var ugc TblRuleSdLstSame2Tbl
	//var wlinfo TblRuleSdLstContent //TblWLLstContent

	for rows.Next() {
		err = rows.Scan(
			&ugc.Id,
			&ugc.Alias,
			&ugc.VarSet,
			//
			//&ugc.Var,
			//&ugc.VarInfo,
			//&ugc.Var2,
			//&ugc.VarInfo2,
			//&ugc.Var3,
			//&ugc.VarInfo3,
			//&ugc.Var4,
			//&ugc.VarInfo4,
			//&ugc.Var5,
			//&ugc.VarInfo5,
			&ugc.Oper,
			&ugc.OperInfo,
			&ugc.TransFunc,
			&ugc.Phase,
			&ugc.Severity,
			&ugc.Accuracy,
			&ugc.Maturity,
			&ugc.Tag,
			&ugc.Details)
		if err != nil {
			fmt.Println("GetRuleSdLst, call rows.Scan failed, err is :", err.Error())
			return err, nil
		}
		wlinfo := new(TblRuleSdLstContent)
		// compose result struct
		wlinfo.Id = ugc.Id
		//aliasEncodedStr := base64.StdEncoding.EncodeToString([]byte(ugc.Alias))
		wlinfo.Alias = ugc.Alias /*aliasEncodedStr*/
		wlinfo.VarSet = ugc.VarSet
		//if len(ugc.Var) > 0 {
		//	VarEle.Var = ugc.Var
		//	//VIStrTmp1 := base64.StdEncoding.EncodeToString([]byte(ugc.VarInfo))
		//	VarEle.VarInfo = ugc.VarInfo /*VIStrTmp1*/
		//	wlinfo.VarSet = append(wlinfo.VarSet, VarEle)
		//}
		//if len(ugc.Var2) > 0 {
		//	VarEle.Var = ugc.Var2
		//	//VIStrTmp2 := base64.StdEncoding.EncodeToString([]byte(ugc.VarInfo2))
		//	VarEle.VarInfo = ugc.VarInfo2 /*VIStrTmp2*/
		//	wlinfo.VarSet = append(wlinfo.VarSet, VarEle)
		//}
		//if len(ugc.Var3) > 0 {
		//	VarEle.Var = ugc.Var3
		//	//VIStrTmp3 := base64.StdEncoding.EncodeToString([]byte(ugc.VarInfo3))
		//	VarEle.VarInfo = ugc.VarInfo3 /*VIStrTmp3*/
		//	wlinfo.VarSet = append(wlinfo.VarSet, VarEle)
		//}
		//if len(ugc.Var4) > 0 {
		//	VarEle.Var = ugc.Var4
		//	//VIStrTmp4 := base64.StdEncoding.EncodeToString([]byte(ugc.VarInfo4))
		//	VarEle.VarInfo = ugc.VarInfo4 /*VIStrTmp4*/
		//	wlinfo.VarSet = append(wlinfo.VarSet, VarEle)
		//}
		//if len(ugc.Var5) > 0 {
		//	VarEle.Var = ugc.Var5
		//	//VIStrTmp5 := base64.StdEncoding.EncodeToString([]byte(ugc.VarInfo5))
		//	VarEle.VarInfo = ugc.VarInfo5 /*VIStrTmp5*/
		//	wlinfo.VarSet = append(wlinfo.VarSet, VarEle)
		///}
		wlinfo.Oper = ugc.Oper
		//OperStrTmp := base64.StdEncoding.EncodeToString([]byte(ugc.OperInfo))
		wlinfo.OperInfo = ugc.OperInfo /*OperStrTmp*/
		wlinfo.TransFunc = ugc.TransFunc
		wlinfo.Phase = ugc.Phase
		wlinfo.Severity = ugc.Severity
		wlinfo.Accuracy = ugc.Accuracy
		wlinfo.Maturity = ugc.Maturity
		//TagStrTmp := base64.StdEncoding.EncodeToString([]byte(ugc.Tag))
		wlinfo.Tag = ugc.Tag /*TagStrTmp*/
		wlinfo.Details = ugc.Details

		list.Elements = append(list.Elements, *wlinfo /*ugc.TblWLLstContent*/)
		cnt++
	}
	if err := rows.Err(); err != nil {
		fmt.Println("GetRuleSdLst, end failed, err is :", err.Error())
		return err, nil
	}
	list.Counts = cnt //GetSSLCount("", para)
	if cnt != 0 {
		//list.Totality = GetWLLstCount(this.GetTableName(), para)
		list.Totality = GetRuleSdLstCountOnCondition(Rule_OL_TableName(), para)
	}

	return nil, &list
}

func GetRuleSdSearchMysqlCmd(datemold, tablename string, para *TblRuleSdSearchPara) []string {
	var qslice_info string
	qslice := make([]string, 0)

	fmt.Println("GetRuleSdSearchMysqlCmd , input para is ", para)

	switch para.Type {
	case "all":
		//qslice_info = fmt.Sprintf(`SELECT id, alias, vars,varsinfo, vars2,varsinfo2, vars3,varsinfo3, vars4,varsinfo4, vars5,varsinfo5,
		//	oper, operinfo, phase, severity, accuracy, maturity, tag, details
		//	FROM %s`,
		//	tablename)
		qslice_info = fmt.Sprintf(`SELECT id, alias, varset,
			oper, operinfo, tfunc, phase, severity, accuracy, maturity, tag, details	
			FROM %s`,
			tablename)
		qslice = append(qslice, qslice_info)
	case "byalias":
		//qslice_info = fmt.Sprintf(`SELECT id, alias, vars,varsinfo, vars2,varsinfo2, vars3,varsinfo3, vars4,varsinfo4, vars5,varsinfo5,
		//	oper, operinfo, phase, severity, accuracy, maturity, tag, details
		//	FROM %s WHERE alias='%s'`,
		//	tablename, para.Alias)
		qslice_info = fmt.Sprintf(`SELECT id, alias, varset,
			oper, operinfo, tfunc, phase, severity, accuracy, maturity, tag, details	
			FROM %s WHERE alias='%s'`,
			tablename, para.Alias)
		qslice = append(qslice, qslice_info)
	case "byid":
		qslice_info = fmt.Sprintf(`SELECT id, alias, varset,
			oper, operinfo, tfunc, phase, severity, accuracy, maturity, tag, details	
			FROM %s WHERE id=%d`,
			tablename, para.Id)
		qslice = append(qslice, qslice_info)
	case "byidset":
		qslice_info = fmt.Sprintf(`SELECT id, alias, varset,
			oper, operinfo, tfunc, phase, severity, accuracy, maturity, tag, details	
			FROM %s WHERE id IN (%s)`,
			tablename, para.RuleSet)
		qslice = append(qslice, qslice_info)
	case "aliasonly":
		qslice_info = fmt.Sprintf(`SELECT id, alias	
			FROM %s`,
			tablename)
		qslice = append(qslice, qslice_info)
	default:
		return qslice
	}

	qslice = append(qslice, " order by id desc ")

	/*if para.LastCount != 0 {
		temp_LC := fmt.Sprintf(" limit %d", para.LastCount)
		qslice = append(qslice, temp_LC)
	} else */if para.Count != 0 {
		para.Page = para.Page - 1
		temp_PC := fmt.Sprintf(" limit %d,%d", para.Page*para.Count, para.Count)
		qslice = append(qslice, temp_PC)
	}
	qslice = append(qslice, ";")

	fmt.Println("DBY GetRuleSdSearchMysqlCmd mysql cmd is ", qslice)
	return qslice
}

func GetRuleSdLstCountOnCondition(tablename string, para *TblRuleSdSearchPara) int64 {
	qslice := make([]string, 0)
	var totalCnt int64
	var qslice_info string

	//fmt.Println("GetRuleSdLstCountOnCondition , input para is ", para)

	switch para.Type {
	case "all":
		qslice_info = fmt.Sprintf(`SELECT count(*)
			FROM %s`,
			tablename)
		qslice = append(qslice, qslice_info)
	case "byalias":
		qslice_info = fmt.Sprintf(`SELECT count(*)
			FROM %s WHERE alias='%s'`,
			tablename, para.Alias)
		qslice = append(qslice, qslice_info)
	case "byid":
		qslice_info = fmt.Sprintf(`SELECT count(*)
			FROM %s WHERE id=%d`,
			tablename, para.Id)
		qslice = append(qslice, qslice_info)
	case "aliasonly":
		qslice_info = fmt.Sprintf(`SELECT count(*)
			FROM %s WHERE alias='%s'`,
			tablename, para.Alias)
		qslice = append(qslice, qslice_info)
	//case "exact":
	default:
		fmt.Println("GetRuleSdLstCountOnCondition , input para cmd error : ", para.Type)
		return 0
	}

	////////////////////////////////////////////////////
	//qslice_count := fmt.Sprintf(`select count(*) as totalnum from %s`,
	//	tablename)
	//qslice = append(qslice, qslice_count)
	qslice = append(qslice, ";")
	query := strings.Join(qslice, "")

	rows, err := db.DB.Query(query)
	if err != nil {
		return 0
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&totalCnt)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("GetRuleSdLstCountOnCondition count mysql cmd:", query)
	return int64(totalCnt)
}

func GetRulesConfText(para *TblOLASearchPara) (RuleBytes []byte, errOutput error) {
	var rulesString string
	qslice := make([]string, 0)
	searchPara := TblRuleSdSearchPara{}
	//var list *TblRuleSdLstData

	fmt.Println("GetRuleConfText , input para is ", para)

	//for idx, _ := range para.RuleSets {
	//	if len(para.RuleSets[idx].Rule) > 0 {
	//singleRStr, err1 := GetSingleRuleConfText(&(para.RuleSets[idx]))
	//		if err1 != nil || singleRStr == "" {
	//			return RuleBytes, err1
	//		}
	//		qslice = append(qslice, singleRStr)
	//		qslice = append(qslice, "\n")
	//	}
	//}

	//searchPara.Id = Rule.Id
	//searchPara.Alias = Rule.Rule
	searchPara.RuleSet = para.RuleSet
	searchPara.Count = 100000
	searchPara.Page = 1
	searchPara.Type = "byidset" /*"byalias"*/

	err, list := GetRuleSdLst(&searchPara)
	if err != nil || list == nil {
		//這個異常沒有人 recover
		//panic(fmt.Sprintf("GetRuleSdLst error:%s", err.Error()))
		fmt.Println("GetSingleRuleConfText, call GetRuleSdLst failed, err is ", err.Error())
		rspMsg := fmt.Sprintf(`规则 %s 不存在。 %s 。`, para.RuleSet, err.Error())
		err = errors.New(rspMsg)
		return RuleBytes, err
	}

	for idx, _ := range list.Elements {
		singleRuleStr, err1 := ComposeSingRuleInStr(&(list.Elements[idx]))
		if err1 != nil || singleRuleStr == "" {
			//這個異常沒有人 recover
			//panic(fmt.Sprintf("GetRuleSdLst error:%s", err1.Error()))
			fmt.Println("GetSingleRuleConfText, call ComposeSingRuleInStr failed, err is ", err1.Error())
			rspMsg := fmt.Sprintf(`规则 % 建立失败。 %s 。`, err1.Error())
			err1 = errors.New(rspMsg)
			return RuleBytes, err1
		}

		qslice = append(qslice, singleRuleStr)
		qslice = append(qslice, "\n")
	}

	//////////////////////////////////////
	rulesString = strings.Join(qslice, "")
	if len(rulesString) == 0 {
		rspMsg := fmt.Sprintf(`规则 %s 不存在，无法建立任务 。`, searchPara.RuleSet)
		err1 := errors.New(rspMsg)
		return RuleBytes, err1
	}

	rulesBytes := []byte(rulesString)

	fmt.Println("DBY GetRuleConfText rulesString is :", rulesString)
	fmt.Println("DBY GetRuleConfText rulesBytes is :", rulesBytes)

	return rulesBytes, nil
}

func GetSingleRuleConfText(Rule *TblRuleSdSet) (outStr string, err error) {
	searchPara := TblRuleSdSearchPara{}
	var list *TblRuleSdLstData
	//var singleRuleStr string

	fmt.Println("GetSingleRuleConfText , input Rule is ", Rule)

	searchPara.Id = Rule.Id
	searchPara.Alias = Rule.Rule
	searchPara.Count = 1
	searchPara.Page = 1
	searchPara.Type = "byid" /*"byalias"*/

	err, list = GetRuleSdLst(&searchPara)
	if err != nil || list == nil {
		//這個異常沒有人 recover
		//panic(fmt.Sprintf("GetRuleSdLst error:%s", err.Error()))
		fmt.Println("GetSingleRuleConfText, call GetRuleSdLst failed, err is ", err.Error())
		rspMsg := fmt.Sprintf(`规则 % 不存在。 %s 。`, err.Error())
		err = errors.New(rspMsg)
		return "", err
	}

	fmt.Println("GetSingleRuleConfText, found rule is ", list)

	if list.Counts == 0 {
		return outStr, nil
	}

	singleRuleStr, err1 := ComposeSingRuleString(&(list.Elements[0]))
	if err1 != nil || singleRuleStr == "" {
		//這個異常沒有人 recover
		//panic(fmt.Sprintf("GetRuleSdLst error:%s", err1.Error()))
		fmt.Println("GetSingleRuleConfText, call ComposeSingRuleString failed, err is ", err1.Error())
		rspMsg := fmt.Sprintf(`规则 % 建立失败。 %s 。`, err1.Error())
		err1 = errors.New(rspMsg)
		return "", err1
	}

	return singleRuleStr, nil
}

func ComposeSingRuleInStr(rule *TblRuleSdLstContent /*TblRuleSdLstPublic*/) (outStr string, err error) {
	// 需求实现的依赖：如果一条 rule 里的变量集有 GEO 变量，那么这条rule就只能包含这一个变量 GEO 。
	// 该依赖在UI 输入界面进行控制。

	VarSetBytes, errVS := base64.StdEncoding.DecodeString(rule.VarSet)
	if nil != errVS {
		rspMsg := fmt.Sprintf(`规则 VarSet 解码错误, alias=%s 。 ComposeSingRuleInStr Err is %s 。`,
			rule.Alias, errVS.Error())
		errEx := errors.New(rspMsg)
		return outStr, errEx
	}
	fmt.Println("Dby  --- ComposeSingRuleInStr, VarSetBytes is ", VarSetBytes)
	fmt.Println("Dby  --- ComposeSingRuleInStr, VarSetBytes string is ", string(VarSetBytes))

	switch strings.Contains(string(VarSetBytes), "GEO") {
	case true:
		singleRuleStr, err1 := ComposeSingRuleStringGEO(rule)
		if err1 != nil || singleRuleStr == "" {
			//這個異常沒有人 recover
			//panic(fmt.Sprintf("GetRuleSdLst error:%s", err1.Error()))
			fmt.Println("ComposeSingRuleInStr, call ComposeSingRuleStringGEO failed, err is ", err1.Error())
			rspMsg := fmt.Sprintf(`规则 % 操作失败。ComposeSingRuleStringGEO %s 。`, err1.Error())
			err1 := errors.New(rspMsg)
			return singleRuleStr, err1
		}
		return singleRuleStr, nil
	case false:
		singleRuleStr, err1 := ComposeSingRuleString(rule)
		if err1 != nil || singleRuleStr == "" {
			//這個異常沒有人 recover
			//panic(fmt.Sprintf("GetRuleSdLst error:%s", err1.Error()))
			fmt.Println("ComposeSingRuleInStr, call ComposeSingRuleString failed, err is ", err1.Error())
			rspMsg := fmt.Sprintf(`规则 % 操作失败。ComposeSingRuleString  %s 。`, err1.Error())
			err1 := errors.New(rspMsg)
			return singleRuleStr, err1
		}
		return singleRuleStr, nil
	default:
		rspMsg := fmt.Sprintf(`规则 % 操作失败。 ComposeSingRuleInStr, strings.Contains err 。`, rule.Alias)
		err1 := errors.New(rspMsg)
		return outStr, err1
	}
}

func ComposeSingRuleString(rule *TblRuleSdLstContent /*TblRuleSdLstPublic*/) (outStr string, err error) {
	rslice := make([]string, 0)
	//var flag int = 0
	//var var_varinfo string

	fmt.Println("ComposeSingRuleString , input rule is ", rule)

	// SecRule
	rslice = append(rslice, "SecRule ")

	// set VarSets, if more than 1, should be as REQ:p|REQHD|HAHA, use "|" to separate
	//for idx, _ := range rule.VarSet {
	//	if len(rule.VarSet[idx].Var) > 0 {
	//		if flag >= 1 {
	//			rslice = append(rslice, "|")
	//		}
	//		rslice = append(rslice, rule.VarSet[idx].Var)
	//		if len(rule.VarSet[idx].VarInfo) > 0 {
	//			var_varinfo = fmt.Sprintf(`:%s`, rule.VarSet[idx].VarInfo)
	//			rslice = append(rslice, var_varinfo)
	//		}
	//	}
	//	flag++
	//}

	VarSetBytes, errVS := base64.StdEncoding.DecodeString(rule.VarSet)
	if nil != errVS {
		rspMsg := fmt.Sprintf(`规则 VarSet 解码错误, alias=%s 。 Err is %s 。`, rule.Alias, errVS.Error())
		errEx := errors.New(rspMsg)
		return outStr, errEx
	}

	var_varset := fmt.Sprintf(` %s `, string(VarSetBytes) /*rule.VarSet*/)
	rslice = append(rslice, var_varset)

	// set oper and operinfo
	if len(rule.Oper) > 0 {
		OpInfoBytes, errOpIn := base64.StdEncoding.DecodeString(rule.OperInfo)
		if nil != errOpIn {
			rspMsg := fmt.Sprintf(`规则 operinfo 解码错误, alias=%s 。 Err is %s 。`, rule.Alias, errOpIn.Error())
			errEx := errors.New(rspMsg)
			return outStr, errEx
		}

		var_operinfo := fmt.Sprintf(` "%s %s" `, rule.Oper, string(OpInfoBytes) /*rule.OperInfo*/)
		rslice = append(rslice, var_operinfo)
	}

	// decode rule Details
	rDBytes, errRD := base64.StdEncoding.DecodeString(rule.Details)
	if nil != errRD {
		fmt.Println(`规则 rule Details 解码错误, alias=%s 。 Err is %s 。`, rule.Alias, errRD.Error())
		//rspMsg := fmt.Sprintf(`规则 rule Details 解码错误, alias=%s 。 Err is %s 。`, rule.Alias, errRD.Error())
		//errEx := errors.New(rspMsg)
		//return outStr, errEx
	}

	var_part2 := fmt.Sprintf(` "log,block,t:%s,severity:'%s',ver:'OWASP_CRS/3.0.0',maturity:'%s', accuracy:'%s',phase:%s,id:%d,tag:'%s',msg:'%s'"`,
		rule.TransFunc, rule.Severity, rule.Maturity, rule.Accuracy, rule.Phase, rule.Id, rule.Tag, string(rDBytes) /*rule.Details*/)

	rslice = append(rslice, var_part2)
	outStr = strings.Join(rslice, "")

	fmt.Println("ComposeSingRuleString , outStr is ", outStr)

	return outStr, nil
}

func ComposeSingRuleStringGEO(rule *TblRuleSdLstContent /*TblRuleSdLstPublic*/) (outStr string, err error) {
	rslice := make([]string, 0)

	fmt.Println("ComposeSingRuleString , input rule is ", rule)

	OpInfoBytes, errOpIn := base64.StdEncoding.DecodeString(rule.OperInfo)
	if nil != errOpIn {
		rspMsg := fmt.Sprintf(`规则 operinfo 解码错误, alias=%s 。 ComposeSingRuleStringGEO Err is %s 。`, rule.Alias, errOpIn.Error())
		errEx := errors.New(rspMsg)
		return outStr, errEx
	}
	OpInfoStr := string(OpInfoBytes)

	VarSetBytes, errVS := base64.StdEncoding.DecodeString(rule.VarSet)
	if nil != errVS {
		rspMsg := fmt.Sprintf(`规则 VarSet 解码错误, alias=%s 。 ComposeSingRuleStringGEO Err is %s 。`, rule.Alias, errVS.Error())
		errEx := errors.New(rspMsg)
		return outStr, errEx
	}
	VarSetStr := string(VarSetBytes /*rule.VarSet*/)

	// Prepare Rule section 1
	ruleStr1 := fmt.Sprintf(`SecRule REMOTE_ADDR "@geoLookup" "chain,id:%d,drop,msg:'Non-%s IP address'"`, rule.Id, OpInfoStr)

	// Prepare Rule section 2
	ruleStr2 := fmt.Sprintf(`SecRule %s "%s %s"`, VarSetStr, rule.Oper, OpInfoStr)

	rslice = append(rslice, ruleStr1)
	rslice = append(rslice, "\n")
	rslice = append(rslice, ruleStr2)

	outStr = strings.Join(rslice, "")

	fmt.Println("ComposeSingRuleStringGEO , outStr is ", outStr)

	return outStr, nil
}
