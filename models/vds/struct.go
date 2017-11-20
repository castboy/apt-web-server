/**********文件威胁检测**********/
package vds

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

/********获取文件威胁数(天)********/
//FTT:"FileThreatTrend"
type TblFTTSearchPara struct {
	Type   string
	PField modelsPublic.TblPublicPara
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
}*/
type TblFTTList struct {
	TblFTCRequest
	//TblFTL
}

/*************************/

/********获取文件威胁分类数(天)********/
//FTC:"FileThreatClassify"
type TblFTCSearchPara struct {
	PField modelsPublic.TblPublicPara
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

/********获取文件威胁详情(天)********/
//FTD:"FileThreatDetails"
type TblFTDSearchPara struct {
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
type TblFTDData struct {
	Totality int64        `json:"total"`
	Counts   int64        `json:"counts"`
	Elements []TblFTDList `json:"elements"`
}
type DupIdList struct {
	TblId  int
	TaskId int
}
type TblFTD struct {
	Id int64
	TblFTDContent
	DupIdList
}

type TblFTDContent struct {
	modelsPublic.AreaMsg
	Time            int64  `json:"Time"`
	LocalVType      string `json:"Attack"`
	LocalVName      string `json:"Local_vname"`
	LocalLogType    string `json:"Local_logtype"`
	ThreatName      string `json:"Threatname"`
	SubFile         string `json:"Subfile"`
	LocalThreatName string `json:"Local_threatname"`
	LocalPlatfrom   string `json:"Local_platfrom"`
	LocalExtent     string `json:"Severity"`
	LocalEngineType string `json:"Local_enginetype"`
	LocalEngineIP   string `json:"Local_engineip"`
	AppFile         string `json:"FilePath"`
	HttpUrl         string `json:"HttpUrl"`
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
