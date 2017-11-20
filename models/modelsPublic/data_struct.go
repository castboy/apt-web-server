package modelsPublic

//	"fmt"

/********共有结构********/
type TblPublicPara struct {
	Start int64
	End   int64
}

type OfflineTage struct {
	Tage       string
	TaskName   string
	CreateTime int64
	MergeTag   string
	DupTag     string
}

type IpInfo struct {
	Country  string `json:"Country"`
	Province string `json:"Province"`
	City     string `json:"City"`
	Lat      string `json:"Lat"`
	Lng      string `json:"Lng"`
	//Organization    string `json:"Organization"`
	//Network         string `json:"Network"`
	//TimeZone        string `json:"TimeZone"`
	//UTC             string `json:"UTC"`
	//RegionalismCode string `json:"RegionalismCode"`
	//PhoneCode       string `json:"PhoneCode"`
	//CountryCode     string `json:"CountryCode"`
	//ContinentCode   string `json:"ContinentCode"`
}
type AreaMsg struct {
	SrcIp      string `json:"Src_ip"`
	SrcPort    string `json:"Src_port"`
	SrcIpInfo  IpInfo `json:"Src_ip_info"`
	DestIp     string `json:"Dest_ip"`
	DestPort   string `json:"Dest_port"`
	DestIpInfo IpInfo `json:"Dest_ip_info"`
	Proto      string `json:"Proto"`
	Operators  string `json:"Operators"`
}

/***********************/
/********获取攻击天数********
type TblAttackSearchPara struct {
	PField TblPublicPara
}

type TblAttack struct {
	Id int64
	TblAttackContent
}

type TblAttackData struct {
	Totality int64           `json:"total"`
	Counts   int64           `json:"counts"`
	Elements []TblAttackList `json:"elements"`
}

type TblAttackContent struct {
	Time  int64 `json:"time"`
	Times int64 `json:"times"`
}

type TblAttackList struct {
	TblAttackContent
}

/**************************/

/********获取网络流量数********
//curl参数
type TblNetFlowSearchPara struct {
	Direction string
	AssetIP   string
	Protocol  string
	Unit      string
	PField    TblPublicPara
}
type TblNetFlow struct {
	Id        int64
	Direction string
	AssetIP   string
	Protocol  string
	TblNetFlowContent
}

//分表数据
type TblNetFlowCount struct {
	Id int64
	TblNetFlowContent
}
type TblNetFlowIP struct {
	Id      int64
	AssetIP string
	TblNetFlowContent
}
type TblNetFlowD struct {
	Id        int64
	Direction string
	TblNetFlowContent
}
type TblNetFlowP struct {
	Id       int64
	Protocol string
	TblNetFlowContent
}
type TblNetFlowIPD struct {
	Id        int64
	AssetIP   string
	Direction string
	TblNetFlowContent
}
type TblNetFlowIPP struct {
	Id       int64
	AssetIP  string
	Protocol string
	TblNetFlowContent
}
type TblNetFlowDP struct {
	Id        int64
	Direction string
	Protocol  string
	TblNetFlowContent
}

//data字段结构
type TblNetFlowData struct {
	Totality int64            `json:"total"`
	Counts   int64            `json:"counts"`
	Elements []TblNetFlowList `json:"elements"`
}

//返回数据
type TblNetFlowContent struct {
	Time int64 `json:"time"`
	Flow int64 `json:"flow"`
}
type TblNetFlowList struct {
	TblNetFlowContent
}

/********************************/

/********获取所有资产********
//type TblAssetSearchPara struct {
//	PField TblPublicPara
//}

type TblAsset struct {
	Id int64
	TblAssetContent
}

type TblAssetData struct {
	//	Totality int64
	Counts   int64          `json:"counts"`
	Elements []TblAssetList `json:"elements"`
}

type TblAssetContent struct {
	ServerIp     string
	OsVersion    string
	OpenPortsNum int
	AttacksNum   int64
}

type TblAssetList struct {
	TblAssetContent
}

/*************************/

/********获取资产开放端口********
//AOP:"AssetOpenPort"
type TblAOPSearchPara struct {
	Asset string
}

type TblAOP struct {
	Id int64
	TblAssetIpContent
}

type TblAOPData struct {
	Counts   int64            `json:"counts"`
	Elements []TblAssetIpList `json:"elements"`
}

type TblAssetIpContent struct {
	ServerIp    string             `json:"serverIp"`
	OpenedPorts []TblAssetPortList `json:"openedPorts"`
}

type TblAssetIpList struct {
	TblAssetIpContent
}
type TblAssetPortContent struct {
	Port int32 `json:"port"`
}

type TblAssetPortList struct {
	TblAssetPortContent
}

/*************************/

/********获取攻击者数(粒度天)********
type TblAttacksNumSearchPara struct {
	PField TblPublicPara
	Asset  string
}

type TblAttacksNum struct {
	Id    int64
	Asset string
	TblAttacksNumContent
}

type TblAttacksNumData struct {
	//	Totality int64
	Counts   int64               `json:"counts"`
	Elements []TblAttacksNumList `json:"elements"`
}

type TblAttacksNumContent struct {
	Time int64
	Num  int64
}

type TblAttacksNumList struct {
	TblAttacksNumContent
}

/*************************/

/********最高威胁资产top********
//ATT:"AssetThreatTop"
type TblATTSearchPara struct {
	Topn int64
}

type TblATT struct {
	Id    int64
	Asset string
	TblATTContent
}

type TblATTData struct {
	//	Totality int64
	Topn     int64        `json:"topn"`
	Elements []TblATTList `json:"elements"`
}

type TblATTContent struct {
	ServerIp string `json:"serverIp"`
	Times    int64  `json:"times"`
}

type TblATTList struct {
	TblATTContent
}

/*************************/

/**********紧急事件**********/
/********获取紧急事件数（天）********
type TblUrgencySearchPara struct {
	Type      string
	LastCount int32
	KeyWord   string
	Page      int32
	Count     int32
	Sort      string
	Order     string
	Unit      string
	PField    TblPublicPara
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

/********获取新紧急事件分类数(天)********
type TblUgcClassifySearchPara struct {
	Type      string
	LastCount int32
	KeyWord   string
	Page      int32
	Count     int32
	Sort      string
	Order     string
	Unit      string
	PField    TblPublicPara
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
/*************威胁模型统计************
type TblUgcMCountSearchPara struct {
	Type      string
	LastCount int32
	KeyWord   string
	Page      int32
	Count     int32
	Sort      string
	Order     string
	Unit      string
	PField    TblPublicPara
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

/********获取紧急事件详情(天)********
//UgcD:"UrgencyDetails"
type TblUgcDSearchPara struct {
	Type      string
	LastCount int32
	KeyWord   string
	Page      int32
	Count     int32
	Sort      string
	Order     string
	PField    TblPublicPara
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

/*************************/

/**********所有安全事件**********/
/*http流量检测*/
/********获取http流量攻击数(天)********
//HFA:"HttpFlowAttackTrend"
type TblHFATSearchPara struct {
	Type   string
	PField TblPublicPara
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
}*
type TblHFATList struct {
	TblHFACRequest
	//TblHFAL
}

/*************************/

/********获取http流量攻击分类数(天)********
//HFAC:"HttpFlowAttackClassify"
type TblHFACSearchPara struct {
	PField TblPublicPara
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
	ScanningPprobe int64 `json:"scaningprobe"`
}
type HFACUniversal struct {
	AttDisclosure   int64 `json:"disclosure"`
	AttDos          int64 `json:"dos"`
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
	Other           int64 `json:"other"`
}
type HFACScaningProbe struct {
	AttReputScanner  int64 `json:"reputation_scanner"`
	AttReputSripting int64 `json:"reputation_scripting"`
	AttReputCrawler  int64 `json:"reputation_crawler"`
}
type TblHFAL struct {
	HFACUniversal
	HFACScaningProbe
}
type HFACContent struct {
	Time int64 `json:"time"`
	TblHFACRequest
}
type TblHFACList struct {
	HFACContent
}

/*************************/

/********获取http流量攻击详情(天)********
//HFAD:"HttpFlowAttackDetails"
type TblHFADSearchPara struct {
	Type      string
	LastCount int32
	KeyWord   string
	Page      int32
	Count     int32
	Sort      string
	Order     string
	PField    TblPublicPara
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

type TblHFADContent struct {
	Time        int64  `json:"time"`
	Client      string `json:"client"`
	HostName    string `json:"hostName"`
	Rev         string `json:"rev"`
	Message     string `json:"msg"`
	Attack      string `json:"attack"`
	Severity    int    `json:"severity"`
	Maturity    int    `json:"maturity"`
	Accuracy    int    `json:"accuracy"`
	Uri         string `json:"uri"`
	UniqueId    string `json:"unique_id"`
	Ref         string `json:"ref"`
	Tags        string `json:"tags"`
	RuleFile    string `json:"ruleFile"`
	RuleLine    int    `json:"ruleLine"`
	RuleId      int    `json:"ruleId"`
	RuleData    string `json:"ruleData"`
	RuleVersion string `json:"ruleVersion"`
	Version     string `json:"version"`
	//SrcPort     int    `json:"srcPort"`
	//DestPort    int    `json:"destPort"`
}
type TblHFADList struct {
	TblHFADContent
}

/*************************/

/**********文件威胁检测**********/
/********获取文件威胁数(天)********
//FTT:"FileThreatTrend"
type TblFTTSearchPara struct {
	Type   string
	PField TblPublicPara
}
type TblFTT struct {
	Id   int64
	Time int64
	TblFTL
}
type TblFTTData struct {
	Counts   int64        `json:"counts"`
	Elements []TblFTTList `json:"elements"`
}

/*
type TblFileThreatContent struct {
	Date  int64 `json:"date"`
	Times int64 `json:"times"`
}*
type TblFTTList struct {
	TblFTCRequest
	//TblFTL
}

/*************************/

/********获取文件威胁分类数(天)********
//FTC:"FileThreatClassify"
type TblFTCSearchPara struct {
	PField TblPublicPara
}

type TblFTC struct {
	Id int64
	TblFTCContent
}

type TblFTCData struct {
	Counts   int64        `json:"counts"`
	Classify []TblFTCList `json:"classify"`
}

type TblFTCContent struct {
	Time int64 `json:"time"`
	TblFTL
}
type TblFTCRequest struct {
	FTCUniversal
}
type FTCUniversal struct {
	BackDoor int `json:"backdoor"`
	Trojan   int `json:"trojan"`
	Spyware  int `json:"spyware"`
	Malware  int `json:"malware"`
	Virus    int `json:"virus"`
	Worm     int `json:"worm"`
	HackTool int `json:"hacktool"`
	Exploit  int `json:"exploit"`
}
type TblFTL struct {
	FTCUniversal
	RiskTool int `json:"risktool"`
	Joke     int `json:"joke"`
	Adware   int `json:"adware"`
	Other    int `json:"other"`
}
type FTCContent struct {
	Time int64 `json:"time"`
	TblFTCRequest
}
type TblFTCList struct {
	FTCContent
}

/*************************/

/********获取文件威胁详情(天)********
//FTD:"FileThreatDetails"
type TblFTDSearchPara struct {
	Type      string
	LastCount int32
	KeyWord   string
	Page      int32
	Count     int32
	Sort      string
	Order     string
	PField    TblPublicPara
	OfflineTage
}
type TblFTDData struct {
	Totality int64        `json:"total"`
	Counts   int64        `json:"counts"`
	Elements []TblFTDList `json:"elements"`
}

type TblFTD struct {
	Id int64
	TblFTDContent
	DupIdList
}

/*
type TblFTDContent struct {
	HackerIp     string `json:"hackerIp"`
	HackerRegion string `json:"hackerRegion"`
	ServerIp     string `json:"serverIp"`
	ServerRegion string `json:"serverRegion"`
	PlatformInfo string `json:"platformInfo"`
	AttackType   string `json:"attackType"`
	AttackName   string `json:"attackName"`
	AttackSample string `json:"attackSample"`
	RiskLevel    string `json:"riskLevel"`
	Time         int64  `json:"time"`
}
*
type TblFTDContent struct {
	ThreatName      string `json:"threatName"`
	SubFile         string `json:"subFile"`
	LocalThreatName string `json:"localThreatName"`
	LocalVType      string `json:"localVType"`
	LocalPlatfrom   string `json:"localPlatfrom"`
	LocalVName      string `json:"localVName"`
	LocalExtent     string `json:"localExtent"`
	LocalEngineType string `json:"localEngineType"`
	LocalLogType    string `json:"localLogType"`
	LocalEngineIP   string `json:"localEngineIP"`
	LogTime         int64  `json:"time"`
	XDR_VDSList
}
type TblFTDList struct {
	TblFTDContent
}

type XDR_VDSList struct {
	SrcIp    string `json:"srcIp"`
	SrcPort  int    `json:"srcPort"`
	DestIp   string `json:"destIp"`
	DestPort int    `json:"destPort"`
	AppFile  string `json:"filePath"`
	HttpUrl  string `json:"url"`
}

/*************************/

/**********恶意流量检测**********/
/********获取恶意流量数(天)********
//MFT:"MaliceFlowTrend"
type TblMFTSearchPara struct {
	Type   string
	PField TblPublicPara
}
type TblMFT struct {
	Id   int64
	Time int64
	TblMFL
	Other int64
}
type TblMFTData struct {
	Counts   int64        `json:"counts"`
	Elements []TblMFTList `json:"elements"`
}

/*
type TblMalFlowContent struct {
	Time  int64 `json:"time"`
	Times int64 `json:"times"`
}*
type TblMFTList struct {
	TblMFL
}

/*************************/

/********获取恶意流量分类数(天)********
//MFC:"MaliceFlowClassify"
type TblMFCSearchPara struct {
	PField TblPublicPara
}

type TblMFC struct {
	Id int64
	TblMFCContent
}

type TblMFCData struct {
	Counts   int64        `json:"counts"`
	Classify []TblMFCList `json:"classify"`
}

type TblMFCContent struct {
	Time int64 `json:"time"`
	TblMFL
	Other int64 `json:"Other"`
}
type TblMFL struct {
	PrivilegeGain     int64 `json:"privilege_gain"`
	DDos              int64 `json:"ddos"`
	InformationLeak   int64 `json:"information_leak"`
	WebAttack         int64 `json:"web_attack"`
	ApplicationAttack int64 `json:"application_attack"`
	CandC             int64 `json:"candc"`
	Malware           int64 `json:"malware"`
	MiscAttack        int64 `json:"misc_attack"`
}

/*
type TblMFL struct {
	AttemptedAdmin   int   `json:"attempted_admin"`
	AttemptedUser    int   `json:"attempted_user"`
	InappropContent  int   `json:"inappropriate_content"` //InappropriateContent
	PolicyViolation  int   `json:"policy_violation"`
	ShellcodeDetect  int   `json:"shellcode_detect"`
	SuccessfulAdmin  int   `json:"successful_admin"`
	SuccessfulUser   int   `json:"successful_user"`
	TrojanActivity   int   `json:"trojan_activity"`
	UnsuccessfulUser int   `json:"unsuccessful_user"`
	WebAppAttack     int   `json:"web_application_attack"` //WebApplicationAttack
	AttemptedDos     int   `json:"attempted_dos"`
	AttemptedRecon   int   `json:"attempted_recon"`
	BadUnknown       int   `json:"bad_unknown"`
	DefLoginAttempt  int   `json:"default_login_attempt"` //DefaultLoginAttempt
	DenialOfService  int   `json:"denial_of_service"`
	MiscAttack       int   `json:"misc_attack"`
	NonStanProto     int   `json:"non_standard_protocol"` //NonStandardProtocol
	RpcPortmapDecode int   `json:"rpc_portmap_decode"`
	SuccessfulDos    int   `json:"successful_dos"`
	SucfReconLarg    int   `json:"successful_recon_largescale"`      //SuccessfulReconLargescale
	SucfReconLim     int   `json:"successful_recon_limited"`       //SuccessfulReconLimited
	SuspFileDetect   int   `json:"suspicious_filename_detect"` //SuspiciousFilenameDetect
	SuspiciousLogin  int   `json:"suspicious_login"`
	SystemCallDetect int   `json:"system_call_detect"`
	UCPC             int   `json:"unusual_client_port_connection"`     //UnusualClientPortConnection
	WebAppActivity   int   `json:"web_application_activity"` //WebApplicationActivity
	IcmpEvent        int   `json:"icmp_event"`
	MiscActivity     int   `json:"misc_activity"`
	NetworkScan      int   `json:"network_scan"`
	NotSuspicious    int   `json:"not_suspicious"`
	ProtoCmdDecode   int   `json:"protocol_command_decode"` //ProtocolCommandDecode
	StringDetect     int   `json:"string_detect"`
	Unknown          int   `json:"unknown"`
	TcpConnection    int   `json:"tcp_connection"`
	Other            int   `json:"other"`
}
*
type TblMFCList struct {
	TblMFCContent
}

/*************************/

/********获取恶意流量详情(天)********
//MFD:"MaliceFlowDetails"
type TblMFDSearchPara struct {
	Type      string
	LastCount int32
	KeyWord   string
	Page      int32
	Count     int32
	Sort      string
	Order     string
	PField    TblPublicPara
}
type TblMFDData struct {
	Totality int64        `json:"total"`
	Counts   int64        `json:"counts"`
	Elements []TblMFDList `json:"elements"`
}

type TblMFD struct {
	Id int64
	TblMFDContent
}

/*原定结构
type TblMFDContent struct {
	HackerIp     string `json:"hackerIp"`
	HackerRegion string `json:"hackerRegion"`
	ServerIp     string `json:"serverIp"`
	ServerRegion string `json:"serverRegion"`
	PlatformInfo string `json:"platformInfo"`
	AttackType   string `json:"attackType"`
	AttackName   string `json:"attackName"`
	AttackSample string `json:"attackSample"`
	RiskLevel    string `json:"riskLevel"`
	Time         int64  `json:"time"`
}
*
type TblMFDContent struct {
	Time       int64  `json:"time"`
	SrcIp      string `json:"srcIp"`
	SrcPort    string `json:"srcPort"`
	DestIp     string `json:"destIp"`
	DestPort   string `json:"destPort"`
	Proto      string `json:"protocol"`
	AttackType string `json:"attackType"`
	Details    string `json:"details"`
	Severity   int    `json:"severity"`
	Engine     string `json:"engine"`
	ByzoroType string `json:"byzoroType"`
}
type TblMFDList struct {
	TblMFDContent
}

/*************************/

/**********暴力破解**********/
/********获取暴力破解数(天)********
type TblBruteForcePara struct {
	Type   string
	PField TblPublicPara
}
type TblBruteForce struct {
	Id   int64
	Type string
	TblBruteForceContent
}
type TblBruteForceData struct {
	Counts   int64               `json:"counts"`
	Elements []TblBruteForceList `json:"elements"`
}
type TblBruteForceContent struct {
	Time  int64 `json:"time"`
	Times int64 `json:"times"`
}
type TblBruteForceList struct {
	TblBruteForceContent
}

/*************************/

/********获取暴力破解分类数(天)********/
//BFC:"BruteForceClassify"
//??
/*************************/

/********获取暴力破解详情(天)********
//BFD:"BruteForceDetails"
type TblBFDSearchPara struct {
	Type      string
	LastCount int32
	KeyWord   string
	Page      int32
	Count     int32
	Sort      string
	Order     string
	PField    TblPublicPara
}
type TblBFDData struct {
	Counts   int64        `json:"counts"`
	Totality int64        `json:"total"`
	Elements []TblBFDList `json:"elements"`
}

type TblBFD struct {
	Id int64
	TblBFDContent
}

type TblBFDContent struct {
	Ip    string `json:"ip"`
	Port  int32  `json:"port"`
	Time  string `json:"time"`
	Count int32  `json:"count"`
	Name  string `json:"name"`
	Level string `json:"level"`
}

type TblBFDList struct {
	TblBFDContent
}

/*************************/

/**********扫描时间**********/
/********获取扫描事件数(天)********
type TblScanEventPara struct {
	Type   string
	PField TblPublicPara
}
type TblScanEvent struct {
	Id   int64
	Type string
	TblScanEventContent
}
type TblScanEventData struct {
	Counts   int64              `json:"counts"`
	Elements []TblScanEventList `json:"elements"`
}
type TblScanEventContent struct {
	Time  int64 `json:"time"`
	Times int64 `json:"times"`
}
type TblScanEventList struct {
	TblScanEventContent
}

/*************************/

/********获取扫描事件详情(天)********
//SED:"ScanEventDetails"
type TblSEDSearchPara struct {
	Type      string
	LastCount int32
	KeyWord   string
	Page      int32
	Count     int32
	Sort      string
	Order     string
	PField    TblPublicPara
}
type TblSEDData struct {
	Counts   int64        `json:"counts"`
	Totality int64        `json:"total"`
	Elements []TblSEDList `json:"elements"`
}

type TblSED struct {
	Id int64
	TblSEDContent
}

type TblSEDContent struct {
	Time      int64  `json:"time"`
	Conntype  string `json:"conntype"`
	Host      string `json:"sip"`
	AlertType string `json:"alerttype"`
	Count     int64  `json:"count"`
}

type TblSEDList struct {
	TblSEDContent
}

/*************************/
/********获取攻击数（天）********
type TblAttackCountSearchPara struct {
	Type      string
	LastCount int32
	KeyWord   string
	Page      int32
	Count     int32
	Sort      string
	Order     string
	Unit      string
	PField    TblPublicPara
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

/*************************/
/********DNS统计(天)********
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

/*************************/
/********SSL证书统计(天)********
//SSLCS:"SecuritySocketLayerCertificateStatistics"
type TblSSLCSSearchPara struct {
	Type      string
	LastCount int32
	KeyWord   string
	Page      int32
	Count     int32
	Sort      string
	Order     string
	Unit      string
	Verify    int
	Ip        string
	ComnName  string
	UnitName  string
	SerialNum string
	PField    TblPublicPara
}
type TblSSLCSPublic struct {
	Time             string `json:"time"`
	SSLCertComnName  string `json:"s_cert_comnname"`
	SSLCertUnitName  string `json:"s_cert_unitname"`
	SSLCertSerialNum string `json:"s_cert_serialnum"`
	Count            int64  `json:"count"`
}

type TblSSLCSData struct {
	Totality int64          `json:"total"`
	Counts   int64          `json:"counts"`
	Elements []TblSSLCSList `json:"elements"`
}
type TblSSLCSIpData struct {
	Totality int64            `json:"total"`
	Counts   int64            `json:"counts"`
	Elements []TblSSLCSIpList `json:"elements"`
}
type TblSSLCSCertData struct {
	Totality int64              `json:"total"`
	Counts   int64              `json:"counts"`
	Elements []TblSSLCSCertList `json:"elements"`
}

type TblSSLCS struct {
	//Id int64
	SSLVerify int
	TblSSLCSContent
}
type TblSSLCSIp struct {
	//Id int64
	SSLVerify int
	TblSSLCSIpContent
}
type TblSSLCSCert struct {
	//Id int64
	SSLVerify int
	TblSSLCSCertContent
}

type TblSSLCSContent struct {
	Ip string `json:"ip"`
	TblSSLCSPublic
}
type TblSSLCSIpContent struct {
	TblSSLCSPublic
}
type TblSSLCSCertContent struct {
	TblSSLCSPublic
}

type TblSSLCSList struct {
	TblSSLCSContent
}
type TblSSLCSIpList struct {
	TblSSLCSIpContent
}
type TblSSLCSCertList struct {
	TblSSLCSCertContent
}

/*************************/
/**********数据迁移任务**********
//DCT:DataCopyTask
type DCTOperatePara struct {
	Cmd      string
	Location string
	Start    string
	End      string
	Date     int64
	DiskPath string
}

type DCTOperateResult struct {
	Result string `json:"result"`
	//Location string `json:"location"`
	//Date     int64  `json:"date"`
}

/******************************/
/**********数据迁移任务列表**********
//DCTL:DataCopyTaskList
type TblDCTLSearchPara struct {
	Location  string
	Date      int64
	Status    string
	Sort      string
	Order     string
	Page      int32
	Count     int32
	LastCount int32
	PField    TblPublicPara
}

type TblDCT struct {
	Id int
	DCTLContent
}

type TblDCTLData struct {
	Totality int64 `json:"total"`
	//Counts   int64        `json:"counts"`
	Elements []TblDCTList `json:"elements"`
}

type DCTLContent struct {
	Location  string `json:"location"`
	DataStart string `json:"start"`
	DataEnd   string `json:"end"`
	Date      int64  `json:"date"`
	DiskPath  string `json:"diskpath"`
	Status    string `json:"status"`
	Details   string `json:"details"`
}

type TblDCTList struct {
	DCTLContent
}

/******************************/
/**********数据迁移任务状态**********
//DCTS:DataCopyTaskStatus
type DCTSPara struct {
	Cmd      string
	Location string
	Date     int64
	Rate     int
	Status   string
	Details  string
}

type DCTSContent struct {
	Location       string `json:"location"`
	Date           int64  `json:"date"`
	Status         string `json:"status"`
	RateOfProgress int    `json:"rate"`
	Details        string `json:"details"`
}

type DCTSList struct {
	DCTSContent
}

type UpdateDCTSResult struct {
	Result string `json:"result"`
}

type DCTSLastTime struct {
	LastTime string `json:"lasttime"`
}

/******************************/
