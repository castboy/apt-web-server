/******** White List Operation /對白名單內容進行各種操作，比如增加、刪除等 ********/
package whiteList

import (
	//"github.com/julienschmidt/httprouter"
	"apt-web-server_v2/controllers/public"
	"apt-web-server_v2/models/whiteList"
	"bytes"
	//"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	//"os"
	//"strconv"
	//_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"

	"errors"
)

type WhiteListOperateController struct{}

var WL_OperateObj = new(WhiteListOperateController)

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
func (this *WhiteListOperateController) Post(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var list *whiteList.CMDResult
	var err error

	input, err := WL_Params(r)
	if err != nil {
		public.Write(w, public.ErrParamentErr, err.Error() /*"操作失败，内部Json格式错误。"*/)
		return
	}
	//fmt.Println("WL Post msg para is ", input)
	//ruleBytes, err := base64.StdEncoding.DecodeString(input.Para.Rule)
	//if nil != err {
	//	//TODO log
	//}
	//return

	switch input.Cmmand {
	case "add":
		var numRepeated int32
		for _, ele := range input.WLOpElement {
			//if (ele.Sip == "") && (ele.Dip == "") && (ele.Sport == "") && (ele.Dport == "") && (ele.Proto == "") {
			if (ele.Sip == "") && (ele.Dip == "") && (ele.Sport == 0) && (ele.Dport == 0) && (ele.Proto == 0) {
				//Write(w, ErrOperErr, "WL Post cmd add error, 5 tuples can't be all blank.")
				//fmt.Println("WL Post cmd add error, 5 tuples can't be all blank.")
				public.Write(w, public.ErrOperErr, "白名单添加失败，白名单信息不能全部为空。")
				return
			}
		}

		err, numRepeated = whiteList.WL_add(&input)
		if err != nil {
			public.Write(w, public.ErrOperErr, "白名单添加失败，发生内部错误。")
			return
		}
		sLen := len(input.WLOpElement)
		if numRepeated == 1 && sLen == 1 {
			public.Write(w, public.ErrOperErr, "白名单添加失败，该白名单已经存在。")
			return
		} else if numRepeated >= 1 && numRepeated == int32(sLen) {
			rspMsg := fmt.Sprintf(`白名单添加失败，这些白名单已经存在。`)
			public.Write(w, public.ErrOperErr, rspMsg)
			return
		} else if numRepeated >= 1 {
			rspMsg := fmt.Sprintf(`白名单添加成功, 添加成功条数：%d , 重复未添加条数：%d 。`,
				(int32(sLen) - numRepeated), numRepeated)
			public.Write(w, public.ErrOkErr, rspMsg)
			return
		}
	case "delete":
		var notFound int32
		err, notFound = whiteList.WL_delete(&input)
		if err != nil {
			public.Write(w, public.ErrOperErr, "白名单删除失败，发生内部错误。")
			return
		}

		sLen := len(input.WLOpElement)

		if notFound == 1 && sLen == 1 {
			public.Write(w, public.ErrOperErr, "白名单删除失败，该白名单不存在。")
			return
		} else if notFound >= 1 && notFound == int32(sLen) {
			rspMsg := fmt.Sprintf(`白名单删除失败，这些白名单不存在。`)
			public.Write(w, public.ErrOperErr, rspMsg)
			return
		} else if notFound >= 1 {
			rspMsg := fmt.Sprintf(`白名单删除成功，删除成功条数：%d , 不存在的条数：%d 。`,
				(int32(sLen) - notFound), notFound)
			public.Write(w, public.ErrOperErr, rspMsg)
			return
		}
	case "clear":
		err = whiteList.WL_clear(&input)
		if err != nil {
			public.Write(w, public.ErrOperErr, "清空白名单操作失败，发生内部错误。")
			return
		}

	default:
		public.Write(w, public.ErrParamentErr, "内部错误，操作类型异常。")
		return
	}

	if err != nil {
		panic(fmt.Sprintf("WL Post intf GetLast error:%s", err.Error()))
		public.Write(w, public.ErrParamentErr, "白名单操作失败，发生了内部错误。")
		return
	}

	public.Write(w, public.ErrOkErr, list)
}

func WL_Params(r *http.Request) (output whiteList.TblWLOperateParaIn, err error) {
	var params whiteList.TblWLOperatePara
	var paramsIn whiteList.TblWLOperateParaIn
	var eleTmp whiteList.TblWLOperElementIn

	bytes := []byte(WLGetDataString(r))

	err = json.Unmarshal(bytes, &params)
	if nil != err {
		fmt.Println("Post Json format is error, and err is :", err)
		return output, err
	}

	//fmt.Println("WL_Params get para string is ", params)
	//return outputs
	//fmt.Println("WL_Params output params struct is ", params)

	paramsIn.Cmmand = params.Cmmand
	paramsIn.OpNum = params.OpNum
	for idx, _ := range params.WLOpElement {
		eleTmp.Sip = params.WLOpElement[idx].Sip
		eleTmp.Sport = params.WLOpElement[idx].Sport
		eleTmp.Dip = params.WLOpElement[idx].Dip
		eleTmp.Dport = params.WLOpElement[idx].Dport
		if false == WLCheckProto(params.WLOpElement[idx].Proto) {
			var err1 error
			err1 = errors.New("协议类型错误")
			return paramsIn, err1
		}
		if params.WLOpElement[idx].Proto == "" {
			eleTmp.Proto = 0
		} else {
			eleTmp.Proto = int32(whiteList.MapProto2Num[params.WLOpElement[idx].Proto])
		}

		paramsIn.WLOpElement = append(paramsIn.WLOpElement, eleTmp)
	}

	return paramsIn, err
}

func WLGetDataString(req *http.Request) string {
	result, err := ioutil.ReadAll(req.Body)
	if err != nil {
	} else {

	}

	strRlt := bytes.NewBuffer(result).String()
	//fmt.Println("WLGetDataString White List Operate Post Para string is ", strRlt)

	return strRlt
}

type WLOperPostInput struct {
	Para whiteList.TblWLOperatePara
}
