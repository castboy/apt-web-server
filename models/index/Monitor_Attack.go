// Monitor_Attack.go
package index

import (
	//"database/sql"
	"apt-web-server_v2/models/modelsPublic"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func map_count(input string) MonitorAttack_st {
	var output MonitorAttack_st
	var info MapCount_st
	//var attackcity []CityName_st
	var dstname CityName_st
	var sql_cmd string

	now_t := time.Now().Unix()
	tmp_t := (now_t + 8*3600) % 86400
	bjzero_t := now_t - tmp_t

	if input == "src_province" {
		sql_cmd = fmt.Sprintf(`select %s,dest_province,count(%s) from ((select %s,
		dest_province from alert_ids where time >= %d and src_country = '中国' and 
		dest_country = '中国' order by id desc limit 20) union all (select %s,
		dest_province from alert_vds where time >= %d and src_country = '中国' and 
		dest_country = '中国' order by id desc limit 20) union all (select %s,
		dest_province from alert_waf where time >= %d and src_country = '中国' and 
		dest_country = '中国' order by id desc limit 20)) as c where src_province 
		<> '中国' group by %s,dest_province`, input, input, input, bjzero_t,
			input, bjzero_t, input, bjzero_t, input)
	} else {
		sql_cmd = fmt.Sprintf(`select %s,dest_country,count(%s) from ((select %s,
		dest_country from alert_ids where time >= %d order by id desc limit 20) 
		union all (select %s,dest_country from alert_vds where time >= %d order 
		by id desc limit 20) union all (select %s,dest_country from alert_waf 
		where time >= %d order by id desc limit 20)) as c group by %s,dest_country`,
			input, input, input, bjzero_t, input, bjzero_t, input, bjzero_t, input)
	}
	//fmt.Println(sql_cmd)
	rows := modelsPublic.Select_mysql(sql_cmd)
	if rows == nil {
		return output
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&(info.Area), &(dstname.Name), &(info.Count)); err == nil {
			var attackcity []CityName_st
			var srcname CityName_st
			if info.Area == "localhost" {
				info.Area = "北京"
			}
			output.MapCount = append(output.MapCount, info)
			srcname.Name = info.Area
			if dstname.Name == "localhost" {
				dstname.Name = "北京"
			}
			attackcity = append(attackcity, srcname)
			attackcity = append(attackcity, dstname)

			output.Attackcity = append(output.Attackcity, attackcity)
			output.Count = output.Count + 1
		}
	}
	return output
}
func attackcount_day(input string) MonitorAttack_st {
	var output MonitorAttack_st
	var info AttackDay_st
	var sql_cmd string
	var source string
	var attack_count int
	sql_cmd = fmt.Sprintf(`select source,sum(attack_count) from attack_days 
		group by %s`, input)
	//fmt.Println(sql_cmd)
	rows := modelsPublic.Select_mysql(sql_cmd)
	if rows == nil {
		return output
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&source, &attack_count); err == nil {
			switch source {
			case "IDS":
				info.IdsCount = attack_count
			case "VDS":
				info.VdsCount = attack_count
			case "WAF":
				info.WafCount = attack_count
			default:
			}
		}
	}
	info.Count = info.IdsCount + info.VdsCount + info.WafCount
	output.Attack_Day = info
	output.Count = 1
	return output
}
func Monitor_Lookup(input StatisticsAttackIn_st) MonitorAttack_st {
	var output MonitorAttack_st
	var info AttackInfo_st
	var sql_cmd string
	if input.Attribute != "" {
		return output
	}

	now_t := time.Now().Unix()
	tmp_t := (now_t + 8*3600) % 86400
	bjzero_t := now_t - tmp_t

	sql_cmd = fmt.Sprintf(`select id,src_ip,src_province,dest_province,dest_ip,proto,
		dest_port,byzoro_type,time,"bruteforce" as from_type from alert_ids where 
		time >= %d order by id desc limit 20`, bjzero_t)
	//fmt.Println(sql_cmd)
	rows := modelsPublic.Select_mysql(sql_cmd)
	if rows == nil {
		return output
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&(info.MysqlId), &(info.SrcIp), &(info.SrcProvince),
			&(info.DestProvince), &(info.DestIp), &(info.Proto),
			&(info.DestPort), &(info.AttackType), &(info.Time), &(info.FromType)); err == nil {
			if info.SrcProvince == "localhost" {
				info.SrcProvince = "北京"
			}
			if info.DestProvince == "localhost" {
				info.DestProvince = "北京"
			}
			output.Ids_Attack = append(output.Ids_Attack, info)
			output.Count = output.Count + 1
		}
	}
	sql_cmd = fmt.Sprintf(`select id,src_ip,src_province,dest_province,dest_ip,proto,
		dest_port,local_vtype,time,"file" as from_type from alert_vds where time >= %d 
		order by id desc limit 20`, bjzero_t)

	//fmt.Println(sql_cmd)
	rows1 := modelsPublic.Select_mysql(sql_cmd)
	if rows1 == nil {
		return output
	}
	defer rows1.Close()
	for rows1.Next() {
		if err := rows1.Scan(&(info.MysqlId), &(info.SrcIp), &(info.SrcProvince),
			&(info.DestProvince), &(info.DestIp), &(info.Proto),
			&(info.DestPort), &(info.AttackType), &(info.Time), &(info.FromType)); err == nil {
			if info.SrcProvince == "localhost" {
				info.SrcProvince = "北京"
			}
			if info.DestProvince == "localhost" {
				info.DestProvince = "北京"
			}
			output.Vds_Attack = append(output.Vds_Attack, info)
			output.Count = output.Count + 1
		}
	}
	sql_cmd = fmt.Sprintf(`select id,src_ip,src_province,dest_province,dest_ip,proto,
		dest_port,attack,time,"index" as from_type from alert_waf where time >= %d 
		order by id desc limit 20`, bjzero_t)
	//fmt.Println(sql_cmd)
	rows2 := modelsPublic.Select_mysql(sql_cmd)
	if rows2 == nil {
		return output
	}
	defer rows2.Close()
	for rows2.Next() {
		if err := rows2.Scan(&(info.MysqlId), &(info.SrcIp), &(info.SrcProvince),
			&(info.DestProvince), &(info.DestIp), &(info.Proto),
			&(info.DestPort), &(info.AttackType), &(info.Time), &(info.FromType)); err == nil {
			if info.SrcProvince == "localhost" {
				info.SrcProvince = "北京"
			}
			if info.DestProvince == "localhost" {
				info.DestProvince = "北京"
			}
			output.Waf_Attack = append(output.Waf_Attack, info)
			output.Count = output.Count + 1
		}
	}
	return output
}
func Monitor_Count(input StatisticsAttackIn_st) MonitorAttack_st {
	var output MonitorAttack_st
	switch input.Attribute {
	case "country":
		output = map_count("src_country")
	case "province":
		output = map_count("src_province")
	case "source":
		output = attackcount_day("source")
	default:
		return output
	}
	return output
}
