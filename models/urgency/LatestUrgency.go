/**********最新紧急事件**********/
package urgency

import (
	"apt-web-server_v2/models/db"
	"apt-web-server_v2/modules/mlog"
	"fmt"
	"strings"
)

func (this *TblUgcL) TableName() string {
	return "urgencymold"
}

func GetUgcLMysqlCmd(tablename string, para *TblUgcLSearchPara) []string {
	//var whereflag int
	qslice := make([]string, 0)
	qslice_tmp := fmt.Sprintf(`SELECT dest_ip,attack_type,time FROM urgencymold 
	    UNION ALL 
	    SELECT client AS dest_ip,attack AS attack_type,time FROM alert_waf 
		WHERE attack IN ('sqli','xss','injection_php','rfi') 
		AND severity IN (0,1,2) 
		ORDER BY time DESC LIMIT %d`, para.LastCount)
	qslice = append(qslice, qslice_tmp)

	qslice = append(qslice, ";")
	fmt.Println(tablename, para)
	return qslice
}

func (this *TblUgcL) GetUrgencyLatest(para *TblUgcLSearchPara) (error, *TblUgcLData) {
	qslice := GetUgcLMysqlCmd(this.TableName(), para)
	query := strings.Join(qslice, "")

	rows, err := db.DB.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	list := TblUgcLData{}
	for rows.Next() {
		ugc := new(TblUgcL)
		err = rows.Scan(
			&ugc.DestIp,
			&ugc.AttackType,
			&ugc.Time)
		if err != nil {
			mlog.Debug(query, err)
			//return err, nil
		}
		list.Elements = append(list.Elements, TblUgcLList{ugc.TblUgcLContent})
		list.Counts++
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	return nil, &list
}
