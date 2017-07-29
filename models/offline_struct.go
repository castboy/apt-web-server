package models

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
	Details    string `json:"details"`
	//	Status string
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
	OfflineTag string `json:"OfflineTag"`
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
	PField     TblPublicPara
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
	TaskContent
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
