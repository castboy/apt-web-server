package main

import (
	"apt-web-server_v2/controllers/attackTopN"
	"apt-web-server_v2/controllers/bruteForce"
	"apt-web-server_v2/controllers/dataCopy"
	"apt-web-server_v2/controllers/dns"
	"apt-web-server_v2/controllers/equipment"
	"apt-web-server_v2/controllers/ids"
	"apt-web-server_v2/controllers/index"
	"apt-web-server_v2/controllers/netFlow"
	"apt-web-server_v2/controllers/offlineAssignment"
	"apt-web-server_v2/controllers/report"
	"apt-web-server_v2/controllers/scanEvent"
	"apt-web-server_v2/controllers/sniff"
	"apt-web-server_v2/controllers/ssl"
	"apt-web-server_v2/controllers/urgency"
	"apt-web-server_v2/controllers/vds"
	"apt-web-server_v2/controllers/waf"
	"apt-web-server_v2/controllers/whiteList"
	"apt-web-server_v2/modules/mconfig"
	"apt-web-server_v2/modules/mlog"
	"net/http"
	"regexp"

	es "github.com/castboy/es_ui_api/modules"

	"github.com/julienschmidt/httprouter"
)

var routerA *httprouter.Router

func init() {
	mlog.SetLogLevel(mlog.LevelDebug)

	es.InitLog()
	es.RegisterKeyword()

	Reg := regexp.MustCompile(`[\S]+`)
	nodes, _ := mconfig.Conf.RawString("es", "nodes")
	nodesSlice := Reg.FindAllString(nodes, -1)
	port, _ := mconfig.Conf.String("es", "port")

	es.Cli(nodesSlice, port)

	route()
}

func route() {
	routerA = httprouter.New()
	routerA.GET("/byzoro.apt.com/urgencies/list/day", urgency.UgcDObj.Get)
	routerA.GET("/byzoro.apt.com/urgency/trend/day", urgency.UrgencyObj.Get)
	routerA.GET("/byzoro.apt.com/urgency/classify/day", urgency.UgcClassifyObj.Get)
	routerA.GET("/byzoro.apt.com/urgencymold/count/day", urgency.UgcMoldCountObj.Get)
	routerA.GET("/byzoro.apt.com/urgencies/latest", urgency.UgcLObj.Get)
	routerA.GET("/byzoro.apt.com/attack/list/day", urgency.AttackCountObj.Get)
	routerA.GET("/byzoro.apt.com/dns/statistics/day", dns.DNSSObj.Get)
	routerA.GET("/byzoro.apt.com/sslcert/statistics/day", ssl.SSLCSObj.Get)
	routerA.GET("/byzoro.apt.com/sslcert/list/search", ssl.SSLCLstObj.Get)
	routerA.GET("/byzoro.apt.com/whitelist/search", whiteList.WL_SearchObj.Get)
	routerA.POST("/byzoro.apt.com/whitelist/operate", whiteList.WL_OperateObj.Post)
	routerA.GET("/byzoro.apt.com/flow/minute", netFlow.NetFlowObj.Get)
	routerA.GET("/byzoro.apt.com/malFlow/list/day", ids.MFDObj.Get)
	routerA.GET("/byzoro.apt.com/malFlow/classify/day", ids.MFCObj.Get)
	routerA.GET("/byzoro.apt.com/malFlow/trend/day", ids.MFTObj.Get)
	routerA.GET("/byzoro.apt.com/fileThreats/list/day", vds.FTDObj.Get)
	routerA.GET("/byzoro.apt.com/fileThreat/classify/day", vds.FTCObj.Get)
	routerA.GET("/byzoro.apt.com/fileThreat/trend/day", vds.FTTObj.Get)
	routerA.GET("/byzoro.apt.com/httpflow/list/day", waf.HFADObj.Get)
	routerA.GET("/byzoro.apt.com/httpflow/classify/day", waf.HFACObj.Get)
	routerA.GET("/byzoro.apt.com/httpflow/trend/day", waf.HFATObj.Get)
	routerA.GET("/byzoro.apt.com/bruteforce/list/day", bruteForce.BFDObj.Get)
	routerA.GET("/byzoro.apt.com/scanEvent/list/day", scanEvent.SEDObj.Get)
	routerA.GET("/byzoro.apt.com/off-line-dispatch/operate", offlineAssignment.OLPOObj.Get)
	routerA.GET("/byzoro.apt.com/off-line-dispatch/tasklist", offlineAssignment.TLObj.Get)
	routerA.GET("/byzoro.apt.com/off-line-dispatch/taskstatus", offlineAssignment.TSLObj.Get)
	routerA.GET("/byzoro.apt.com/sniffplcy/creat", sniff.SniffcreatObj.Get)
	routerA.GET("/byzoro.apt.com/sniffplcy/select", sniff.SniffselectObj.Get)
	routerA.GET("/byzoro.apt.com/sniffplcy/issued", sniff.SniffissuedObj.Get)
	routerA.GET("/byzoro.apt.com/sniffplcy/delete", sniff.SniffdeleteObj.Get)
	routerA.GET("/byzoro.apt.com/sniffplcy/show", sniff.SniffstatusObj.Get)
	routerA.POST("/byzoro.apt.com/off-line-dispatch/rule/operate", offlineAssignment.OLPOObj.Post)
	routerA.GET("/byzoro.apt.com/off-line-dispatch/rule/tasklist", offlineAssignment.TLObj.Rule)
	routerA.GET("/byzoro.apt.com/off-line-dispatch/rule/taskstatus", offlineAssignment.TSLObj.Rule)
	routerA.GET("/byzoro.apt.com/off-line-dispatch/rule/request", offlineAssignment.RULEObj.Get)
	//the rule operation for customer-defined-screening, such as add-rule, delete-rule , etc. in file offlineRuleOperate.go, class is OLROController
	routerA.POST("/byzoro.apt.com/off-line-dispatch/ruleoperate", offlineAssignment.OLRSDOObj.Post)
	routerA.GET("/byzoro.apt.com/off-line-dispatch/rulelist", offlineAssignment.OLRSDOObj.Get)
	/*start 暂时未用*/
	routerA.POST("/byzoro.apt.com/equipment/createdep", equipment.EquipmentcreatedepObj.Post)
	routerA.POST("/byzoro.apt.com/equipment/updatedep", equipment.EquipmentupdatedepObj.Post)
	/*end 暂时未用*/
	routerA.POST("/byzoro.apt.com/equipment/import", equipment.EquipmentimportObj.Post)
	routerA.GET("/byzoro.apt.com/equipment/manage", equipment.EquipmentmanageObj.Get)
	routerA.POST("/byzoro.apt.com/equipment/managedept", equipment.EquipmentmanagedeptObj.Post)
	routerA.GET("/byzoro.apt.com/equipment/department", equipment.EquipmentdepartmentObj.Get)
	routerA.GET("/byzoro.apt.com/datacopy/operate", dataCopy.DCTObj.Get)
	routerA.GET("/byzoro.apt.com/datacopy/list", dataCopy.DCTLObj.Get)
	routerA.GET("/byzoro.apt.com/datacopy/status", dataCopy.DCTSObj.Get)
	routerA.GET("/byzoro.apt.com/attack/topn", attackTopN.AttackObj.Get)
	routerA.GET("/byzoro.apt.com/index/attackstatistics", index.AttackStatisticsObj.Get)
	routerA.GET("/byzoro.apt.com/index/attackmonitor", index.AttackMonitorObj.Get)
	routerA.GET("/byzoro.apt.com/Report/Security", report.SecurityObj.Get)

	routerA.GET("/byzoro.apt.com/es-alert-search", es.Server)

	return
}

func main() {
	mconfig.Writepid()
	mlog.Debug("SSO start running...")
	port, _ := mconfig.Conf.String("server", "HttpPort")
	addr := ":" + port
	mlog.Debug("http listen on ", port, "...")
	http.ListenAndServe(addr, routerA)
	return
}
