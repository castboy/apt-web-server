// Report.go
package report

type Attribute_st struct {
	Info []string
}
type Security_st struct {
	Score int `json:"score"`
	//Statistics_1 []int           `json:"statistics1"`
	Statistics []Statistics_st `json:"statistics"`
	Event      []Event_st      `json:"event"`
}
type Statistics_st struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}
type Event_st struct {
	Type string `json:"type"`
	Ip   string `json:"ip"`
}
