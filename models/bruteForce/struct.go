/********获取暴力破解详情(天)********/
//BFD:"BruteForceDetails"
package bruteForce

/********共有结构********/
type TblPublicPara struct {
	Start int64
	End   int64
}
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
