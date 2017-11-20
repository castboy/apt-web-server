/********White List 查询********/
package whiteList

import (
	"apt-web-server_v2/controllers/public"
	"apt-web-server_v2/models/whiteList"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	//"time"

	"github.com/julienschmidt/httprouter"
)

type WhiteListSearchController struct{}

var WL_SearchObj = new(WhiteListSearchController)

//func: Get func of WhiteListSearchController, to fetch White list detail info and response to PHP
//curl : keytype= [comnname|serialnum]
//       comnname=      , if keytype is 'comnname'
//       serialnum=    , if keytype is 'serialnum'
//       page=         , fetch ssl detail info of this page num
//       count=        , num of ssl info fetched per page
//       lastcount=    , the num of ssl info fetched regardless of page and count , its priority is higher than page and count
//
//w : Json format like   : {"code":10000,"msg":"success","data":{"total":4,"counts":4,"elements":[]}
//    Json elements like : {"s_cert_comnname":"www.161.com","s_cert_origname":"TenncetUnit","s_cert_unitname":"Tencent","s_cert_serialnum":"abcdef001","s_cert_notbefore":"2017-09-13 17:29:51","s_cert_notafter":"2017-09-25 07:16:31","s_cert_version":"3"}
func (this *WhiteListSearchController) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	input := WLSearchGetInput{}

	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		public.Write(w, public.ErrParamentErr, "白名单查询失败，内部消息解码失败。")
		return
	}

	//fmt.Println("DBY White List Get intf received 1 is : ", r.URL.RawQuery)
	//fmt.Println("DBY White List Get intf received 2 is : ", queryForm)

	//Write(w, ErrParamentErr, "This is for test, balabala.")
	//return

	//searchtype := queryForm.Get("keytype")
	//if searchtype != "all" && searchtype != "exact" {
	//	Write(w, ErrParamentErr, "para keytype is illegal.")
	//	return
	//} else {
	//	input.Para.Type = searchtype
	//}

	//if "exact" == searchtype {

	//if anyone of 5 tuples is ommited, handle it as *  during searching White list.
	sip := queryForm.Get("sip")
	if sip != "" {
		input.Para.Sip = sip
	}

	s_port := queryForm.Get("sport")
	if s_port != "" {
		sport, err := strconv.Atoi(s_port)
		if err != nil || sport < 0 {
			public.Write(w, public.ErrParamentErr, "白名单查询失败，源端口号是非法值。")
			return
		}
		input.Para.Sport = int32(sport)
	} else {
		input.Para.Sport = 0
	}

	dip := queryForm.Get("dip")
	if dip != "" {
		input.Para.Dip = dip
	}

	d_port := queryForm.Get("dport")
	if d_port != "" {
		dport, err := strconv.Atoi(d_port)
		if err != nil || dport < 0 {
			public.Write(w, public.ErrParamentErr, "白名单查询失败，目的端口号是非法值。")
			return
		}
		input.Para.Dport = int32(dport)
	} else {
		input.Para.Dport = 0
	}

	protocol := queryForm.Get("proto")
	if protocol != "" {
		//protonum, err := strconv.Atoi(protocol)
		//if err != nil || protonum < 0 {
		//	Write(w, ErrParamentErr, "白名单查询失败，协议号是非法值。")
		//	return
		//}
		//input.Para.Proto = int32(protonum)
		if false == WLCheckProto(protocol) {
			public.Write(w, public.ErrParamentErr, "白名单查询失败，协议是非法值。")
			return
		}
		input.Para.Proto = int32(whiteList.MapProto2Num[protocol])
	} else {
		input.Para.Proto = 0
	}

	if ("" == input.Para.Sip) && ("" == input.Para.Dip) && (0 == input.Para.Sport) && (0 == input.Para.Dport) && (0 == input.Para.Proto) {
		input.Para.Type = "all"
	} else {
		input.Para.Type = "exact"
	}

	page, err := strconv.Atoi(queryForm.Get("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	input.Para.Page = int32(page)

	//count, err := strconv.Atoi(queryForm.Get("count"))
	//if err != nil || count <= 0 {
	//	Write(w, ErrParamentErr, "para count is illegal.")
	//	return
	//}
	count := queryForm.Get("count")
	if count != "" {
		countnum, err := strconv.Atoi(count)
		if err != nil || countnum <= 0 {
			public.Write(w, public.ErrParamentErr, "白名单查询失败，内部参数count值非法。")
			return
		}
		input.Para.Count = int32(countnum)
	} else {
		input.Para.Count = 2000
	}

	//fmt.Println("WL search Get parameter input is : ", input.Para)

	//lastCount, err := strconv.Atoi(queryForm.Get("lastcount"))
	//if err != nil || lastCount <= 0 {
	//	Write(w, ErrParamentErr, "para lastCount is illegal.")
	//	return
	//}
	//input.Para.LastCount = int32(lastCount)

	var list *whiteList.TblWLLstData
	err, list = new(whiteList.TblWLLst).GetWLLCLst(&input.Para)
	if err != nil || list == nil {
		//這個異常沒有人 recover
		panic(fmt.Sprintf("cert_data GetWLLCLst error:%s", err.Error()))
		return
	}

	//response
	//fmt.Println("Get intf of White List, response is : ", list)
	public.Write(w, public.ErrOkErr, list)
	return
}

func WLCheckProto(pro string) bool {
	//空字符代表协议是通配，是合法值
	if pro == "" {
		return true
	}
	if whiteList.MapProto2Num[pro] == 0 {
		return false
	} else {
		return true
	}
}

type WLSearchGetInput struct {
	Para whiteList.TblWLSearchPara
}
