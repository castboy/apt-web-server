/********SSL证书统计(天)********/
package ssl

import (
	"apt-web-server_v2/models/modelsPublic"
)

/********共有结构********
type TblPublicPara struct {
	Start int64
	End   int64
}*/

/********SSL证书详情查询********/ /*子元素結構共用SSL證書統計的結構*/
//SSLCLst:"SecuritySocketLayerCertificateList"
type TblSSLCLstSearchPara struct {
	Type      string //summary , comnname
	LastCount int32
	KeyWord   string
	Page      int32
	Count     int32
	Sort      string
	Order     string
	Unit      string
	Verify    int
	Ip        string
	ComnName  string
	UnitName  string
	Summary   string
	SerialNum string
	PField    modelsPublic.TblPublicPara
}

type TblSSLCLstData struct {
	Totality int64               `json:"total"`  //total num searched
	Counts   int64               `json:"counts"` //num returned to client this time
	Elements []TblSSLCLstContent `json:"elements"`
}

type TblSSLCLstContent struct {
	TblSSLCLstPublic
}

type TblSSLCLstPublic struct {
	SSLCertComnName  string `json:"s_cert_comnname"`
	SSLCertOrigName  string `json:"s_cert_origname"`
	SSLCertUnitName  string `json:"s_cert_unitname"`
	SSLCertSerialNum string `json:"s_cert_serialnum"`
	SSLCertNotBefore string `json:"s_cert_notbefore"`
	SSLCertNotAfter  string `json:"s_cert_notafter"`
	SSLCertVersion   string `json:"s_cert_version"`
}

type TblSSLCLstList struct {
	TblSSLCLstContent
}

type TblSSLCLst struct {
	//Id int64
	//SSLVerify int
	TblSSLCLstContent
}

/*************************/
//SSLCS:"SecuritySocketLayerCertificateStatistics"
type TblSSLCSSearchPara struct {
	Type      string
	LastCount int32
	KeyWord   string
	Page      int32
	Count     int32
	Sort      string
	Order     string
	Unit      string
	Verify    int
	Ip        string
	ComnName  string
	UnitName  string
	SerialNum string
	PField    modelsPublic.TblPublicPara
}
type TblSSLCSPublic struct {
	Time             string `json:"time"`
	SSLCertComnName  string `json:"s_cert_comnname"`
	SSLCertUnitName  string `json:"s_cert_unitname"`
	SSLCertSerialNum string `json:"s_cert_serialnum"`
	Count            int64  `json:"count"`
}

type TblSSLCSData struct {
	Totality int64          `json:"total"`
	Counts   int64          `json:"counts"`
	Elements []TblSSLCSList `json:"elements"`
}
type TblSSLCSIpData struct {
	Totality int64            `json:"total"`
	Counts   int64            `json:"counts"`
	Elements []TblSSLCSIpList `json:"elements"`
}
type TblSSLCSCertData struct {
	Totality int64              `json:"total"`
	Counts   int64              `json:"counts"`
	Elements []TblSSLCSCertList `json:"elements"`
}

type TblSSLCS struct {
	//Id int64
	SSLVerify int
	TblSSLCSContent
}
type TblSSLCSIp struct {
	//Id int64
	SSLVerify int
	TblSSLCSIpContent
}
type TblSSLCSCert struct {
	//Id int64
	SSLVerify int
	TblSSLCSCertContent
}

type TblSSLCSContent struct {
	Ip string `json:"ip"`
	TblSSLCSPublic
}
type TblSSLCSIpContent struct {
	TblSSLCSPublic
}
type TblSSLCSCertContent struct {
	TblSSLCSPublic
}

type TblSSLCSList struct {
	TblSSLCSContent
}
type TblSSLCSIpList struct {
	TblSSLCSIpContent
}
type TblSSLCSCertList struct {
	TblSSLCSCertContent
}
