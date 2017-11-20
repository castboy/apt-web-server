package offlineAssignment

import (
	"apt-web-server_v2/models/modelsPublic"
)

/********共有结构********
type TblPublicPara struct {
	Start int64
	End   int64
}
*/
/************离线任务存储结构*************/
//OLA:"Off_Line Assignment"
type TblOLASearchPara struct {
	Cmmand     string `json:"cmd"`
	Name       string `json:"name"`
	Time       int64  `json:"time"`
	Type       string `json:"type"`
	Start      string `json:"start"`
	End        string `json:"end"`
	Weight     int    `json:"weight"`
	OfflineTag string
	Rule       string `json:"rule"`
	//RuleSets   []TblRuleSdSet `json:"ruleset"` //规则获取方式变更，原来是字符串，现在是规则名，由规则名找到对应的规则，然后进行后续处理;存放的是以,分割的rule的id，如 1,2,111。
	RuleSet string `json:"ruleset"` //规则获取方式变更，原来是字符串，现在是规则名，由规则名找到对应的规则，然后进行后续处理;存放的是以,分割的rule的id，如 1,2,111。
	Details string `json:"details"`
	//	Status string
}

type TblRuleSdSet struct {
	Id   int64  `json:"id"`
	Rule string `json:"rule"`
}

type BasePara struct {
	Cmmand string
	TblOLAPrimaryData
}

type OLACreatPara struct {
	BasePara
}
type OLADeletePara struct {
	BasePara
}
type OLAStartPara struct {
	BasePara
}
type OLAStopPara struct {
	BasePara
}
type TblOLA struct {
	Id int
	TblOLAPrimaryData
	TblOLAParameterData
	TblOLAStatusData
	Topic   string
	Details string
}

type TblOLAPrimaryData struct {
	Name string `json:"name"`
	Time int64  `json:"time"`
}
type TblOLAParameterData struct {
	Type   string `json:"type"`
	Start  string `json:"start"`
	End    string `json:"end"`
	Weight int    `json:"weight"`
}
type TblOLAStatusData struct {
	Status string `json:"status"`
}
type ReadStatus struct {
}
type CMDResult struct {
	Result string `json:"result"`
}

type AgentPara struct {
	Engine     string `json:"Engine"`
	Topic      string `json:"Topic"`
	Weight     int    `json:"Weight"`
	SignalType string `json:"SignalType"`
}
type ETCDPicker struct {
	Type   string
	Start  string
	End    string
	topic  string
	State  string
	Offset string
	Count  int
	Total  int
	Now    int64
}

type ETCDAgent struct {
	First  int
	Engine int
	Err    int
	Cache  int
	Last   int
	Weight int
}
type TaskStatus struct {
	Status      string `json:"status"`
	PickerCount int    `json:"pickerCount"`
	PickerTotal int    `json:"pickerTotal"`
	AgentCount  int    `json:"agentCount"`
	AgentTotal  int    `json:"agentTotal"`
}

//get task list
type TblTaskSearchPara struct {
	Name       string
	Time       int64
	Type       string
	Status     string
	Sort       string
	Order      string
	Page       int32
	Count      int32
	LastCount  int32
	PField     modelsPublic.TblPublicPara
	OfflineTag string
}

type TblTaskDetails struct {
	Id int64
	TaskContent
}
type TblTaskData struct {
	Totality int64         `json:"total"`
	Counts   int64         `json:"counts"`
	Elements []TblTaskList `json:"elements"`
}

type TaskContent struct {
	Name       string `json:"name"`
	RuleSet    string `json:"ruleset"`
	CreateTime int64  `json:"time"`
	Type       string `json:"type"`
	DataStart  string `json:"start"`
	DataEnd    string `json:"end"`
	Weight     int    `json:"weight"`
	Topic      string `json:"topic"`
	Status     string `json:"status"`
	Details    string `json:"details"`
}

type TaskContent2UI struct {
	Name string `json:"name"`
	//RuleSets   []TblRuleSdSet `json:"ruleset"`
	RuleSets   string `json:"ruleset"`
	CreateTime int64  `json:"time"`
	Type       string `json:"type"`
	DataStart  string `json:"start"`
	DataEnd    string `json:"end"`
	Weight     int    `json:"weight"`
	Topic      string `json:"topic"`
	Status     string `json:"status"`
	Details    string `json:"details"`
}

type TblTaskList struct {
	//TaskContent
	TaskContent2UI
}

//get task status struct
type TaskListPara struct {
	TaskList   string
	OfflineTag string
}
type TaskStatusRequestPara struct {
	Name string
	Time int64
}

type TaskEtcdResData struct {
	PickerCount float32 `json:"pickerCount"`
	PickerTotal float32 `json:"pickerTotal"`
	AgentCount  float32 `json:"agentCount"`
	AgentTotal  float32 `json:"agentTotal"`
}

type TaskStatusContent struct {
	Status        string  `json:"status"`
	PickerPercent float32 `json:"pickerPercent"`
	AgentPercent  float32 `json:"agentPercent"`
	TaskStatusRequestPara
}
type TaskStatusList struct {
	TaskStatusContent
}

/************自定义规则的存储结构*************/
type TblRuleOperPara struct {
	Command string `json:"cmd"`   //add,del,mod
	Alias   string `json:"alias"` //add,del,mod
	Id      int64  `json:"id"`
	IdSet   string `json:"idset"` //string format, store the ids of rule to be delete, such as "1,200,3000"
	//VarSet   []StVarset `json:"varsets"`
	VarSet    string `json:"varset"`
	Oper      string `json:"oper"`
	OperInfo  string `json:"Operinfo"`
	TransFunc string `json:"tfunc"`
	Phase     string `json:"phase"` // string of number , such as "1","5"
	Severity  string `json:"severity"`
	Accuracy  string `json:"accuracy"` // string of number , such as "1","5"
	Maturity  string `json:"maturity"` // string of number , such as "1","5"
	Tag       string `json:"tag"`
	Details   string `json:"details"`
}

type StVarset struct {
	Var     string `json:"var"`
	VarInfo string `json:"varinfo"`
}

/*************************/
/********自定义规则 rule 信息查询********/
type TblRuleSdSearchPara struct {
	Type      string //
	Alias     string // 规则别名
	Id        int64
	RuleSet   string //存放rule id集合的字符串，用，分割，比如  "1,2,100"
	LastCount int32
	Page      int32
	Count     int32
}

type TblRuleSdLstData struct {
	Totality int64                 `json:"total"`  //total num searched
	Counts   int64                 `json:"counts"` //num returned to client for this time
	Elements []TblRuleSdLstContent `json:"elements"`
}

type TblRuleSdLstContent struct {
	TblRuleSdLstPublic
}

type TblRuleSdLstPublic struct {
	Id    int32  `json:"id"`
	Alias string `json:"alias"` //add,del,mod
	//VarSet   []TblRuleSdLstVarSet `json:"varsets"`
	VarSet    string `json:"varset"`
	Oper      string `json:"oper"`
	OperInfo  string `json:"Operinfo"`
	TransFunc string `json:"tfunc"`
	Phase     string `json:"phase"` // string of number , such as "1","5"
	Severity  string `json:"severity"`
	Accuracy  string `json:"accuracy"` // string of number , such as "1","5"
	Maturity  string `json:"maturity"` // string of number , such as "1","5"
	Tag       string `json:"tag"`
	Details   string `json:"details"`
}

type TblRuleSdLstVarSet struct {
	Var     string `json:"var"`
	VarInfo string `json:"varinfo"`
}

type TblRuleSdLstSame2Tbl struct {
	Id     int32
	Alias  string `json:"alias"` //add,del,mod
	VarSet string `json:"varset"`
	//Var      string
	//VarInfo  string
	//Var2     string
	//VarInfo2 string
	//Var3     string
	//VarInfo3 string
	//Var4     string
	//VarInfo4 string
	//Var5     string
	//VarInfo5 string
	Oper      string `json:"oper"`
	OperInfo  string `json:"Operinfo"`
	TransFunc string `json:"tfunc"`
	Phase     string `json:"phase"` // string of number , such as "1","5"
	Severity  string `json:"severity"`
	Accuracy  string `json:"accuracy"` // string of number , such as "1","5"
	Maturity  string `json:"maturity"` // string of number , such as "1","5"
	Tag       string `json:"tag"`
	Details   string `json:"details"`
}
