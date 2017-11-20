// Sniffxml_creat.go
package sniff

import (
	"apt-web-server_v2/models/db"
	"encoding/xml"
	"fmt"
	//"io/ioutil"
	"database/sql"
	"os"
	//"os/exec"

	_ "github.com/go-sql-driver/mysql"
)

type sub_condition struct {
	S3_Count int    `xml:"count,attr"`
	Content  string `xml:"content"`
}
type condition struct {
	S2_Count      int           `xml:"count,attr"`
	Type          string        `xml:"type,attr"`
	Sub_condition sub_condition `xml:"sub_condition"`
}
type exec_time struct {
	Start_date  int    `xml:"start_date,attr"`
	End_date    int    `xml:"end_date,attr"`
	Period_type string `xml:"period_type,attr"`
}
type control_condition struct {
	S1_Count  int         `xml:"count,attr"`
	Condition []condition `xml:"condition"`
}
type plcy struct {
	Plcy_id           int               `xml:"plcy_id,attr"`
	Plcy_type         string            `xml:"plcy_type,attr"`
	Priority          int               `xml:"priority,attr"`
	Control_condition control_condition `xml:"control_condition"`
	Exec_time         exec_time         `xml:"exec_time"`
}
type sniffplcy struct {
	XMLName      xml.Name `xml:"plcy_set"`
	Plcy_count   int      `xml:"plcy_count,attr"`
	Plcy_version int      `xml:"plcy_version"`
	Plcy         []plcy   `xml:"plcy"`
}
type condition_type struct {
	Type string
	Vlue string
}
type row struct {
	Plcy_name        string
	Affect_starttime int
	Affect_endtime   int
	Pcap_path        string
	Plcy_status      int
	Plcy_date        string
	Condtion_type    [5]condition_type
}

func deal_subcondition(a sub_condition, Count int, Content string) sub_condition {
	a.S3_Count = Count
	a.Content = Content
	return a
}
func deal_condition(a condition, Xml_sub_condition sub_condition, Count int, Type string, Content string) condition {
	Xml_sub_condition = deal_subcondition(Xml_sub_condition, 1, Content)
	a.S2_Count = Count
	a.Type = Type
	a.Sub_condition = Xml_sub_condition
	return a
}
func deal_control_condition(a control_condition, Xml_condition condition, Count int, Condtion_type condition_type) control_condition {
	var Xml_sub_condition sub_condition
	Xml_condition = deal_condition(Xml_condition, Xml_sub_condition, 1, Condtion_type.Type, Condtion_type.Vlue)

	a.S1_Count = Count
	a.Condition = append(a.Condition, Xml_condition)
	return a
}
func deal_plcy(a plcy, Xml_control_condition control_condition, Plcy_id int, Plcy_type string, Priority int, Rows row) plcy {
	var Xml_condition condition
	count := 0
	for i, _ := range Rows.Condtion_type {
		if Rows.Condtion_type[i].Vlue != "" {
			count = count + 1
		}
	}
	for i, _ := range Rows.Condtion_type {
		if Rows.Condtion_type[i].Vlue != "" {
			Xml_control_condition = deal_control_condition(Xml_control_condition, Xml_condition, count, Rows.Condtion_type[i])
		}
	}

	a.Exec_time.Start_date = Rows.Affect_starttime
	a.Exec_time.End_date = Rows.Affect_endtime
	a.Exec_time.Period_type = "none"
	a.Plcy_id = Plcy_id
	a.Plcy_type = Plcy_type
	a.Priority = Priority
	a.Control_condition = Xml_control_condition
	return a
}
func creat_pcaprule(db *sql.DB) bool {
	var pcap_plcy []string
	var tmp string
	var tmp1 string
	var tmp2 string
	rows, err := db.Query("select plcy_name,pcap_path,pcap_plcy from sniff_plcy_ui where plcy_status = 1")
	defer rows.Close()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	for rows.Next() {
		if err = rows.Scan(&tmp, &tmp1, &tmp2); err == nil {
			tmp = fmt.Sprintf("%s|%s|%s", tmp, tmp1, tmp2)
			pcap_plcy = append(pcap_plcy, tmp)
		}
	}
	//fileName := "test.dat"
	dstFile, err := os.Create("./rules/rule")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer dstFile.Close()
	for i, _ := range pcap_plcy {
		dstFile.WriteString(pcap_plcy[i] + "\n")
	}
	return true
}
func deal_sniffplcy(a sniffplcy, Xml_plcy plcy, Plcy_version int) sniffplcy {
	//mysql 查找plcy_id 有一个row[]
	var Rows row
	var plcy_id int
	var Plcy_count int
	/*
		db, err := sql.Open("mysql", "root:byzoro@tcp(10.66.2.114:3306)/test?charset=utf8")
		if err != nil {
			fmt.Printf("connect err")
		}
		defer db.Close()
	*/
	rows1, err2 := db.DB.Query("select COUNT(plcy_name) from sniff_plcy_dpi where plcy_status = 1")
	defer rows1.Close()
	if err2 != nil {
		fmt.Println(err2.Error())
		return a
	}
	for rows1.Next() {
		if err := rows1.Scan(&Plcy_count); err != nil {
			return a
		}
	}

	Mysql := fmt.Sprintf(`select plcy_id,plcy_name,src_ip,dst_ip,src_port,dst_port,
	proto,affect_time_start,affect_time_end,plcy_status from sniff_plcy_dpi where 
	plcy_status = 1`)
	rows, err1 := db.DB.Query(Mysql)
	defer rows.Close()
	if err1 != nil {
		fmt.Println(err1.Error())
		return a
	}

	rst := creat_pcaprule(db.DB)
	if rst == false {
		return a
	}

	a.Plcy_count = Plcy_count
	a.Plcy_version = Plcy_version
	for rows.Next() {
		if err := rows.Scan(&plcy_id, &(Rows.Plcy_name),
			&(Rows.Condtion_type[0].Vlue), &(Rows.Condtion_type[1].Vlue),
			&(Rows.Condtion_type[2].Vlue), &(Rows.Condtion_type[3].Vlue),
			&(Rows.Condtion_type[4].Vlue), &(Rows.Affect_starttime),
			&(Rows.Affect_endtime), &(Rows.Plcy_status)); err == nil {
			//fmt.Println("=================")
			/*
				fmt.Println(Rows.Plcy_name)
				fmt.Println(Rows.Condtion_type[0].Vlue)
				fmt.Println(Rows.Condtion_type[1].Vlue)
				fmt.Println(Rows.Condtion_type[2].Vlue)
				fmt.Println(Rows.Condtion_type[3].Vlue)
				fmt.Println(Rows.Condtion_type[4].Vlue)
				fmt.Println(Rows.Pcap_path)
				fmt.Println(Rows.Affect_starttime)
				fmt.Println(Rows.Affect_endtime)
				fmt.Println(Rows.Plcy_status)
				fmt.Println(Rows.Plcy_date)
			*/
			Rows.Condtion_type[0].Type = "src_ipv4"
			Rows.Condtion_type[1].Type = "dst_ipv4"
			Rows.Condtion_type[2].Type = "src_port"
			Rows.Condtion_type[3].Type = "dst_port"
			Rows.Condtion_type[4].Type = "l4_proto"
			//回调生成抓包条件
			var Xml_control_condition control_condition
			Plcy_type := "sniff"
			Priority := 600
			Xml_plcy = deal_plcy(Xml_plcy, Xml_control_condition, plcy_id, Plcy_type, Priority, Rows)

			a.Plcy = append(a.Plcy, Xml_plcy)
		}
	}
	return a
}

func Sniff_xml_creat(creat_path string) bool {
	var Xml_sniffplcy sniffplcy
	var Xml_plcy plcy
	plcy_version := 110
	Xml_sniffplcy = deal_sniffplcy(Xml_sniffplcy, Xml_plcy, plcy_version)
	//fmt.Println(Xml_sniffplcy)
	//fmt.Println(v)
	if Xml_sniffplcy.Plcy_version == 0 {
		return false
	}
	output, err := xml.MarshalIndent(Xml_sniffplcy, " ", " ")
	if err != nil {
		fmt.Printf("error %Xml_sniffplcy", err)
		return false
	}
	file, _ := os.Create(creat_path)
	defer file.Close()
	file.Write([]byte(xml.Header))
	file.Write(output)
	//fmt.Println("Hello World!")
	/*tar -zcvf */
	//cmd := exec.Command("tar", "-zcvf", "PM_Upgrade_resc.tar.gz", "./xml/*.xml", "2>&1", ">/dev/null ")
	//_, err3 := cmd.Output()
	//if err3 != nil {
	//	fmt.Println(err3.Error())
	//}
	/*tar end*/
	return true
}
