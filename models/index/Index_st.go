// Statistics_st.go
package index

type StatisticsAttackIn_st struct {
	Method    string
	Attribute string
	Strength  int
}
type StatisticsAttackOut_st struct {
	Count       int           `json:"counts"`
	Info        []Info_st     `json:"elements"`
	Top5Info    []Top5Info_st `json:"equip_top5"`
	CountryInfo []Country_st  `json:"country_percent"`
}
type Info_st struct {
	Country      string `json:"country"`
	Province     string `json:"province"`
	Attack_Type  string `json:"attack_type"`
	IP           string `json:"ip"`
	Source       string `json:"source"`
	Attack_Count int    `json:"attack_count"`
}
type Top5Info_st struct {
	Ip              string `json:"ip"`
	Ids_attackcount int    `json:"ids"`
	Vds_attackcount int    `json:"vds"`
	Waf_attackcount int    `json:"waf"`
}
type Country_st struct {
	Home   int `json:"home"`
	Abroad int `json:"abroad"`
}
type MapCount_st struct {
	Area  string `json:"name"`
	Count int    `json:"value"`
}
type CityName_st struct {
	Name string `json:"name"`
}
type AttackInfo_st struct {
	MysqlId string `json:"id"`
	SrcIp   string `json:"src_ip"`
	//SrcCountry  string `json:"src_country"`
	SrcProvince  string `json:"src_province"`
	DestProvince string `json:"dest_province"`
	DestIp       string `json:"dest_ip"`
	Proto        string `json:"proto"`
	DestPort     int    `json:"dest_port"`
	AttackType   string `json:"attack_type"`
	//Servity      int    `json:"servity"`
	Time     string `json:"time"`
	FromType string `json:"fromtype"`
}
type AttackDay_st struct {
	Count    int `json:"count"`
	WafCount int `json:"waf_count"`
	IdsCount int `json:"ids_count"`
	VdsCount int `json:"vds_count"`
}
type MonitorAttack_st struct {
	Count      int             `json:"counts"`
	MapCount   []MapCount_st   `json:"mapcount"`
	Attackcity [][]CityName_st `json:"attackcity"`
	Ids_Attack []AttackInfo_st `json:"ids_attack"`
	Vds_Attack []AttackInfo_st `json:"vds_attack"`
	Waf_Attack []AttackInfo_st `json:"waf_attack"`
	Attack_Day AttackDay_st    `json:attack_day`
}
