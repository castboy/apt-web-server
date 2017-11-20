/********获取紧急事件数（天）********/
package urgency

import (
	"apt-web-server_v2/models/db"
	"apt-web-server_v2/models/modelsPublic"
	"fmt"
	"strings"
)

func (this *TblUrgency) TableName() string {
	return "tbl_urgency"
}
func GetUgcSeconds(unitType string) int64 {
	var seconds int64
	switch unitType {
	case "minute":
		seconds = 60
	case "quarter":
		seconds = (15 * 60)
	case "hour":
		seconds = (60 * 60)
	case "day":
		seconds = (24 * 60 * 60)
	case "month":
		seconds = (31 * 24 * 60 * 60)
	}
	return seconds
}
func GetUrgencyMysqlCmd(datemold string, para *TblUrgencySearchPara) []string {
	qslice := make([]string, 0)
	qslice_mysql := fmt.Sprintf(`CALL allurgencycount('%s',%d,%d)`, datemold, para.PField.Start, para.PField.End)
	qslice = append(qslice, qslice_mysql)
	qslice = append(qslice, ";")
	fmt.Println(qslice)
	return qslice
}

func (this *TblUrgency) GetUrgencyC(para *TblUrgencySearchPara) (error, *TblUrgencyData) {
	var datemold string
	var count_tmp int64
	switch para.Unit {
	case "day":
		datemold = "%Y-%m-%d"
	case "month":
		datemold = "%Y-%m"
	default:
		datemold = "%Y-%m-%d"
	}
	list := TblUrgencyData{}
	qslice := GetUrgencyMysqlCmd(datemold, para)
	query := strings.Join(qslice, "")
	fmt.Sprintf(query)
	rows, err := db.DB.Query(query)
	ugc := new(TblUrgency)
	for rows.Next() {
		err = rows.Scan(
			&ugc.Date,
			&ugc.Times)
		if err != nil {
			return err, nil
		}
		count_tmp += ugc.Times
		list.Elements = append(list.Elements, TblUrgencyList{ugc.TblUrgencyContent})
	}
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	list.Counts = count_tmp
	return nil, &list
}

func GetUrgencyCounts(tablename string, para *TblUrgencySearchPara) int64 {
	qslice, _ := modelsPublic.DefaultParaCmd("getcounts", tablename, &para.PField)
	qslice = append(qslice, ";")
	query := strings.Join(qslice, "")
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
