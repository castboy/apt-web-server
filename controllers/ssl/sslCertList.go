/********SSL Cert 查询********/
package ssl

import (
	"apt-web-server_v2/controllers/public"
	"apt-web-server_v2/models/ssl"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	//"time"

	"github.com/julienschmidt/httprouter"
)

type SSLCLstController struct{}

var SSLCLstObj = new(SSLCLstController)

//func: Get for SSLCLstController, to fetch SSL detail info and response to PHP
//curl : keytype= [comnname|serialnum]
//       comnname=      , if keytype is 'comnname'
//       serialnum=    , if keytype is 'serialnum'
//       page=         , fetch ssl detail info of this page num
//       count=        , num of ssl info fetched per page
//       lastcount=    , the num of ssl info fetched regardless of page and count , its priority is higher than page and count
//
//w : Json format like   : {"code":10000,"msg":"success","data":{"total":4,"counts":4,"elements":[]}
//    Json elements like : {"s_cert_comnname":"www.161.com","s_cert_origname":"TenncetUnit","s_cert_unitname":"Tencent","s_cert_serialnum":"abcdef001","s_cert_notbefore":"2017-09-13 17:29:51","s_cert_notafter":"2017-09-25 07:16:31","s_cert_version":"3"}
func (this *SSLCLstController) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	input := SSLCLstGetInput{}

	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		public.Write(w, public.ErrParamentErr, "操作失败，内部消息解析失败。")
		fmt.Println("DBY SSL Cert Input parse err, Get intf received is : ", queryForm)
		return
	}

	//fmt.Println("DBY SSL Cert List Get intf received is : ", r.URL.RawQuery)
	//fmt.Println("DBY SSL Cert List Get intf received is : ", queryForm)
	//Write(w, ErrParamentErr, "This is for test, balabala.")
	//return

	searchtype := queryForm.Get("keytype")
	if searchtype != "serialnum" && searchtype != "comnname" {
		public.Write(w, public.ErrParamentErr, "操作失败，内部消息keytype值非法。")
		return
	} else {
		input.Para.Type = searchtype
	}

	if "comnname" == searchtype {

		comnName := queryForm.Get("comnname")
		if comnName == "" {
			public.Write(w, public.ErrParamentErr, "操作失败，内部消息commname值为空。")
			return
		} else {
			input.Para.ComnName = comnName
		}
	} else {

		serialnum := queryForm.Get("serialnum")
		if serialnum == "" {
			public.Write(w, public.ErrParamentErr, "操作失败，内部消息serialnum值为空。")
			return
		} else {
			input.Para.SerialNum = serialnum //serialnum is stored as string
		}
	}

	page, err := strconv.Atoi(queryForm.Get("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	input.Para.Page = int32(page)

	lastCount, err := strconv.Atoi(queryForm.Get("lastcount"))
	if err != nil {
		lastCount = 0
	} else if lastCount <= 0 {
		public.Write(w, public.ErrParamentErr, "操作失败，内部消息lastCount值非法。")
		return
	}
	input.Para.LastCount = int32(lastCount)

	count, err := strconv.Atoi(queryForm.Get("count"))
	if err != nil {
		count = 10
	} else if count <= 0 {
		public.Write(w, public.ErrParamentErr, "para count is not larger than 0.")
		return
	}
	input.Para.Count = int32(count)

	var list *ssl.TblSSLCLstData

	err, list = new(ssl.TblSSLCLst).GetSSLCLst(&input.Para)
	if err != nil {
		//這個異常沒有人 recover
		panic(fmt.Sprintf("cert_data GetSSLCLst error:%s", err.Error()))
		return
	}

	//response
	//fmt.Println("Get intf of SSL_Cert_List, response is : ", list)
	public.Write(w, public.ErrOkErr, list)
	return
}

type SSLCLstGetInput struct {
	Para ssl.TblSSLCLstSearchPara
}
