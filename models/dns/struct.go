/********DNS统计(天)********/
package dns

/********共有结构********/
type TblPublicPara struct {
	Start int64
	End   int64
}

//DNSS:"DomainNameServerStatistics"
type TblDNSSSearchPara struct {
	Type      string
	LastCount int32
	KeyWord   string
	Page      int32
	Count     int32
	Sort      string
	Order     string
	Unit      string
	Ip        string
	Domain    string
	PField    TblPublicPara
}
type TblDNSSPublic struct {
	Time  string `json:"time"`
	Count int64  `json:"count"`
}

type TblDNSSData struct {
	Totality int64         `json:"total"`
	Counts   int64         `json:"counts"`
	Elements []TblDNSSList `json:"elements"`
}
type TblDNSSIpData struct {
	Totality int64           `json:"total"`
	Counts   int64           `json:"counts"`
	Elements []TblDNSSIpList `json:"elements"`
}
type TblDNSSDomainData struct {
	Totality int64               `json:"total"`
	Counts   int64               `json:"counts"`
	Elements []TblDNSSDomainList `json:"elements"`
}

type TblDNSS struct {
	//Id int64
	TblDNSSContent
}
type TblDNSSIp struct {
	//Id int64
	TblDNSSIpContent
}
type TblDNSSDomain struct {
	//Id int64
	TblDNSSDomainContent
}

type TblDNSSContent struct {
	Ip     string `json:"ip"`
	Domain string `json:"domain"`
	TblDNSSPublic
}
type TblDNSSIpContent struct {
	Domain string `json:"domain"`
	TblDNSSPublic
}
type TblDNSSDomainContent struct {
	Domain string `json:"domain"`
	TblDNSSPublic
}

type TblDNSSList struct {
	TblDNSSContent
}
type TblDNSSIpList struct {
	TblDNSSIpContent
}
type TblDNSSDomainList struct {
	TblDNSSDomainContent
}
