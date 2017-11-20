// Statistics_lib.go
package index

import (
	"fmt"
)

func Snprintf(table string, Attribute string, top int) string {
	var sql_cmd string
	if top > 0 {
		if (Attribute == "country") || (Attribute == "source") {
			sql_cmd = fmt.Sprintf(`select country,province,attack_type,ip,source,
						sum(attack_count) as attack_count from ((select 
						country,province,attack_type,ip,source,sum(attack_count) 
						as attack_count from attack_days group by %s) 
						union all (select country,province,attack_type,ip,source,
						sum(attack_count) as attack_count from %s group 
						by %s)) as c group by %s order by attack_count desc 
						limit %d`, Attribute, table, Attribute, Attribute, top)
		} else {
			sql_cmd = fmt.Sprintf(`select country,province,attack_type,ip,source,
						sum(attack_count) as attack_count from ((select 
						country,province,attack_type,ip,source,sum(attack_count) 
						as attack_count from attack_days group by %s) 
						union all (select country,province,attack_type,ip,source,
						sum(attack_count) as attack_count from %s group 
						by %s)) as c where country = '中国' and province <>'中国' 
						group by %s order by attack_count desc limit %d`, Attribute,
				table, Attribute, Attribute, top)
		}

	} else if top == 0 {
		if (Attribute == "country") || (Attribute == "source") {
			sql_cmd = fmt.Sprintf(`select country,province,attack_type,ip,source,
						sum(attack_count) as attack_count from ((select 
						country,province,attack_type,ip,source,sum(attack_count) 
						as attack_count from attack_days group by %s) 
						union all (select country,province,attack_type,ip,source,
						sum(attack_count) as attack_count from %s group 
						by %s)) as c group by %s`, Attribute, table, Attribute, Attribute)
		} else {
			sql_cmd = fmt.Sprintf(`select country,province,attack_type,ip,source,
						sum(attack_count) as attack_count from ((select 
						country,province,attack_type,ip,source,sum(attack_count) 
						as attack_count from attack_days group by %s) 
						union all (select country,province,attack_type,ip,source,
						sum(attack_count) as attack_count from %s group 
						by %s)) as c where country='中国' and province <>'中国' 
						group by %s`, Attribute, table, Attribute, Attribute)
		}

	}
	//fmt.Println(sql_cmd)
	return sql_cmd
}
func Snprintf_top5(table string, Attribute string, top int) string {
	var sql_cmd string
	if top > 0 {
		sql_cmd = fmt.Sprintf(`select country,province,attack_type,ip,source,
						sum(attack_count) as attack_count from ((select 
						country,province,attack_type,ip,source,sum(attack_count) 
						as attack_count from attack_days group by %s) 
						union all (select country,province,attack_type,ip,source,
						sum(attack_count) as attack_count from %s group 
						by %s)) as c group by %s order by attack_count desc 
						limit %d`, Attribute, table, Attribute, Attribute, top)
	} else if top == 0 {
		sql_cmd = fmt.Sprintf(`select country,province,attack_type,ip,source,
						sum(attack_count) as attack_count from ((select 
						country,province,attack_type,ip,source,sum(attack_count) 
						as attack_count from attack_days group by %s) 
						union all (select country,province,attack_type,ip,source,
						sum(attack_count) as attack_count from %s group 
						by %s)) as c group by %s`, Attribute, table, Attribute, Attribute)
	}
	fmt.Println(sql_cmd)
	return sql_cmd
}
func Snprintf_top10(table string, Attribute string, top int) string {
	var sql_cmd string
	if top > 0 {
		sql_cmd = fmt.Sprintf(`select country,province,attack_type,ip,source,
					sum(attack_count) as attack_count from ((select 
					country,province,attack_type,ip,source,sum(attack_count) 
					as attack_count from attack_days group by %s) 
					union all (select country,province,attack_type,ip,source,
					sum(attack_count) as attack_count from %s group 
					by %s)) as c group by %s order by attack_count desc 
					limit %d`, Attribute, table, Attribute, Attribute, top)
	} else if top == 0 {
		sql_cmd = fmt.Sprintf(`select country,province,attack_type,ip,source,
					sum(attack_count) as attack_count from ((select 
					country,province,attack_type,ip,source,sum(attack_count) 
					as attack_count from attack_days group by %s) 
					union all (select country,province,attack_type,ip,source,
					sum(attack_count) as attack_count from %s group 
					by %s)) as c group by %s`, Attribute, table, Attribute, Attribute)
	}
	//fmt.Println(sql_cmd)
	return sql_cmd
}
