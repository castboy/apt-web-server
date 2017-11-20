/********白名單 whitelist 信息查询********/
package whiteList

import (
	"apt-web-server_v2/models/modelsPublic"
)

//SSLCLst:"SecuritySocketLayerCertificateList"
type CMDResult struct {
	Result string `json:"result"`
}

var MapProto2Num = map[string]int{
	"TCP": 6,
	"UDP": 17}

var MapNum2Proto = map[int32]string{
	6:  "TCP",
	17: "UDP"}

type TblWLSearchPara struct {
	Type      string // ( all | exact )
	Sip       string
	Sport     int32
	Dip       string
	Dport     int32
	Proto     int32
	LastCount int32
	Page      int32
	Count     int32
	PField    modelsPublic.TblPublicPara
}

type TblWLLstData struct {
	Totality int64             `json:"total"`  //total num searched
	Counts   int64             `json:"counts"` //num returned to client for this time
	Elements []TblWLLstContent `json:"elements"`
}

type TblWLLstInfoTmp struct {
	TblWLLstContent
}

type TblWLLstContent struct {
	TblWLLstPublic
}

type TblWLLstPublic struct {
	//WLId    int64  `json:"id"`
	WLSip   string `json:"sip"`
	WLSport int32  `json:"sport"`
	WLDip   string `json:"dip"`
	WLDport int32  `json:"dport"`
	WLProto string `json:"proto"` /*int32*/
}

type TblWLLstList struct {
	TblWLLstContent
}

type TblWLLst struct {
	//Id int64
	//SSLVerify int
	//TblWLLstContent

	//WLId    int64  `json:"id"`
	WlSip   string `json:"sip"`
	WlSport int32  `json:"sport"`
	WlDip   string `json:"dip"`
	WlDport int32  `json:"dport"`
	WlProto int32  `json:"proto"`
}

/*************************/
/******** whitelist operate / 對白名單的各種操作 ********/
//WLOperate:"WhiteListOperate"

//WL:"White_List"
type TblWLOperatePara struct {
	Cmmand      string             `json:"cmd"` //add|delete|clear
	OpNum       int32              `json:"num"`
	WLOpElement []TblWLOperElement `json:"list"`
}

type TblWLOperateParaIn struct {
	Cmmand      string               `json:"cmd"` //add|delete|clear
	OpNum       int32                `json:"num"`
	WLOpElement []TblWLOperElementIn `json:"list"`
}

//type TblWLOperElement struct {
//	Sip   string `json:"sip"`
//	Sport int32  `json:"sport"`
//	Dip   string `json:"dip"`
//	Dport int32  `json:"dport"`
//	Proto int32  `json:"proto"`
//}

type TblWLOperElement struct {
	Sip   string `json:"sip"`
	Sport int32  `json:"sport"`
	Dip   string `json:"dip"`
	Dport int32  `json:"dport"`
	Proto string `json:"proto"` //Proto int32  `json:"proto"`
}

type TblWLOperElementIn struct {
	Sip   string `json:"sip"`
	Sport int32  `json:"sport"`
	Dip   string `json:"dip"`
	Dport int32  `json:"dport"`
	Proto int32  `json:"proto"` //Proto int32  `json:"proto"`
}

/*************************/
