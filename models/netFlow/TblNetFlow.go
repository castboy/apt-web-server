package netFlow

import (
	"apt-web-server_v2/models/db"
	"apt-web-server_v2/models/modelsPublic"
	"fmt"
	"strings"
)

func (this *TblNetFlow) TableName(unitType string) string {
	switch unitType {
	case "minute":
		return "netflow"
	case "quarter":
		return "netflow_q"
	case "hour":
		return "netflow_h"
	case "day":
		return "netflow_d"
	default:
		return "netflow_q"
	}
}

/*add func*/
func (this *TblNetFlowIP) TableName(unitType string) string {
	switch unitType {
	case "minute":
		return "netflowip"
	case "quarter":
		return "netflowip_q"
	case "hour":
		return "netflowip_h"
	case "day":
		return "netflowip_d"
	default:
		return "netflowip_q"
	}
}
func (this *TblNetFlowD) TableName(unitType string) string {
	switch unitType {
	case "minute":
		return "netflowd"
	case "quarter":
		return "netflowd_q"
	case "hour":
		return "netflowd_h"
	case "day":
		return "netflowd_d"
	default:
		return "netflowd_q"
	}
}

func (this *TblNetFlowP) TableName(unitType string) string {
	switch unitType {
	case "minute":
		return "netflowp"
	case "quarter":
		return "netflowp_q"
	case "hour":
		return "netflowp_h"
	case "day":
		return "netflowp_d"
	default:
		return "netflowp_q"
	}
}
func (this *TblNetFlowIPD) TableName(unitType string) string {
	switch unitType {
	case "minute":
		return "netflowipd"
	case "quarter":
		return "netflowipd_q"
	case "hour":
		return "netflowipd_h"
	case "day":
		return "netflowipd_d"
	default:
		return "netflowipd_q"
	}
}
func (this *TblNetFlowIPP) TableName(unitType string) string {
	switch unitType {
	case "minute":
		return "netflowipp"
	case "quarter":
		return "netflowipp_q"
	case "hour":
		return "netflowipp_h"
	case "day":
		return "netflowipp_d"
	default:
		return "netflowipp_q"
	}
}
func (this *TblNetFlowDP) TableName(unitType string) string {
	switch unitType {
	case "minute":
		return "netflowdp"
	case "quarter":
		return "netflowdp_q"
	case "hour":
		return "netflowdp_h"
	case "day":
		return "netflowdp_d"
	default:
		return "netflowdp_q"
	}
}

/*add func end*/

func GetSeconds(unitType string) int64 {
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
	default:
		seconds = (60 * 60)
	}
	return seconds
}
func GetColumn(tablename string) string {
	var column string
	switch tablename {
	case "netflowip", "netflowip_q", "netflowip_h", "netflowip_d":
		column = ",assetIP"
	case "netflowd", "netflowd_q", "netflowd_h", "netflowd_d":
		column = ",direction"
	case "netflowp", "netflowp_q", "netflowp_h", "netflowp_d":
		column = ",protocol"
	case "netflowipd", "netflowipd_q", "netflowipd_h", "netflowipd_d":
		column = ",assetIP,direction"
	case "netflowipp", "netflowipp_q", "netflowipp_h", "netflowipp_d":
		column = ",assetIP,protocol"
	case "netflowdp", "netflowdp_q", "netflowdp_h", "netflowdp_d":
		column = ",direction,protocol"
	}
	return column
}
func GetMysqlSlice(tablename string, para *TblNetFlowSearchPara) []string {
	qslice := make([]string, 0)
	whereflag := 0
	cmd := fmt.Sprintf(`SELECT`)
	qslice = append(qslice, cmd)
	qslice = append(qslice, " time,SUM(flow) as flow")
	//column := GetColumn(tablename)
	//qslice = append(qslice, column)
	fromTable := fmt.Sprintf(` FROM %s`, tablename)
	qslice = append(qslice, fromTable)
	if para.Protocol != "" {
		if whereflag == 1 {
			temp_p := fmt.Sprintf(" AND protocol IN ('%s')", para.Protocol)
			qslice = append(qslice, temp_p)
		} else {
			temp_p := fmt.Sprintf(" WHERE protocol IN ('%s')", para.Protocol)
			qslice = append(qslice, temp_p)
			whereflag = 1
		}
	}
	if para.AssetIP != "" {
		if whereflag == 1 {
			temp_a := fmt.Sprintf(" AND assetIP IN ('%s')", para.AssetIP)
			qslice = append(qslice, temp_a)
		} else {
			temp_a := fmt.Sprintf(" WHERE assetIP IN ('%s')", para.AssetIP)
			qslice = append(qslice, temp_a)
			whereflag = 1
		}
	}
	if para.Direction != "" {
		if whereflag == 1 {
			temp_d := fmt.Sprintf(" AND direction IN ('%s')", para.Direction)
			qslice = append(qslice, temp_d)
		} else {
			temp_d := fmt.Sprintf(" WHERE direction IN ('%s')", para.Direction)
			qslice = append(qslice, temp_d)
			whereflag = 1
		}
	}
	qslice = append(qslice, " GROUP BY time ORDER BY time;")
	//	fmt.Println(tablename, para)
	return qslice
}
func GetNetFlowList(tablename string, para *TblNetFlowSearchPara) (error, *TblNetFlowData) {
	qslice := GetMysqlSlice(tablename, para)
	query := strings.Join(qslice, "")
	rows, err := db.DB.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	//	list := []TblNetFlowList{}
	list := TblNetFlowData{}
	for rows.Next() {
		ugc := new(TblNetFlow)
		err = rows.Scan(
			&ugc.Time,
			&ugc.Flow)
		if err != nil {
			return err, nil
		}
		list.Elements = append(list.Elements, TblNetFlowList{ugc.TblNetFlowContent})
		list.Counts++
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	list.Totality = GetNetFlowCounts(tablename, para)
	return nil, &list
}
func (this *TblNetFlow) GetNetFlow(para *TblNetFlowSearchPara) (error, *TblNetFlowData) {
	tablename := this.TableName(para.Unit)
	err, list := GetNetFlowList(tablename, para)
	return err, list
}

func (this *TblNetFlowIP) GetNetFlowIP(para *TblNetFlowSearchPara) (error, *TblNetFlowData) {
	tablename := this.TableName(para.Unit)
	err, list := GetNetFlowList(tablename, para)
	return err, list
}
func (this *TblNetFlowD) GetNetFlowD(para *TblNetFlowSearchPara) (error, *TblNetFlowData) {
	tablename := this.TableName(para.Unit)
	err, list := GetNetFlowList(tablename, para)
	return err, list
}

func (this *TblNetFlowP) GetNetFlowP(para *TblNetFlowSearchPara) (error, *TblNetFlowData) {
	tablename := this.TableName(para.Unit)
	err, list := GetNetFlowList(tablename, para)
	return err, list
}
func (this *TblNetFlowIPD) GetNetFlowIPD(para *TblNetFlowSearchPara) (error, *TblNetFlowData) {
	tablename := this.TableName(para.Unit)
	err, list := GetNetFlowList(tablename, para)
	return err, list
}
func (this *TblNetFlowIPP) GetNetFlowIPP(para *TblNetFlowSearchPara) (error, *TblNetFlowData) {
	tablename := this.TableName(para.Unit)
	err, list := GetNetFlowList(tablename, para)
	return err, list
}
func (this *TblNetFlowDP) GetNetFlowDP(para *TblNetFlowSearchPara) (error, *TblNetFlowData) {
	tablename := this.TableName(para.Unit)
	err, list := GetNetFlowList(tablename, para)
	return err, list
}

func GetNetFlowCounts(tablename string, para *TblNetFlowSearchPara) int64 {
	qslice, whereflag := modelsPublic.DefaultParaCmd("getcounts", tablename, &para.PField)
	if tablename != "" {
		if para.AssetIP != "" {
			if whereflag != 0 {
				temp_a := fmt.Sprintf(" AND assetIP='%s'", para.AssetIP)
				qslice = append(qslice, temp_a)
			} else {
				temp_a := fmt.Sprintf(" WHERE assetIP='%s'", para.AssetIP)
				qslice = append(qslice, temp_a)
			}

		}
		if para.Protocol != "" {
			if whereflag != 0 {
				temp_p := fmt.Sprintf(" AND protocol='%s'", para.Protocol)
				qslice = append(qslice, temp_p)
			} else {
				temp_p := fmt.Sprintf(" WHERE protocol='%s'", para.Protocol)
				qslice = append(qslice, temp_p)
			}

		}
		if para.Direction != "" {
			if whereflag != 0 {
				temp_d := fmt.Sprintf(" AND direction='%s'", para.Direction)
				qslice = append(qslice, temp_d)
			} else {
				temp_d := fmt.Sprintf(" WHERE direction='%s'", para.Direction)
				qslice = append(qslice, temp_d)
			}

		}

	}
	qslice = append(qslice, ";")
	query := strings.Join(qslice, "")
	fmt.Println(query)
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
