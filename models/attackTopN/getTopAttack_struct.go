package attackTopN

type AttackSearchPara struct {
	Type  string
	Start int64
	End   int64
	Unit  string
	Count int
}

type Attack struct {
	Id    int64
	Total int
	AttackContent
}

type AttackData struct {
	//Totality int64        `json:"total"`
	Counts   int          `json:"counts"`
	Elements []AttackList `json:"elements"`
}

type AttackContent struct {
	Type string `json:"name"`
	//Total int           `json:"tatal"`
	Data []AttackCount `json:"data"`
}
type AttackCount struct {
	Time  string `json:"name"`
	Times int64  `json:"y"`
}

type AttackList struct {
	AttackContent
}
