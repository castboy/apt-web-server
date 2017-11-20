/********获取恶意流量分类数(天)********/
package ids

import (
	"apt-web-server_v2/models/db"
	"apt-web-server_v2/models/modelsPublic"
	"fmt"
	"strings"
)

func (this *TblMFC) TableName() string {
	return "byzoro_ids_count"
	//return "ids_count_day"
}

func GetMFCMysqlCmd(tablename string, para *TblMFCSearchPara) []string {
	qslice, _ := modelsPublic.DefaultParaCmd("getlist", tablename, &para.PField)
	qslice = append(qslice, ";")
	fmt.Println(tablename, para)
	return qslice
}

func (this *TblMFC) GetUrgencyDetails(para *TblMFCSearchPara) (error, *TblMFCData) {
	qslice := GetMFCMysqlCmd(this.TableName(), para)
	query := strings.Join(qslice, "")
	rows, err := db.DB.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	list := TblMFCData{}
	for rows.Next() {
		ugc := new(TblMFC)
		err = rows.Scan(
			&ugc.Id,
			&ugc.Time,
			&ugc.PrivilegeGain,
			&ugc.DDos,
			&ugc.InformationLeak,
			&ugc.WebAttack,
			&ugc.ApplicationAttack,
			&ugc.CandC,
			&ugc.Malware,
			&ugc.MiscAttack,
			&ugc.Other)
		if err != nil {
			return err, nil
		}
		list.Classify = append(list.Classify, TblMFCList{ugc.TblMFCContent})
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	list.Counts = GetMFCCounts(this.TableName(), para)

	return nil, &list
}

func GetMFCCounts(tablename string, para *TblMFCSearchPara) int64 {
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

func (this *TblMFC) CreateSql() string {
	return fmt.Sprintf(
		`CREATE TABLE %s (
		id   integer unsigned  AUTO_INCREMENT NOT NULL,
		time BIGINT NOT NULL DEFAULT 0,
		attempted_admin BIGINT NOT NULL DEFAULT 0,
		attempted_user BIGINT NOT NULL DEFAULT 0,
		inappropriate_content BIGINT NOT NULL DEFAULT 0,
		policy_violation BIGINT NOT NULL DEFAULT 0,
		shellcode_detect BIGINT NOT NULL DEFAULT 0,
		successful_admin BIGINT NOT NULL DEFAULT 0,
		successful_user BIGINT NOT NULL DEFAULT 0,
		trojan_activity BIGINT NOT NULL DEFAULT 0,
		unsuccessful_user BIGINT NOT NULL DEFAULT 0,
		web_application_attack BIGINT NOT NULL DEFAULT 0,
		attempted_dos BIGINT NOT NULL DEFAULT 0,
		attempted_recon BIGINT NOT NULL DEFAULT 0,
		bad_unknown BIGINT NOT NULL DEFAULT 0,
		default_login_attempt BIGINT NOT NULL DEFAULT 0,
		denial_of_service BIGINT NOT NULL DEFAULT 0,
		misc_attack BIGINT NOT NULL DEFAULT 0,
		non_standard_protocol BIGINT NOT NULL DEFAULT 0,
		rpc_portmap_decode BIGINT NOT NULL DEFAULT 0,
		successful_dos BIGINT NOT NULL DEFAULT 0,
		successful_recon_largescale BIGINT NOT NULL DEFAULT 0,
		successful_recon_limited BIGINT NOT NULL DEFAULT 0,
		suspicious_filename_detect BIGINT NOT NULL DEFAULT 0,
		suspicious_login BIGINT NOT NULL DEFAULT 0,
		system_call_detect BIGINT NOT NULL DEFAULT 0,
		unusual_client_port_connection BIGINT NOT NULL DEFAULT 0,
		web_application_activity BIGINT NOT NULL DEFAULT 0,
		icmp_event BIGINT NOT NULL DEFAULT 0,
		misc_activity BIGINT NOT NULL DEFAULT 0,
		network_scan BIGINT NOT NULL DEFAULT 0,
		not_suspicious BIGINT NOT NULL DEFAULT 0,
		protocol_command_decode BIGINT NOT NULL DEFAULT 0,
		string_detect BIGINT NOT NULL DEFAULT 0,
		unknown BIGINT NOT NULL DEFAULT 0,
		tcp_connection BIGINT NOT NULL DEFAULT 0,
		other BIGINT NOT NULL DEFAULT 0,
		PRIMARY KEY (id)
		)ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
		this.TableName())
}
