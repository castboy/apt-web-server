// Sniffplcy_st.go
package sniff

type Sniffplcy_st struct {
	Plcy_name         string
	Attack_ip         string
	Victim_ip         string
	Dst_port          string
	Proto             string
	Affect_time_start int
	Affect_time_end   int
	Plcy_status       int
	Issued_status     int
	Plcy_date         int
	Pcap_plcy         string
	Pcap_path         string
}
type Sniffselect_st struct {
	Type    string
	Page    int
	Count   int
	Lies    string
	Orderby string
}
type Sniffshow_st struct {
	Count  int            `json:"counts"`
	Plcy_s []Sniffplcy_st `json:"elements"`
}
type Sniffissued_st struct {
	Plcy_name [10]string
}
type Msg_st struct {
	Status int    `json:"status"`
	Log    string `json:"msg"`
}
type Sniffstatus_st struct {
	Type string
	Msg  Msg_st
}
