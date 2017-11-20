// Equipment_st.go
package equipment

type Detailinfo_st struct {
	Service_ip       string `json:"service_ip"`       //服务器IP
	Service_name     string `json:"service_name"`     //服务名称
	Service_type     string `json:"service_type"`     //TCP/UDP
	Service_version  string `json:"service_version"`  //服务器版本
	Service_platform string `json:"service_platform"` //服务器平台(linux/windows)
	Service_port     int    `json:"service_port"`     //服务器端口
	Service_banner   string `json:"service_banner"`   //服务内容信息
}

type Showdetailinfo_st struct {
	Counts     int             `json:"counts"`
	Detailinfo []Detailinfo_st `json:"elements"`
}
type Departmentip_st struct {
	DepartmentId   int      `json:"departmentId"`   //部门id
	Departmentname string   `json:"departmentname"` //部门名称
	Ip             []string `json:"ip"`             //部门资产ip
}

type Departmentshowinfo_st struct {
	DepartmentId   int    `json:"departmentId"`   //部门id
	Departmentname string `json:"departmentname"` //部门名称
	Equipmentcount int    `json:"equipmentcount"` //资产数
}

type Departmentshow_st struct {
	Count int                     `json:"counts"`
	Info  []Departmentshowinfo_st `json:"elements"`
}

type Equipmentshowinfo_st struct {
	Ip             string `json:"ip"`
	Os_type        string `json:"os_type"`   //系统类型
	Alias          string `json:"alias"`     //标签
	Authority      int    `json:"authority"` //权限
	Attack_count   int    `json:"attack_count"`
	Data_source    int    `json:"data_source"`    //来源
	Departmentname string `json:"departmentname"` //部门
	Service_port   string `json:"service_port"`   //服务端口
}

type Equipmentshow_st struct {
	Count int                    `json:"counts"`
	Info  []Equipmentshowinfo_st `json:"elements"`
}

type Limt_st struct {
	Type         string
	Page         int
	Count        int
	Lies         string
	Orderby      string
	Ip           string
	Os_type      string
	Data_source  int
	DepartmentId int
}

type EquipmentIP_st struct {
	Ip [10]string
}

type Equipmentinfo_st struct {
	Ip           string `json:"ip"`
	Os_type      string `json:"os_type"`
	Alias        string `json:"alias"`
	DepartmentId int    `json:"departmentId"` //部门id
}

type Importdate_st struct {
	Counts    int                `json:"counts"`
	Equipment []Equipmentinfo_st `json:"elements"`
}
