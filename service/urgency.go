package service

type Urgency struct {
	Time        int64
	Type        string
	ServerIp    string
	Description string
}

func (this *Urgency) List(para *UgcListPara) (error, []Urgency) {
	if para.LastCount 
}

type UgcList struct {
}

type UgcListPara struct {
	Start     int64
	End       int64
	Type      string
	LastCount int32
	KeyWord   string
	Page      int32
	Count     int32
}

