/**********紧急事件**********/
package urgency

import (
	"apt-web-server_v2/models/modelsPublic"
)

/********共有结构********
type TblPublicPara struct {
	Start int64
	End   int64
}
*/
/********获取紧急事件数（天）********/
type TblUrgencySearchPara struct {
	Type      string
	LastCount int32
	KeyWord   string
	Page      int32
	Count     int32
	Sort      string
	Order     string
	Unit      string
	PField    modelsPublic.TblPublicPara
}

type TblUrgency struct {
	Id int64
	TblUrgencyContent
}

type TblUrgencyData struct {
	//	Totality int64
	Counts   int64            `json:"counts"`
	Elements []TblUrgencyList `json:"elements"`
}

type TblUrgencyContent struct {
	Date  string `json:"date"`
	Times int64  `json:"times"`
}

type TblUrgencyList struct {
	TblUrgencyContent
}

/*************************/

/********获取新紧急事件分类数(天)********/
type TblUgcClassifySearchPara struct {
	Type      string
	LastCount int32
	KeyWord   string
	Page      int32
	Count     int32
	Sort      string
	Order     string
	Unit      string
	PField    modelsPublic.TblPublicPara
}

type TblUgcClassify struct {
	Id int64
	TblUgcClassifyContent
}
type TblUgcClassifyData struct {
	Totality int64                `json:"total"`
	Counts   int64                `json:"counts"`
	Elements []TblUgcClassifyList `json:"elements"`
}
type TblUgcClassifyContent struct {
	AttackType string `json:"attack_type"`
	Severity   int    `json:"severity"`
	Time       int64  `json:"time"`
	Srcip      string `json:"src"`
	Destip     string `json:"dest"`
	Details    string `json:"details"`
}

type TblUgcClassifyList struct {
	TblUgcClassifyContent
}

type TblUgcCCount struct {
	Id int64
	TblUgcCCountContent
}
type TblUgcCCountData struct {
	//	Totality int64
	Counts   int64              `json:"counts"`
	Elements []TblUgcCCountList `json:"elements"`
}

type TblUgcCCountContent struct {
	Time               int64 `json:"time"`
	Webshell           int64 `json:"webshell"`
	ExceptionalVisit   int64 `json:"exceptionalvisit"`
	AbnormalConnection int64 `json:"abnormal_connection"`
	Sqli               int64 `json:"sqli"`
	Xss                int64 `json:"xss"`
	InjectionPHP       int64 `json:"injection_php"`
	Rfi                int64 `json:"rfi"`
}

type TblUgcCCountList struct {
	TblUgcCCountContent
}

/*
type TblUgcClassifyContent struct {
	SqlInject         int64 `json:"sqlInject"`
	ChickenBehave     int64 `json:"chickenBehave"`
	WormVirus         int64 `json:"wormVirus"`
	TrojanBack        int64 `json:"trojanBack"`
	Ddos              int64 `json:"ddos"`
	UnlicenseDownload int64 `json:"unlicenseDownload"`
	Webshell          int64 `json:"webshell"`
}
*/
/*************威胁模型统计************/
type TblUgcMCountSearchPara struct {
	Type      string
	LastCount int32
	KeyWord   string
	Page      int32
	Count     int32
	Sort      string
	Order     string
	Unit      string
	PField    modelsPublic.TblPublicPara
}
type TblUgcMCount struct {
	Id int64
	TblUgcMCountContent
}
type TblUgcMCountData struct {
	//	Totality int64
	Counts   int64              `json:"counts"`
	Elements []TblUgcMCountList `json:"elements"`
}

type TblUgcMCountContent struct {
	Time               int64 `json:"time"`
	Webshell           int64 `json:"webshell"`
	ExceptionalVisit   int64 `json:"exceptionalvisit"`
	AbnormalConnection int64 `json:"abnormal_connection"`
	BruteForce         int64 `json:"bruteforce"`
	PortScan           int64 `json:"portscan"`
}

type TblUgcMCountList struct {
	TblUgcMCountContent
}

/*************************/

/********获取紧急事件详情(天)********/
//UgcD:"UrgencyDetails"
type TblUgcDSearchPara struct {
	Type      string
	LastCount int32
	KeyWord   string
	Page      int32
	Count     int32
	Sort      string
	Order     string
	PField    modelsPublic.TblPublicPara
}
type TblUgcD struct {
	Id int64
	TblUgcDContent
}
type TblUgcDData struct {
	Totality int64         `json:"total"`
	Counts   int64         `json:"counts"`
	Elements []TblUgcDList `json:"elements"`
}
type TblUgcDContent struct {
	Time       int64  `json:"time"`
	SrcIp      string `json:"srcip"`
	SrcPort    int    `json:"srcport"`
	DestIp     string `json:"destip"`
	DestPort   int    `json:"destport"`
	Proto      string `json:"protocol"`
	ServerName string `json:"servername"`
	AttackType string `json:"attacktype"`
	Serverity  string `json:"serverity"`
	AttackerOS string `json:"attackeros"`
	AttackedOS string `json:"attackedos"`
	Details    string `json:"details"`
}
type TblUgcDList struct {
	TblUgcDContent
}

/********获取攻击数（天）********/
type TblAttackCountSearchPara struct {
	Type      string
	LastCount int32
	KeyWord   string
	Page      int32
	Count     int32
	Sort      string
	Order     string
	Unit      string
	PField    modelsPublic.TblPublicPara
}

type TblAttackCount struct {
	Id int64
	TblAttackCountContent
}

type TblAttackCountData struct {
	Totality int64                `json:"allipnum"`
	Counts   int64                `json:"allattackednum"`
	Elements []TblAttackCountList `json:"elements"`
}

type TblAttackCountContent struct {
	Time        string `json:"time"`
	IpCount     int64  `json:"ipnum"`
	AttackCount int64  `json:"attackcount"`
}

type TblAttackCountList struct {
	TblAttackCountContent
}

/**********最新紧急事件**********/
//UgcL:"UrgencyLatest"
type TblUgcLSearchPara struct {
	Type      string
	LastCount int32
	KeyWord   string
	Page      int32
	Count     int32
	Sort      string
	Order     string
	PField    modelsPublic.TblPublicPara
}
type TblUgcL struct {
	Id int64
	TblUgcLContent
}
type TblUgcLData struct {
	//Totality int64         `json:"total"`
	Counts   int64         `json:"counts"`
	Elements []TblUgcLList `json:"elements"`
}
type TblUgcLContent struct {
	Time       int64  `json:"time"`
	DestIp     string `json:"destip"`
	AttackType string `json:"attacktype"`
}
type TblUgcLList struct {
	TblUgcLContent
}
