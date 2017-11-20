/**********恶意流量检测**********/
package ids

import (
	"apt-web-server_v2/models/modelsPublic"
)

/********共有结构********
type TblPublicPara struct {
	Start int64
	End   int64
}
*/
/********获取恶意流量数(天)********/
//MFT:"MaliceFlowTrend"
type TblMFTSearchPara struct {
	Type   string
	PField modelsPublic.TblPublicPara
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
}*/
type TblMFTList struct {
	TblMFL
}

/*************************/

/********获取恶意流量分类数(天)********/
//MFC:"MaliceFlowClassify"
type TblMFCSearchPara struct {
	PField modelsPublic.TblPublicPara
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
*/
type TblMFCList struct {
	TblMFCContent
}

/*************************/

/********获取恶意流量详情(天)********/
//MFD:"MaliceFlowDetails"
type TblMFDSearchPara struct {
	Type      string
	LastCount int32
	KeyWord   string
	Id        int32
	Page      int32
	Count     int32
	Sort      string
	Order     string
	PField    modelsPublic.TblPublicPara
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

type TblMFDContent struct {
	modelsPublic.AreaMsg
	Time       int64  `json:"Time"`
	ByzoroType string `json:"Attack"`
	AttackType string `json:"Attack_type"`
	Details    string `json:"Details"`
	Severity   int    `json:"Severity"`
	Engine     string `json:"Engine"`
}
type TblMFDList struct {
	TblMFDContent
}

/*************************/
