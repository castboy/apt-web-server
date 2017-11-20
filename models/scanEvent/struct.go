/**********扫描事件**********/
package scanEvent

import (
	"apt-web-server_v2/models/modelsPublic"
)

/********共有结构********
type TblPublicPara struct {
	Start int64
	End   int64
}
*/
/********获取扫描事件数(天)********/
type TblScanEventPara struct {
	Type   string
	PField modelsPublic.TblPublicPara
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

/********获取扫描事件详情(天)********/
//SED:"ScanEventDetails"
type TblSEDSearchPara struct {
	Type      string
	LastCount int32
	KeyWord   string
	Page      int32
	Count     int32
	Sort      string
	Order     string
	PField    modelsPublic.TblPublicPara
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
