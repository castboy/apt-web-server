package models

import (
	"fmt"
	"strings"
)

func (this *TblNetFlow) TableName(unitType string) string {
	switch unitType {
	case "quarter":
		return "tbl_netflow_quarter"
	case "hour":
		return "tbl_netflow_hour"
	case "day":
		return "tbl_netflow_day"
	default:
		return "tbl_netflow_quarter"
	}

}

/*add func*/
func (this *TblNetFlowIP) TableName(unitType string) string {
	switch unitType {
	case "quarter":
		return "tbl_netflowip_quarter"
	case "hour":
		return "tbl_netflowip_hour"
	case "day":
		return "tbl_netflowip_day"
	default:
		return "tbl_netflowip"
	}
	//return "tbl_netflowip"
}
func (this *TblNetFlowD) TableName(unitType string) string {
	switch unitType {
	case "quarter":
		return "tbl_netflowd_quarter"
	case "hour":
		return "tbl_netflowd_hour"
	case "day":
		return "tbl_netflowd_day"
	default:
		return "tbl_netflowd"
	}
	//	return "tbl_netflowd"
}

func (this *TblNetFlowP) TableName(unitType string) string {
	switch unitType {
	case "quarter":
		return "tbl_netflowp_quarter"
	case "hour":
		return "tbl_netflowp_hour"
	case "day":
		return "tbl_netflowp_day"
	default:
		return "tbl_netflowp"
	}
	//	return "tbl_netflowp"
}
func (this *TblNetFlowIPD) TableName(unitType string) string {
	switch unitType {
	case "quarter":
		return "tbl_netflowipd_quarter"
	case "hour":
		return "tbl_netflowipd_hour"
	case "day":
		return "tbl_netflowipd_day"
	default:
		return "tbl_netflowipd"
	}
	//	return "tbl_netflowipd"
}
func (this *TblNetFlowIPP) TableName(unitType string) string {
	switch unitType {
	case "quarter":
		return "tbl_netflowipp_quarter"
	case "hour":
		return "tbl_netflowipp_hour"
	case "day":
		return "tbl_netflowipp_day"
	default:
		return "tbl_netflowipp"
	}
	//	return "tbl_netflowipp"
}
func (this *TblNetFlowDP) TableName(unitType string) string {
	switch unitType {
	case "quarter":
		return "tbl_netflow_quarterdp"
	case "hour":
		return "tbl_netflow_hourdp"
	case "day":
		return "tbl_netflow_daydp"
	default:
		return "tbl_netflowdp"
	}
	//	return "tbl_netflowdp"
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
func GetMysqlSlice(tablename string, para *TblNetFlowSearchPara) []string {
	//	paraDef := TblPublicPara{
	//		Start: para.Start,
	//		End:   para.End}
	qslice, whereflag := DefaultParaCmd("getlist", tablename, &para.PField)

	if para.Protocol != "" {
		if whereflag == 1 {
			temp_p := fmt.Sprintf(" and protocol='%s'", para.Protocol)
			qslice = append(qslice, temp_p)
		} else {
			temp_p := fmt.Sprintf(" where protocol='%s'", para.Protocol)
			qslice = append(qslice, temp_p)
			whereflag = 1
		}
	}
	if para.AssetIP != "" {
		if whereflag == 1 {
			temp_a := fmt.Sprintf(" and assetIP='%s'", para.AssetIP)
			qslice = append(qslice, temp_a)
		} else {
			temp_a := fmt.Sprintf(" where assetIP='%s'", para.AssetIP)
			qslice = append(qslice, temp_a)
			whereflag = 1
		}
	}
	if para.Direction != "" {
		if whereflag == 1 {
			temp_d := fmt.Sprintf(" and direction='%s'", para.Direction)
			qslice = append(qslice, temp_d)
		} else {
			temp_d := fmt.Sprintf(" where direction='%s'", para.Direction)
			qslice = append(qslice, temp_d)
			whereflag = 1
		}
	}
	qslice = append(qslice, " group by time order by time;")
	//	fmt.Println(tablename, para)
	return qslice
}

/*
//func (this *TblNetFlow) GetNetFlow(para *TblNetFlowSearchPara) (error, []TblNetFlowList) {
func (this *TblNetFlow) GetNetFlow(para *TblNetFlowSearchPara) (error, *TblNetFlowData) {
	var wg sync.WaitGroup
	seconds := GetSeconds(para.Unit)
	var start, end int64
	start = para.PField.Start
	end = para.PField.End
	count := (end - start) / seconds
	fmt.Println(start, end, para.PField.Start, para.PField.End, count)
	list := TblNetFlowData{}
	var i int64
	ch := make(chan TblNetFlow, count+1)
	for i = 0; i <= count; i++ {
		wg.Add(1)
		para.PField.Start = start + (i * seconds)
		para.PField.End = para.PField.Start + seconds
		if para.PField.End > end {
			para.PField.End = end
		}
		netFlow := TblNetFlow{}
		netFlow.Time = para.PField.End
		//qslice := GetMysqlSlice(this.TableName(), para)
		query := fmt.Sprintf(`select sum(flow) from %s where time > %d and time < %d;`,
			this.TableName(),
			para.PField.Start,
			para.PField.End)
		//query := strings.Join(qslice, "")
		go func() {
			defer wg.Done()
			rows, err := db.Query(query)
			if err != nil {
				//return err, nil
			}
			defer rows.Close()
			for rows.Next() {
				ugc := new(TblNetFlow)
				err = rows.Scan(
					//&ugc.Id,
					//&ugc.Time,
					&ugc.Flow)
				//&ugc.AssetIP,
				//&ugc.Direction,
				//&ugc.Protocol)
				if err != nil {
					//	return err, nil
				}
				netFlow.Flow = ugc.Flow
			}
			if err := rows.Err(); err != nil {
				//	return err, nil
			}
			ch <- netFlow
		}()

	}

	list.Counts = i
	for flow := range ch {
		list.Elements = append(list.Elements, TblNetFlowList{flow.TblNetFlowContent})
		i--
		if i == 0 {
			close(ch)
		}
	}
	wg.Wait()
	//list.Counts = GetNetFlowCounts("", para)
	//list.Totality = GetNetFlowCounts(this.TableName(), para)
	return nil, &list
}
*/

func (this *TblNetFlow) GetNetFlow(para *TblNetFlowSearchPara) (error, *TblNetFlowData) {
	qslice := GetMysqlSlice(this.TableName(para.Unit), para)
	query := strings.Join(qslice, "")
	rows, err := db.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	//	list := []TblNetFlowList{}
	list := TblNetFlowData{}
	for rows.Next() {
		ugc := new(TblNetFlowCount)
		err = rows.Scan(
			&ugc.Id,
			&ugc.Time,
			&ugc.Flow)
		if err != nil {
			return err, nil
		}
		list.Elements = append(list.Elements, TblNetFlowList{ugc.TblNetFlowContent})
		//		list = append(list, TblNetFlowList{ugc.TblNetFlowContent})
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	list.Counts = GetNetFlowCounts("", para)
	list.Totality = GetNetFlowCounts(this.TableName(para.Unit), para)
	return nil, &list
}

func (this *TblNetFlowIP) GetNetFlowIP(para *TblNetFlowSearchPara) (error, *TblNetFlowData) {
	qslice := GetMysqlSlice(this.TableName(para.Unit), para)
	query := strings.Join(qslice, "")
	rows, err := db.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()

	list := TblNetFlowData{}
	for rows.Next() {
		ugc := new(TblNetFlowIP)
		err = rows.Scan(
			&ugc.Id,
			&ugc.Time,
			&ugc.Flow,
			&ugc.AssetIP)
		if err != nil {
			return err, nil
		}
		list.Elements = append(list.Elements, TblNetFlowList{ugc.TblNetFlowContent})
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	list.Counts = GetNetFlowCounts("", para)
	list.Totality = GetNetFlowCounts(this.TableName(para.Unit), para)
	return nil, &list
}
func (this *TblNetFlowD) GetNetFlowD(para *TblNetFlowSearchPara) (error, *TblNetFlowData) {
	qslice := GetMysqlSlice(this.TableName(para.Unit), para)
	query := strings.Join(qslice, "")
	rows, err := db.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()

	list := TblNetFlowData{}
	for rows.Next() {
		ugc := new(TblNetFlowD)
		err = rows.Scan(
			&ugc.Id,
			&ugc.Time,
			&ugc.Flow,
			&ugc.Direction)
		if err != nil {
			return err, nil
		}
		list.Elements = append(list.Elements, TblNetFlowList{ugc.TblNetFlowContent})
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	list.Counts = GetNetFlowCounts("", para)
	list.Totality = GetNetFlowCounts(this.TableName(para.Unit), para)
	return nil, &list
}

func (this *TblNetFlowP) GetNetFlowP(para *TblNetFlowSearchPara) (error, *TblNetFlowData) {
	qslice := GetMysqlSlice(this.TableName(para.Unit), para)
	query := strings.Join(qslice, "")
	rows, err := db.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()

	list := TblNetFlowData{}
	for rows.Next() {
		ugc := new(TblNetFlowP)
		err = rows.Scan(
			&ugc.Id,
			&ugc.Time,
			&ugc.Flow,
			&ugc.Protocol)
		if err != nil {
			return err, nil
		}
		list.Elements = append(list.Elements, TblNetFlowList{ugc.TblNetFlowContent})
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	list.Counts = GetNetFlowCounts("", para)
	list.Totality = GetNetFlowCounts(this.TableName(para.Unit), para)
	return nil, &list
}
func (this *TblNetFlowIPD) GetNetFlowIPD(para *TblNetFlowSearchPara) (error, *TblNetFlowData) {
	qslice := GetMysqlSlice(this.TableName(para.Unit), para)
	query := strings.Join(qslice, "")
	rows, err := db.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()

	list := TblNetFlowData{}
	for rows.Next() {
		ugc := new(TblNetFlowIPD)
		err = rows.Scan(
			&ugc.Id,
			&ugc.Time,
			&ugc.Flow,
			&ugc.AssetIP,
			&ugc.Direction)
		if err != nil {
			return err, nil
		}
		list.Elements = append(list.Elements, TblNetFlowList{ugc.TblNetFlowContent})
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	list.Counts = GetNetFlowCounts("", para)
	list.Totality = GetNetFlowCounts(this.TableName(para.Unit), para)
	return nil, &list
}
func (this *TblNetFlowIPP) GetNetFlowIPP(para *TblNetFlowSearchPara) (error, *TblNetFlowData) {
	qslice := GetMysqlSlice(this.TableName(para.Unit), para)
	query := strings.Join(qslice, "")
	rows, err := db.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()

	list := TblNetFlowData{}
	for rows.Next() {
		ugc := new(TblNetFlowIPP)
		err = rows.Scan(
			&ugc.Id,
			&ugc.Time,
			&ugc.Flow,
			&ugc.AssetIP,
			&ugc.Protocol)
		if err != nil {
			return err, nil
		}
		list.Elements = append(list.Elements, TblNetFlowList{ugc.TblNetFlowContent})
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	list.Counts = GetNetFlowCounts("", para)
	list.Totality = GetNetFlowCounts(this.TableName(para.Unit), para)
	return nil, &list
}
func (this *TblNetFlowDP) GetNetFlowDP(para *TblNetFlowSearchPara) (error, *TblNetFlowData) {
	qslice := GetMysqlSlice(this.TableName(para.Unit), para)
	query := strings.Join(qslice, "")
	rows, err := db.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()

	list := TblNetFlowData{}
	for rows.Next() {
		ugc := new(TblNetFlowDP)
		err = rows.Scan(
			&ugc.Id,
			&ugc.Time,
			&ugc.Flow,
			&ugc.Direction,
			&ugc.Protocol)
		if err != nil {
			return err, nil
		}
		list.Elements = append(list.Elements, TblNetFlowList{ugc.TblNetFlowContent})
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	list.Counts = GetNetFlowCounts("", para)
	list.Totality = GetNetFlowCounts(this.TableName(para.Unit), para)
	return nil, &list
}

func GetNetFlowCounts(tablename string, para *TblNetFlowSearchPara) int64 {
	//	paraDef := TblPublicPara{
	//		Start: para.Start,
	//		End:   para.End}
	qslice, whereflag := DefaultParaCmd("getcounts", tablename, &para.PField)
	if tablename != "" {
		if para.AssetIP != "" {
			if whereflag != 0 {
				temp_a := fmt.Sprintf(" and assetIP='%s'", para.AssetIP)
				qslice = append(qslice, temp_a)
			} else {
				temp_a := fmt.Sprintf(" where assetIP='%s'", para.AssetIP)
				qslice = append(qslice, temp_a)
			}

		}
		if para.Protocol != "" {
			if whereflag != 0 {
				temp_p := fmt.Sprintf(" and protocol='%s'", para.Protocol)
				qslice = append(qslice, temp_p)
			} else {
				temp_p := fmt.Sprintf(" where protocol='%s'", para.Protocol)
				qslice = append(qslice, temp_p)
			}

		}
		if para.Direction != "" {
			if whereflag != 0 {
				temp_d := fmt.Sprintf(" and direction='%s'", para.Direction)
				qslice = append(qslice, temp_d)
			} else {
				temp_d := fmt.Sprintf(" where direction='%s'", para.Direction)
				qslice = append(qslice, temp_d)
			}

		}

	}
	qslice = append(qslice, ";")
	query := strings.Join(qslice, "")
	fmt.Println(query)
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

func (this *TblNetFlow) CreateSql() string {
	return fmt.Sprintf(
		`CREATE TABLE %s (
		id   integer unsigned  AUTO_INCREMENT NOT NULL,
		time   BIGINT NOT NULL DEFAULT 0,
		flow   BIGINT NOT NULL DEFAULT 0,
		direction varchar(20) NOT NULL DEFAULT '',
		assetIP varchar(20) NOT NULL DEFAULT '',
		protocol varchar(50) NOT NULL DEFAULT '',
		PRIMARY KEY (Id)
	)ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
		this.TableName(""))
}
