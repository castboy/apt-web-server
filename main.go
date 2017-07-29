package main

import (
	"apt-web-server/controllers"
	"apt-web-server/modules/mconfig"
	"apt-web-server/modules/mlog"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var routerA *httprouter.Router

func init() {
	mlog.SetLogLevel(mlog.LevelDebug)

	route()
}

func route() {
	routerA = httprouter.New()
	routerA.GET("/byzoro.apt.com/urgencies/list/day", controllers.UgcDObj.Get)
	routerA.GET("/byzoro.apt.com/urgency/trend/day", controllers.UrgencyObj.Get)
	routerA.GET("/byzoro.apt.com/urgency/classify/day", controllers.UgcClassifyObj.Get)
	routerA.GET("/byzoro.apt.com/urgencymold/count/day", controllers.UgcMoldCountObj.Get)
	routerA.GET("/byzoro.apt.com/attack/list/day", controllers.AttackCountObj.Get)
	routerA.GET("/byzoro.apt.com/dns/statistics/day", controllers.DNSSObj.Get)
	routerA.GET("/byzoro.apt.com/flow/minute", controllers.NetFlowObj.Get)
	routerA.GET("/byzoro.apt.com/malFlow/list/day", controllers.MFDObj.Get)
	routerA.GET("/byzoro.apt.com/malFlow/classify/day", controllers.MFCObj.Get)
	routerA.GET("/byzoro.apt.com/malFlow/trend/day", controllers.MFTObj.Get)
	routerA.GET("/byzoro.apt.com/fileThreats/list/day", controllers.FTDObj.Get)
	routerA.GET("/byzoro.apt.com/fileThreat/classify/day", controllers.FTCObj.Get)
	routerA.GET("/byzoro.apt.com/fileThreat/trend/day", controllers.FTTObj.Get)
	routerA.GET("/byzoro.apt.com/httpflow/list/day", controllers.HFADObj.Get)
	routerA.GET("/byzoro.apt.com/httpflow/classify/day", controllers.HFACObj.Get)
	routerA.GET("/byzoro.apt.com/httpflow/trend/day", controllers.HFATObj.Get)
	routerA.GET("/byzoro.apt.com/bruteforce/list/day", controllers.BFDObj.Get)
	routerA.GET("/byzoro.apt.com/scanEvent/list/day", controllers.SEDObj.Get)
	routerA.GET("/byzoro.apt.com/off-line-dispatch/operate", controllers.OLPOObj.Get)
	routerA.GET("/byzoro.apt.com/off-line-dispatch/tasklist", controllers.TLObj.Get)
	routerA.GET("/byzoro.apt.com/off-line-dispatch/taskstatus", controllers.TSLObj.Get)
	routerA.POST("/byzoro.apt.com/off-line-dispatch/rule/operate", controllers.OLPOObj.Post)
	routerA.GET("/byzoro.apt.com/off-line-dispatch/rule/tasklist", controllers.TLObj.Rule)
	routerA.GET("/byzoro.apt.com/off-line-dispatch/rule/taskstatus", controllers.TSLObj.Rule)
	routerA.GET("/byzoro.apt.com/off-line-dispatch/rule/request", controllers.RULEObj.Get)

	return
}

func main() {
	mlog.Debug("SSO start running...")
	port, _ := mconfig.Conf.String("server", "HttpPort")
	addr := ":" + port
	mlog.Debug("http listen on ", port, "...")
	http.ListenAndServe(addr, routerA)
	return
}
