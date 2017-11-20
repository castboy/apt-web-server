/********White List 查找********/
package whiteList

import (
	"apt-web-server_v2/models/db"
	"fmt"
	"strings"
	//"time"
)

func (this *TblWLLst) GetTableName() string {
	//return "cert_data"
	return "abnconn_whitelist"
}

func GetWLSearchMysqlCmd(datemold, tablename string, para *TblWLSearchPara) []string {
	//var flag int
	var qslice_info string
	qslice := make([]string, 0)

	switch para.Type {
	case "all":
		qslice_info = fmt.Sprintf(`SELECT (CASE WHEN src_ip is NULL THEN '' ELSE src_ip END) as src_ip, 
			src_port, (CASE WHEN dest_ip is NULL THEN '' ELSE dest_ip END) as dest_ip, dest_port, proto
			FROM %s`,
			tablename)
		qslice = append(qslice, qslice_info)
	case "exact":
		qslice_info = fmt.Sprintf(`SELECT (CASE WHEN src_ip is NULL THEN '' ELSE src_ip END) as src_ip, 
			src_port, (CASE WHEN dest_ip is NULL THEN '' ELSE dest_ip END) as dest_ip, dest_port, proto 
			FROM %s WHERE `,
			tablename)

		qslice = append(qslice, qslice_info)

		//set 'where' conditions by one or some of sip sport dip dport proto
		var muls int
		if para.Sip != "" {
			qslice_sip := fmt.Sprintf(`src_ip='%s' `,
				para.Sip)
			qslice = append(qslice, qslice_sip)
			muls = 1
		}
		if para.Sport != 0 {
			if muls == 1 {
				qslice_sport := fmt.Sprintf(`and src_port=%d `,
					para.Sport)
				qslice = append(qslice, qslice_sport)
			} else {
				qslice_sport := fmt.Sprintf(`src_port=%d `,
					para.Sport)
				qslice = append(qslice, qslice_sport)
			}
			muls = 1
		}
		if para.Dip != "" {
			if muls == 1 {
				qslice_dip := fmt.Sprintf(`and dest_ip='%s' `,
					para.Dip)
				qslice = append(qslice, qslice_dip)
			} else {
				qslice_dip := fmt.Sprintf(`dest_ip='%s' `,
					para.Dip)
				qslice = append(qslice, qslice_dip)
			}
			muls = 1
		}
		if para.Dport != 0 {
			if muls == 1 {
				qslice_dport := fmt.Sprintf(`and dest_port=%d `,
					para.Dport)
				qslice = append(qslice, qslice_dport)
			} else {
				qslice_dport := fmt.Sprintf(`dest_port=%d `,
					para.Dport)
				qslice = append(qslice, qslice_dport)
			}
			muls = 1
		}
		if para.Proto != 0 {
			if muls == 1 {
				qslice_pro := fmt.Sprintf(`and proto=%d `,
					para.Proto)
				qslice = append(qslice, qslice_pro)
			} else {
				qslice_pro := fmt.Sprintf(`proto=%d `,
					para.Proto)
				qslice = append(qslice, qslice_pro)
			}
			muls = 1
		}
	default:
		return qslice
	}

	//qslice = append(qslice, qslice_info)
	// order
	qslice = append(qslice, " order by id desc ")

	if para.LastCount != 0 {
		temp_LC := fmt.Sprintf(" limit %d", para.LastCount)
		qslice = append(qslice, temp_LC)
	} else if para.Count != 0 {
		para.Page = para.Page - 1
		temp_PC := fmt.Sprintf(" limit %d,%d", para.Page*para.Count, para.Count)
		qslice = append(qslice, temp_PC)
	}
	qslice = append(qslice, ";")

	fmt.Println("DBY the mysql cmd is ", qslice)
	return qslice
}

func (this *TblWLLst) GetWLLCLst(para *TblWLSearchPara) (error, *TblWLLstData) {
	//datemold := GetDateMold(para.Unit)  datemold = "%Y-%m-%d"
	datemold := "%Y-%m-%d %H:%i:%s"
	var cnt int64 = 0
	list := TblWLLstData{}

	qslice := GetWLSearchMysqlCmd(datemold, this.GetTableName(), para)
	query := strings.Join(qslice, "")

	fmt.Println(qslice)
	fmt.Println(query)

	rows, err := db.DB.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	//ugc := new(TblWLLst)
	var ugc TblWLLst
	var wlinfo TblWLLstContent
	//wlinfo := new(TblWLLstContent /*TblWLLstInfoTmp*/)
	for rows.Next() {
		err = rows.Scan(
			//&ugc.WLId,
			&ugc.WlSip,
			&ugc.WlSport,
			&ugc.WlDip,
			&ugc.WlDport,
			&ugc.WlProto)
		if err != nil {
			return err, nil
		}

		if false == WLCheckProtoNum(ugc.WlProto) {
			fmt.Println("GetWLLCLst white list proto num is illegal. ", ugc.WlProto)
			return err, nil
		}

		wlinfo.WLSip = ugc.WlSip
		wlinfo.WLSport = ugc.WlSport
		wlinfo.WLDip = ugc.WlDip
		wlinfo.WLDport = ugc.WlDport
		if ugc.WlProto == 0 {
			wlinfo.WLProto = ""
		} else {
			wlinfo.WLProto = MapNum2Proto[ugc.WlProto]
		}

		cnt++
		list.Elements = append(list.Elements, wlinfo /*ugc.TblWLLstContent*/)
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	list.Counts = cnt //GetSSLCount("", para)
	if cnt != 0 {
		//list.Totality = GetWLLstCount(this.GetTableName(), para)
		list.Totality = GetWLLstCountOnCondition(this.GetTableName(), para)
	}

	return nil, &list
}

func GetWLLstCount(tablename string, para *TblWLSearchPara) int64 {
	qslice := make([]string, 0)
	var totalCnt int64

	qslice_count := fmt.Sprintf(`select count(*) as totalnum from %s`,
		tablename)
	qslice = append(qslice, qslice_count)
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
	fmt.Println("count mysql cmd:", query)
	return int64(totalCnt)
}

func GetWLLstCountOnCondition(tablename string, para *TblWLSearchPara) int64 {
	qslice := make([]string, 0)
	var totalCnt int64
	////////////////////////////////////////////////////
	var qslice_info string

	//fmt.Println("GetWLSearchMysqlCmd , input para is ", para)

	switch para.Type {
	case "all":
		qslice_info = fmt.Sprintf(`SELECT count(*)
			FROM %s`,
			tablename)
		qslice = append(qslice, qslice_info)
	case "exact":
		qslice_info = fmt.Sprintf(`SELECT count(*) 
			FROM %s WHERE `,
			tablename)

		qslice = append(qslice, qslice_info)

		//set 'where' conditions by one or some of sip sport dip dport proto
		var muls int
		if para.Sip != "" {
			qslice_sip := fmt.Sprintf(`src_ip='%s' `,
				para.Sip)
			qslice = append(qslice, qslice_sip)
			muls = 1
		}
		if para.Sport != 0 {
			if muls == 1 {
				qslice_sport := fmt.Sprintf(`and src_port=%d `,
					para.Sport)
				qslice = append(qslice, qslice_sport)
			} else {
				qslice_sport := fmt.Sprintf(`src_port=%d `,
					para.Sport)
				qslice = append(qslice, qslice_sport)
			}
			muls = 1
		}
		if para.Dip != "" {
			if muls == 1 {
				qslice_dip := fmt.Sprintf(`and dest_ip='%s' `,
					para.Dip)
				qslice = append(qslice, qslice_dip)
			} else {
				qslice_dip := fmt.Sprintf(`dest_ip='%s' `,
					para.Dip)
				qslice = append(qslice, qslice_dip)
			}
			muls = 1
		}
		if para.Dport != 0 {
			if muls == 1 {
				qslice_dport := fmt.Sprintf(`and dest_port=%d `,
					para.Dport)
				qslice = append(qslice, qslice_dport)
			} else {
				qslice_dport := fmt.Sprintf(`dest_port=%d `,
					para.Dport)
				qslice = append(qslice, qslice_dport)
			}
			muls = 1
		}
		if para.Proto != 0 {
			if muls == 1 {
				qslice_pro := fmt.Sprintf(`and proto=%d `,
					para.Proto)
				qslice = append(qslice, qslice_pro)
			} else {
				qslice_pro := fmt.Sprintf(`proto=%d `,
					para.Proto)
				qslice = append(qslice, qslice_pro)
			}
			muls = 1
		}
	default:
		fmt.Println("GetWLLstCountOnCondition , input para cmd error : ", para.Type)
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
	fmt.Println("count mysql cmd:", query)
	return int64(totalCnt)
}

func WLCheckProtoNum(proNum int32) bool {
	if proNum == 0 {
		return true
	}

	if MapNum2Proto[proNum] == "" {
		return false
	} else {
		return true
	}
}
