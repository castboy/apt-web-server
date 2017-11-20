/********http流量检测********/
package waf

import (
	"apt-web-server_v2/models/modelsPublic"
)

/********共有结构********
type TblPublicPara struct {
	Start int64
	End   int64
}
*/
type OfflineTage struct {
	Tage       string
	TaskName   string
	CreateTime int64
	MergeTag   string
	DupTag     string
}

/********获取http流量攻击数(天)********/
//HFA:"HttpFlowAttackTrend"
type TblHFATSearchPara struct {
	Type   string
	PField modelsPublic.TblPublicPara
}
type TblHFAT struct {
	Id   int64
	Time int64
	TblHFAL
}
type TblHFATData struct {
	//	Type     string       `json:"type"`
	Counts   int64         `json:"counts"`
	Elements []TblHFATList `json:"elements"`
}

/*
type TblHFAContent struct {
	Time  int64 `json:"time"`
	Times int64 `json:"times"`
}*/
type TblHFATList struct {
	TblHFACRequest
	//TblHFAL
}

/*************************/

/********获取http流量攻击分类数(天)********/
//HFAC:"HttpFlowAttackClassify"
type TblHFACSearchPara struct {
	PField modelsPublic.TblPublicPara
}

type TblHFAC struct {
	Id int64
	TblHFACContent
}

type TblHFACData struct {
	Counts   int64         `json:"counts"`
	Classify []TblHFACList `json:"classify"`
}

type TblHFACContent struct {
	Time int64 `json:"time"`
	TblHFAL
}
type TblHFACRequest struct {
	HFACUniversal
}
type HFACUniversal struct {
	AttDisclosure   int64 `json:"disclosure"`
	AttDdos         int64 `json:"ddos"`
	AttReputationIp int64 `json:"reputation_ip"`
	AttLfi          int64 `json:"lfi"`
	AttSqli         int64 `json:"sqli"`
	AttXSS          int64 `json:"xss"`
	AttInjectionPHP int64 `json:"injection_php"`
	AttGeneric      int64 `json:"generic"`
	AttRce          int64 `json:"rce"`
	AttProtocol     int64 `json:"protocol"`
	AttRfi          int64 `json:"rfi"`
	AttFixation     int64 `json:"fixation"`
	Scaning         int64 `json:"scaning"`
	Other           int64 `json:"other"`
}
type HFACScaningProbe struct {
	AttReputScanner  int64 `json:"reputation_scanner"`
	AttReputSripting int64 `json:"reputation_scripting"`
	AttReputCrawler  int64 `json:"reputation_crawler"`
}
type TblHFAL struct {
	HFACUniversal
	//HFACScaningProbe
}
type HFACContent struct {
	Time int64 `json:"time"`
	TblHFACRequest
}
type TblHFACList struct {
	HFACContent
}

/*************************/

/********获取http流量攻击详情(天)********/
//HFAD:"HttpFlowAttackDetails"
type TblHFADSearchPara struct {
	Id        int32
	Type      string
	LastCount int32
	KeyWord   string
	Page      int32
	Count     int32
	Sort      string
	Order     string
	PField    modelsPublic.TblPublicPara
	OfflineTage
}
type DupIdList struct {
	TblId  int
	TaskId int
}
type TblHFAD struct {
	Id int64
	TblHFADContent
	DupIdList
}
type TblHFADData struct {
	Totality int64         `json:"total"`
	Counts   int64         `json:"counts"`
	Elements []TblHFADList `json:"elements"`
}

type WafAlertRule struct {
	Version string `json:"Ver"`
	Data    string `json:"Data"`
	File    string `json:"File"`
	Line    int    `json:"Line"`
	Id      int    `json:"Id"`
}
type TblHFADContent struct {
	modelsPublic.AreaMsg
	Time     int64        `json:"Time"`
	Attack   string       `json:"Attack"`
	Client   string       `json:"Client"`
	Rev      string       `json:"Rev"`
	Severity int          `json:"Severity"`
	Maturity int          `json:"Maturity"`
	Accuracy int          `json:"Accuracy"`
	HostName string       `json:"Hostname"`
	UniqueId string       `json:"Unique_id"`
	Ref      string       `json:"Ref"`
	Tags     string       `json:"Tags"`
	Rule     WafAlertRule `json:"Rule"`
	Version  string       `json:"Version"`
	Request  string       `json:"Request"`
	Response string       `json:"Response"`
	Message  string       `json:"Msg"`
	Uri      string       `json:"Uri"`
}
type TblHFADList struct {
	TblHFADContent
}
