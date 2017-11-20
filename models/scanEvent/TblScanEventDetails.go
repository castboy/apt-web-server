/********获取扫描事件详情********/
package scanEvent

import (
	"apt-web-server_v2/models/db"
	"apt-web-server_v2/models/modelsPublic"
	"fmt"
	"strings"
)

func (this *TblSED) TableName() string {
	return "alert_portscan"
}

func GetSEDMysqlCmd(tablename string, para *TblSEDSearchPara) []string {
	//qslice, whereflag := modelsPublic.DefaultParaCmd("getlist", tablename, &para.PField)
	var whereflag int
	qslice := make([]string, 0)
	qslice_tmp := fmt.Sprintf(`SELECT time,conntype,host,scantype,count FROM %s`,
		tablename)
	qslice = append(qslice, qslice_tmp)
	if para.PField.Start != 0 && para.PField.End != 0 {
		qslice_time := fmt.Sprintf(" WHERE time BETWEEN %d AND %d",
			para.PField.Start, para.PField.End)
		qslice = append(qslice, qslice_time)
		whereflag = 1
	}
	if para.Type != "" {
		if whereflag == 1 {
			temp_t := fmt.Sprintf(" AND conntype IN ('%s')", para.Type)
			qslice = append(qslice, temp_t)
		} else {
			temp_t := fmt.Sprintf(" WHERE conntype IN ('%s')", para.Type)
			qslice = append(qslice, temp_t)
			whereflag = 1
		}
	}
	if para.Sort != "" {
		temp_s := fmt.Sprintf(" ORDER BY %s %s", para.Sort, para.Order)
		qslice = append(qslice, temp_s)
	}
	if para.LastCount != 0 {
		temp_LC := fmt.Sprintf(" LIMIT %d", para.LastCount)
		qslice = append(qslice, temp_LC)
	} else {
		para.Page = para.Page - 1
		temp_PC := fmt.Sprintf(" LIMIT %d,%d", para.Page*para.Count, para.Count)
		qslice = append(qslice, temp_PC)
	}
	qslice = append(qslice, ";")
	fmt.Println(tablename, para)
	fmt.Println(qslice)
	return qslice
}

func (this *TblSED) GetScanEventDetails(para *TblSEDSearchPara) (error, *TblSEDData) {
	qslice := GetSEDMysqlCmd(this.TableName(), para)
	query := strings.Join(qslice, "")
	rows, err := db.DB.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	list := TblSEDData{}
	for rows.Next() {
		ugc := new(TblSED)
		err = rows.Scan(
			&ugc.Time,
			&ugc.Conntype,
			&ugc.Host,
			&ugc.AlertType,
			&ugc.Count)
		if err != nil {
			return err, nil
		}
		list.Elements = append(list.Elements, TblSEDList{ugc.TblSEDContent})
		list.Counts++
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	//list.Counts = GetScanEventCounts("", para)
	list.Totality = GetScanEventCounts(this.TableName(), para)

	return nil, &list
}

func GetScanEventCounts(tablename string, para *TblSEDSearchPara) int64 {
	qslice, whereflag := modelsPublic.DefaultParaCmd("getcounts", tablename, &para.PField)
	if tablename != "" {
		if para.Type != "" {
			if whereflag != 0 {
				temp_t := fmt.Sprintf(" AND conntype IN ('%s')", para.Type)
				qslice = append(qslice, temp_t)
			} else {
				temp_t := fmt.Sprintf(" WHERE conntype IN ('%s')", para.Type)
				qslice = append(qslice, temp_t)
			}
		}
	}
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
