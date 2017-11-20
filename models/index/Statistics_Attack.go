// Statistics_Attack.go
package index

import (
	//"database/sql"
	"apt-web-server_v2/models/modelsPublic"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func percent_select(table string) StatisticsAttackOut_st {
	var output StatisticsAttackOut_st
	var info Country_st
	var sql_cmd string
	count := 0

	sql_cmd = fmt.Sprintf(`select sum(a.attack_count) from ((select 
				country,sum(attack_count)as attack_count from attack_days 
				where country='中国') union all (select country,sum(attack_count) 
				as attack_count from %s where country='中国')) as a where 
				attack_count <> ''`, table)
	//fmt.Println(sql_cmd)
	rows := modelsPublic.Select_mysql(sql_cmd)
	if rows == nil {
		return output
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&count); err == nil {

		}
	}
	info.Home = count
	count = 0
	sql_cmd = fmt.Sprintf(`select sum(a.attack_count) from ((select 
				country,sum(attack_count)as attack_count from attack_days 
				where country<>'中国') union all (select country,sum(attack_count) 
				as attack_count from %s where country<>'中国')) as a where 
				a.attack_count <> ''`, table)
	rows1 := modelsPublic.Select_mysql(sql_cmd)
	if rows1 == nil {
		return output
	}
	defer rows1.Close()
	for rows1.Next() {
		if err := rows1.Scan(&count); err == nil {

		}
	}
	info.Abroad = count
	output.CountryInfo = append(output.CountryInfo, info)
	output.Count = output.Count + 1
	return output
}
func Statistics_Percent(input StatisticsAttackIn_st) StatisticsAttackOut_st {
	var output StatisticsAttackOut_st
	switch input.Strength {
	case 7:
		output = percent_select("attack_6days")
	case 30:
		output = percent_select("attack_29days")
	default:
		return output
	}
	return output
}
func Statistics_Count(input StatisticsAttackIn_st) StatisticsAttackOut_st {
	var output StatisticsAttackOut_st
	var info Info_st
	var sql_cmd string
	switch input.Strength {
	case 7:
		sql_cmd = Snprintf("attack_6days", input.Attribute, 0)
	case 30:
		sql_cmd = Snprintf("attack_29days", input.Attribute, 0)
	default:
		return output
	}
	rows := modelsPublic.Select_mysql(sql_cmd)
	if rows == nil {
		return output
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&(info.Country), &(info.Province),
			&(info.Attack_Type), &(info.IP),
			&(info.Source), &(info.Attack_Count)); err == nil {
			output.Info = append(output.Info, info)
			output.Count = output.Count + 1
		}
	}
	//fmt.Println(output)
	return output
}
func Statistics_Top10(input StatisticsAttackIn_st) StatisticsAttackOut_st {
	var output StatisticsAttackOut_st
	var info Info_st
	var sql_cmd string
	switch input.Strength {
	case 7:
		sql_cmd = Snprintf_top10("attack_6days", input.Attribute, 5)
	case 30:
		sql_cmd = Snprintf_top10("attack_29days", input.Attribute, 5)
	default:
		return output
	}
	rows := modelsPublic.Select_mysql(sql_cmd)
	if rows == nil {
		return output
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&(info.Country), &(info.Province),
			&(info.Attack_Type), &(info.IP),
			&(info.Source), &(info.Attack_Count)); err == nil {
			output.Info = append(output.Info, info)
			output.Count = output.Count + 1
		}
	}
	//fmt.Println(output)
	return output
}
func top5_select(table string) StatisticsAttackOut_st {
	var output StatisticsAttackOut_st

	var sql_cmd string
	var ip string
	var source string
	attack_count := 0
	sql_cmd1 := fmt.Sprintf(`select ip from ((select ip,sum(attack_count)as attack_count 
				from attack_days group by ip) union all (select ip,
				sum(attack_count) as attack_count from %s group 
				by ip)) as c group by ip order by sum(attack_count) desc 
				limit 5`, table)
	//fmt.Println(sql_cmd1)
	rows := modelsPublic.Select_mysql(sql_cmd1)
	if rows == nil {
		return output
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&ip); err == nil {
			// fmt.Println(ip)
			sql_cmd = fmt.Sprintf(`select source,sum(attack_count)as 
						attack_count from ((select ip,source,sum(attack_count)
						as attack_count from attack_days where ip = '%s' 
						group by source) union all (select ip,source,sum(attack_count)
						as attack_count from %s where ip = '%s' group 
						by source)) as c where ip = '%s' group by source order 
						by source asc`, ip, table, ip, ip)
			//fmt.Println(sql_cmd)
			var info Top5Info_st
			info.Ip = ip
			rows2 := modelsPublic.Select_mysql(sql_cmd)
			if rows2 == nil {
				return output
			}
			defer rows2.Close()
			for rows2.Next() {
				if err := rows2.Scan(&(source),
					&(attack_count)); err == nil {
					//fmt.Println(err.Error())
					//fmt.Println("info", info)
					switch source {
					case "IDS":
						info.Ids_attackcount = attack_count
					case "VDS":
						info.Vds_attackcount = attack_count
					case "WAF":
						info.Waf_attackcount = attack_count
					default:
					}
				}
			}
			output.Top5Info = append(output.Top5Info, info)
			output.Count = output.Count + 1
		}
	}
	return output
}
func Statistics_Top5(input StatisticsAttackIn_st) StatisticsAttackOut_st {
	var output StatisticsAttackOut_st
	if input.Attribute != "source" {
		return output
	}
	switch input.Strength {
	case 7:
		output = top5_select("attack_6days")
	case 30:
		output = top5_select("attack_29days")
	default:
		return output
	}
	//fmt.Println(output)
	return output
}
