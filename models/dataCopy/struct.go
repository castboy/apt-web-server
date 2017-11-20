package dataCopy

import (
	"apt-web-server_v2/models/modelsPublic"
)

/********共有结构********
type TblPublicPara struct {
	Start int64
	End   int64
}
*/
/**********数据迁移任务**********/
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
/**********数据迁移任务列表**********/
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
	PField    modelsPublic.TblPublicPara
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
/**********数据迁移任务状态**********/
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
