/********获取网络流量数********/
package netFlow

import (
	"apt-web-server_v2/models/modelsPublic"
)

/********共有结构********
type TblPublicPara struct {
	Start int64
	End   int64
}
*/
//curl参数
type TblNetFlowSearchPara struct {
	Direction string
	AssetIP   string
	Protocol  string
	Unit      string
	PField    modelsPublic.TblPublicPara
}
type TblNetFlow struct {
	Id        int64
	Direction string
	AssetIP   string
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

/**********分表数据**********/
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

/******************************/
