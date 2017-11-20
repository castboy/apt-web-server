/******统计攻击类型的topN数据******/
package attackTopN

import (
	"apt-web-server_v2/models/db"
	"fmt"
	"strings"
	"time"
)

func (this *AttackCount) TableName(attackType string) string {
	switch attackType {
	case "waf_type",
		"disclosure",
		"ddos",
		"reputation_ip",
		"lfi",
		"sqli",
		"xss",
		"injection_php",
		"generic",
		"rce",
		"protocol",
		"rfi",
		"fixation",
		"scaning":
		return "alert_waf"
	case "abnormal_connection",
		"exceptionalvisit",
		"webshell":
		return "urgencymold"
	default:
		return ""
	}
}

func GetTypeList(attackType string) []string {
	idsList := []string{"privilege_gain",
		"ddos",
		"information_leak",
		"web_attack",
		"application_attack",
		"candc",
		"misc_attack"}
	vdsList := []string{"backdoor",
		"trojan",
		"risktool",
		"spyware",
		"malware",
		"virus",
		"worm",
		"joke",
		"adware",
		"hacktool",
		"exploit"}
	wafList := []string{"disclosure",
		"ddos",
		"reputation_ip",
		"lfi",
		"sqli",
		"xss",
		"injection_php",
		"generic",
		"rce",
		"protocol",
		"rfi",
		"fixation",
		"scaning"}
	urgencyList := []string{"abnormal_connection",
		"exceptionalvisit",
		"webshell",
		"sqli",
		"xss",
		"injection_php",
		"rfi"}
	switch attackType {
	case "ids_type":
		return idsList
	case "vds_type":
		return vdsList
	case "waf_type":
		return wafList
	default:
		return urgencyList
	}
}

func SQLOfGetTopN(para *AttackSearchPara) []string {
	var attack string
	var timeTag, tblName string
	switch para.Type {
	case "ids_type":
		tblName = "alert_ids"
		attack = "byzoro_type"
		timeTag = "time"
	case "vds_type":
		tblName = "alert_vds"
		attack = "local_vtype"
		timeTag = "time"
	}
	qslice := make([]string, 0)
	qslice_topN := fmt.Sprintf(`SELECT %s,COUNT(id) AS count 
	                            FROM %s 
								WHERE %s BETWEEN %d AND %d 
								GROUP BY %s 
								ORDER BY count DESC,%s 
								LIMIT %d;`,
		attack, tblName, timeTag, para.Start, para.End, attack, attack, para.Count)
	qslice_WAFtopN := fmt.Sprintf(`SELECT attack,count from 
								(SELECT attack,count(id) AS count FROM alert_waf 
								WHERE time BETWEEN %d AND %d 
								GROUP BY attack) AS a 
								ORDER BY count DESC,attack LIMIT %d;`,
		para.Start, para.End, para.Count)
	qslice_urgency := fmt.Sprintf(`SELECT attack,count FROM
                                (SELECT attack,COUNT(id) AS count,time 
								FROM alert_waf 
								WHERE time BETWEEN %d AND %d 
								AND attack in('sqli','xss','injection_php','rfi') 
								AND severity in (0,1,2) 
								GROUP BY attack
                                UNION ALL
                                SELECT attack_type as attack,COUNT(id) AS count,time 
								FROM urgencymold 
								WHERE time BETWEEN %d AND %d 
								GROUP BY attack) AS a order by count DESC,attack LIMIT %d;`,
		para.Start, para.End, para.Start, para.End, para.Count)
	switch para.Type {
	case "urgent_type":
		qslice = append(qslice, qslice_urgency)
	case "waf_type":
		qslice = append(qslice, qslice_WAFtopN)
	default:
		qslice = append(qslice, qslice_topN)
	}
	return qslice
}

func SQLOfAttackCount(datemold string, tblName string, attackType string, para *AttackSearchPara) []string {
	var attack string
	switch para.Type {
	case "ids_type":
		tblName = "alert_ids"
		attack = "byzoro_type"
	case "vds_type":
		tblName = "alert_vds"
		attack = "local_vtype"
	case "waf_type":
		tblName = "alert_waf"
		attack = "attack"
	case "urgent_type":
		if tblName == "urgencymold" {
			attack = "attack_type"
		} else {
			attack = "attack"
		}
	}
	qslice := make([]string, 0)
	qslice_count := fmt.Sprintf(`SELECT FROM_UNIXTIME(a.time,'%s') AS times,
	    COUNT(a.id) AS num FROM (SELECT time,id FROM %s WHERE %s IN ('%s') AND time BETWEEN %d AND %d`,
		datemold, tblName, attack, attackType, para.Start, para.End)
	qslice = append(qslice, qslice_count)
	if para.Type == "urgent_type" {
		qslice_uType := fmt.Sprintf(` AND severity IN (0,1,2)`)
		qslice = append(qslice, qslice_uType)
	}
	qslice = append(qslice, ") AS a GROUP BY times;")
	return qslice
}
func GetTopNData(tblName string, datemold string, attackType string, para *AttackSearchPara) (error, []AttackCount) {
	var countDay int64
	var catchStruct AttackCount
	timeC := para.Start
	oneDay := int64(60 * 60 * 24)
	countData := make([]AttackCount, 0)
	getData := SQLOfAttackCount(datemold, tblName, attackType, para)
	query := strings.Join(getData, "")
	fmt.Println("sqlGetData=", getData)
	rows, err := db.DB.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	ugc := new(AttackCount)
	for rows.Next() {
		err = rows.Scan(
			&ugc.Time,
			&ugc.Times)
		if err != nil {
			return err, nil
		}
	CHECKTIME:
		timeC = para.Start + oneDay*countDay
		tim := time.Unix(timeC, 0)
		if ugc.Time == tim.Format("01-02") {
			countData = append(countData, *ugc)
			countDay++
		} else {
			catchStruct.Time = tim.Format("01-02")
			countData = append(countData, catchStruct)
			if timeC <= para.End {
				countDay++
				goto CHECKTIME
			}
		}
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	timeC = para.Start + oneDay*countDay
	for timeC <= para.End {
		timO := time.Unix(timeC, 0)
		catchStruct.Time = timO.Format("01-02")
		countData = append(countData, catchStruct)
		countDay++
		timeC = para.Start + oneDay*countDay
	}
	return nil, countData
}
func (this *AttackCount) GetAttackTopN(para *AttackSearchPara) (error, *AttackData) {
	var datemold string
	switch para.Unit {
	case "day":
		datemold = "%m-%d"
	case "month":
		datemold = "%Y-%m"
	case "hour":
		datemold = "%Y-%m-%d %H"
	case "minute":
		datemold = "%Y-%m-%d %H-%i"
	default:
		datemold = "%m-%d"
	}

	list := AttackData{}
	getTopN := SQLOfGetTopN(para)
	query := strings.Join(getTopN, "")
	fmt.Println("sqlGetTopN=", getTopN)
	rows, err := db.DB.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	ugc := new(Attack)
	for rows.Next() {
		err = rows.Scan(
			&ugc.Type,
			&ugc.Total)
		if err != nil {
			ugc.Type = ""
			//return err, nil
		}
		if ugc.Type != "" {
			_, ugc.Data = GetTopNData(this.TableName(ugc.Type), datemold, ugc.Type, para)
			ugc.Type = strings.ToLower(ugc.Type)
			list.Elements = append(list.Elements, AttackList{ugc.AttackContent})
			list.Counts++
		}
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	return nil, &list
}
