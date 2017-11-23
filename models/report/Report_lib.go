// Report_lib
package report

import (
	//"database/sql"
	"apt-web-server_v2/models/modelsPublic"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func Security_Reportall() Security_st {
	var output Security_st
	sql_cmd := fmt.Sprintf(`select attack_type,sum(attack_count) from attack_days group by
							%s`, "attack_type")
	//fmt.Println(sql_cmd)
	rows := modelsPublic.Select_mysql(sql_cmd)
	if rows == nil {
		return output
	}
	defer rows.Close()
	for rows.Next() {
		var info Statistics_st
		if err := rows.Scan(&(info.Type), &(info.Count)); err == nil {
			output.Statistics = append(output.Statistics, info)
		}
	}
	sql_cmd = fmt.Sprintf(`select attack_type,ip from attack_days group by
							ip order by %s desc limit 5`, "sum(attack_count)")
	//fmt.Println(sql_cmd)
	rows = modelsPublic.Select_mysql(sql_cmd)
	if rows == nil {
		return output
	}
	defer rows.Close()
	for rows.Next() {
		var info Event_st
		if err := rows.Scan(&(info.Type), &(info.Ip)); err == nil {
			output.Event = append(output.Event, info)
		}
	}
	//fmt.Println(output)
	return output
}
func Security_Reportcondition(input Attribute_st) Security_st {
	var output Security_st
	var mysql_out Attribute_st
	sql_cmd := fmt.Sprintf(`select attack_type,sum(attack_count) from attack_days group 
		by attack_type order by %s desc`, "sum(attack_count)")
	fmt.Println(sql_cmd)
	rows := modelsPublic.Select_mysql(sql_cmd)
	if rows == nil {
		return output
	}
	defer rows.Close()
	for rows.Next() {
		var info Statistics_st
		if err := rows.Scan(&(info.Type), &(info.Count)); err == nil {
			for j, _ := range input.Info {
				if info.Type == input.Info[j] {
					mysql_out.Info = append(mysql_out.Info, info.Type)
					output.Statistics = append(output.Statistics, info)
				}
			}
		}
	}
	for i, _ := range input.Info {
		var q int
		for k, _ := range mysql_out.Info {
			if input.Info[i] == mysql_out.Info[k] {
				q = k - 1
				break
			}
			q = k
		}
		if q >= len(mysql_out.Info)-1 {
			var info1 Statistics_st
			info1.Type = input.Info[i]
			output.Statistics = append(output.Statistics, info1)
		}
	}

	sql_cmd = fmt.Sprintf(`select attack_type,ip from attack_days where 
		attack_type <> "other" group by ip order by %s desc limit 5`,
		"sum(attack_count)")
	//fmt.Println(sql_cmd)
	rows = modelsPublic.Select_mysql(sql_cmd)
	if rows == nil {
		return output
	}
	defer rows.Close()
	for rows.Next() {
		var info Event_st
		if err := rows.Scan(&(info.Type), &(info.Ip)); err == nil {
			output.Event = append(output.Event, info)
		}
	}
	output.Score = Security_Score("attack_days")
	//output.Score = Security_Score("attack_6days")
	//output.Score = Security_Score("attack_29days")
	//fmt.Println(output)
	return output
}
func Security_Score(input string) int {
	var output float32
	sql_cmd := fmt.Sprintf(`select sum(attack_count) from %s where attack_type='disclosure' 
	or attack_type='sqli' or attack_type='rce' or attack_type='privilege_gain' or 
	attack_type='information_leak' or attack_type='candc' or attack_type='trojan' 
	or attack_type='exploit' or attack_type='worm' or attack_type='webshell' or 
	attack_type='rfi'`, input)
	//fmt.Println(sql_cmd)
	rows := modelsPublic.Select_mysql(sql_cmd)
	if rows == nil {
		return int(output)
	}
	defer rows.Close()
	for rows.Next() {
		var count int
		if err := rows.Scan(&count); err == nil {
			fmt.Println(count)
			if count >= 0 && count <= 100 {
				output = 5*float32(count)/100 + output
			} else if count > 100 && count <= 500 {
				output = 10*float32(count-100)/400 + 5 + output
			} else {
				output = 15 + 5 + 10 + output
			}
		}
	}
	sql_cmd = fmt.Sprintf(`select sum(attack_count) from %s where attack_type='malware' 
	or attack_type='backdoor' or attack_type='spyware' or 
	attack_type='virus' or attack_type='hacktool' or attack_type='application_attack'`, input)
	//fmt.Println(sql_cmd)
	rows = modelsPublic.Select_mysql(sql_cmd)
	if rows == nil {
		return int(output)
	}
	defer rows.Close()
	for rows.Next() {
		var count int
		if err := rows.Scan(&count); err == nil {
			fmt.Println(count)
			if count >= 0 && count <= 100 {
				output = 2*float32(count)/100 + output
			} else if count > 100 && count <= 500 {
				output = 4*float32(count-100)/400 + 2 + output
			} else {
				output = 6 + 2 + 4 + output
			}
		}
	}
	sql_cmd = fmt.Sprintf(`select sum(attack_count) from %s where attack_type='dos' 
	or attack_type='reputation_ip' or attack_type='lfi' or attack_type='xss' or 
	attack_type='injection_php' or attack_type='generic' or attack_type='protocol' 
	or attack_type='fixation' or attack_type='scaningprobe' or attack_type='reputation_scanner' 
	or attack_type='reputation_scripting' or attack_type='reputation_crawler' or
	attack_type='ddos' or attack_type='misc_attack'`, input)
	//fmt.Println(sql_cmd)
	rows = modelsPublic.Select_mysql(sql_cmd)
	if rows == nil {
		return int(output)
	}
	defer rows.Close()
	for rows.Next() {
		var count int
		if err := rows.Scan(&count); err == nil {
			fmt.Println(count)
			if count >= 0 && count <= 100 {
				output = 10*float32(count)/100 + output
			} else if count > 100 && count <= 500 {
				output = 15*float32(count-100)/400 + 1 + output
			} else {
				output = 4 + 1 + 2 + output
			}
		}
	}
	//fmt.Println(output)
	if float32(int(output)) < output {
		output = output + 1
	}
	if output > 30 {
		return 70
	}
	return 100 - int(output)
}
