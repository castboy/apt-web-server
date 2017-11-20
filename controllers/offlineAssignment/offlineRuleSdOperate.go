/********对于自定义规则的规则进行操作，比如建立、删除等等********/
//package controllers
package offlineAssignment

import (
	"apt-web-server_v2/controllers/public"
	"apt-web-server_v2/models/offlineAssignment"
	//"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/julienschmidt/httprouter"
	//"os"
	//"io/ioutil"
	//"bytes"
	//"encoding/base64"
)

// offline rule operate controller
type OLROController struct{}

var OLRSDOObj = new(OLROController)

func (this *OLROController) Post(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var list *offlineAssignment.CMDResult
	//var err error

	fmt.Println("OLROObj.Post, r is : ", r)

	input, err := RuleOperParams(r)
	if err != nil {
		public.Write(w, public.ErrParamentErr, err.Error() /*"操作失败，内部Json格式错误。"*/)
		return
	}
	//ruleBytes, err := base64.StdEncoding.DecodeString(input.Para.Rule)
	//if nil != err {
	//	//TODO log
	//}

	switch input.Command {
	case "add":
		err = offlineAssignment.Rule_Defined_Add(&input)
		if err != nil {
			public.Write(w, public.ErrOperErr, err.Error())
			return
		}
	case "mod":
		err = offlineAssignment.Rule_Defined_Mod(&input)
		if err != nil {
			public.Write(w, public.ErrOperErr, err.Error())
			return
		}
	case "del":
		err = offlineAssignment.Rule_Defined_Del(&input)
		if err != nil {
			public.Write(w, public.ErrOperErr, err.Error())
			return
		}
	default:
		rspMsg := fmt.Sprintf(`自定义规则操作失败，操作命令 %s 非法。`, input.Command)
		err1 := errors.New(rspMsg)
		public.Write(w, public.ErrOperErr, err1.Error())
		return
	}

	public.Write(w, public.ErrOkErr, list)
}

func RuleOperParams(r *http.Request) (output offlineAssignment.TblRuleOperPara, err error) {
	var params offlineAssignment.TblRuleOperPara
	var err1 error
	var rspMsg string

	bytes := []byte(GetDataString(r))

	fmt.Println("RuleOperParams input r is ", r)

	err = json.Unmarshal(bytes, &params)
	if nil != err {
		fmt.Println("RuleParams Post Json format is error, and err is :", err)
		return params, err
	}

	//check post para
	if (params.Command != "add") && (params.Command != "del") && (params.Command != "mod") {
		rspMsg = fmt.Sprintf(`规则操作命令 %s 不合法。`, params.Command)
		err1 = errors.New(rspMsg)
		return params, err1
	}

	// 删除规则时，只关注 id，用 id 进行删除即可;
	if params.Command == "del" {
		return params, err
	}

	//AliasBytes, errAlias := base64.StdEncoding.DecodeString(params.Alias)
	//if nil != errAlias {
	//	rspMsg = fmt.Sprintf(`规则别名解码错误, %s 不合法。 Err is %s 。`, params.Alias, errAlias.Error())
	//	err1 = errors.New(rspMsg)
	//	return params, err1
	//}
	//params.Alias = string(AliasBytes)
	// 添加规则时，规则别名是必须携带的；并且此时没有 id ;
	if params.Command == "add" {
		if 0 == len(params.Alias) || len(params.Alias) > 64 {
			rspMsg = fmt.Sprintf(`规则别名长度不合法，别名长度为 %d 。`, len(params.Alias))
			err1 = errors.New(rspMsg)
			return params, err1
		}
	}

	varsetLen := len(params.VarSet)
	if varsetLen == 0 {
		//rspMsg = fmt.Sprintf(`规则携带的变量集个数不合法，最大值为5，实际个数为 %d 。`, varsetLen)
		rspMsg = fmt.Sprintf(`规则携带的变量集参数不能为空 。`)
		err1 = errors.New(rspMsg)
		return params, err1
	}
	//for _, ele := range params.VarSet {
	//	varLen := len(ele.Var)
	//	if varLen == 0 || varLen > 32 {
	//		rspMsg = fmt.Sprintf(`规则携带的变量集名长度不合法，最大长度为32，实际长度为 %d 。`, varLen)
	//		err1 = errors.New(rspMsg)
	//		return params, err1
	//	}
	//	VarInfoBytes, errVarInfo := base64.StdEncoding.DecodeString(ele.VarInfo)
	//	if nil != errAlias {
	//		rspMsg = fmt.Sprintf(`规则变量集的参数解码错误, %s 不合法。 Err is %s 。`, ele.VarInfo, errVarInfo.Error())
	//		err1 = errors.New(rspMsg)
	//		return params, err1
	//	}
	//	ele.VarInfo = string(VarInfoBytes)
	//	// 变量集名称的取值暂不检查
	//	varinfoLen := len(ele.VarInfo)
	//	if /*varinfoLen == 0 || */ varinfoLen > 32 {
	//		rspMsg = fmt.Sprintf(`规则携带的变量集名的参数长度不合法，最大长度为32，实际长度为 %d 。`, varinfoLen)
	//		err1 = errors.New(rspMsg)
	//		return params, err1
	//	}
	//}
	if len(params.VarSet) > 65535 {
		rspMsg = fmt.Sprintf(`规则携带的变量集长度不合法，最大长度为65535，实际长度为 %d 。`, params.VarSet)
		err1 = errors.New(rspMsg)
		return params, err1
	}

	// 操作符在UI界面是下拉框规定取值范围的，所以web service不做合法性判断了
	//if "@pm" != params.Oper && "@rx" != params.Oper {
	//	rspMsg = fmt.Sprintf(`规则携带的操作符类型不合法，实际为 %s 。`, params.Oper)
	//	err1 = errors.New(rspMsg)
	//	return params, err1
	//}
	operLen := len(params.Oper)
	if operLen == 0 || operLen > 23 {
		rspMsg := fmt.Sprintf(`规则携带的操作符的长度不合法，最大长度为23，实际长度为 %d 。`, operLen)
		err1 = errors.New(rspMsg)
		return params, err1
	}

	//OperInfoBytes, errOperInfo := base64.StdEncoding.DecodeString(params.OperInfo)
	//if nil != errOperInfo {
	//	rspMsg = fmt.Sprintf(`操作符的参数解码错误, %s 不合法。 Err is %s 。`, params.OperInfo, errOperInfo.Error())
	//	err1 = errors.New(rspMsg)
	//	return params, err1
	//}
	//params.OperInfo = string(OperInfoBytes)
	operinfoLen := len(params.OperInfo)
	if operinfoLen == 0 || operinfoLen > 512 {
		rspMsg := fmt.Sprintf(`规则携带的操作符的参数长度不合法，最大长度为512，实际长度为 %d 。`, operinfoLen)
		err1 = errors.New(rspMsg)
		return params, err1
	}

	// 事件函数在UI界面中是通过下拉框规定取值范围的，所以web service不做合法性判断了
	tfLen := len(params.TransFunc)
	if /*tfLen == 0 ||*/ tfLen > 23 {
		rspMsg := fmt.Sprintf(`规则携带的事件函数的长度不合法，最大长度为23，实际长度为 %d 。`, tfLen)
		err1 = errors.New(rspMsg)
		return params, err1
	}

	switch params.Phase {
	case "1":
	case "2":
	case "3":
	case "4":
	case "5":
		break
	default:
		rspMsg := fmt.Sprintf(`规则携带的处理阶段参数不合法，实际携带的是 %s 。`, params.Phase)
		err1 = errors.New(rspMsg)
		return params, err1
	}

	switch params.Severity {
	case "EMERGENCY":
	case "ALERT":
	case "CRITICAL":
	case "ERROR":
	case "WARNING":
	case "NOTICE":
	case "INFO":
	case "DEBUG":
		break
	default:
		rspMsg := fmt.Sprintf(`规则携带的危害等级参数不合法，实际携带的是 %s 。`, params.Severity)
		err1 = errors.New(rspMsg)
		return params, err1
	}

	switch params.Accuracy {
	case "0":
	case "1":
	case "2":
	case "3":
	case "4":
	case "5":
	case "6":
	case "7":
	case "8":
	case "9":
		break
	default:
		rspMsg := fmt.Sprintf(`规则携带的精确度参数不合法，实际携带的是 %s 。`, params.Accuracy)
		err1 = errors.New(rspMsg)
		return params, err1
	}

	switch params.Maturity {
	case "0":
	case "1":
	case "2":
	case "3":
	case "4":
	case "5":
	case "6":
	case "7":
	case "8":
	case "9":
		break
	default:
		rspMsg := fmt.Sprintf(`规则携带的成熟度参数不合法，实际携带的是 %s 。`, params.Maturity)
		err1 = errors.New(rspMsg)
		return params, err1
	}

	tagLen := len(params.Tag)
	if tagLen == 0 || tagLen > 32 {
		rspMsg := fmt.Sprintf(`规则携带的操作符的参数长度不合法，最大长度为32，实际长度为 %d 。`, tagLen)
		err1 = errors.New(rspMsg)
		return params, err1
	}

	// details 是规则的说明， PHP进行Base64编码，本地保存成Base64。查询时PHP得到后进行Base64解码后再显示在UI界面上
	msgLen := len(params.Details)
	if /*msgLen == 0 ||*/ msgLen > 512 {
		rspMsg := fmt.Sprintf(`规则携带的详细说明的长度不合法，最大长度为512字节，实际长度为 %d 。`, msgLen)
		err1 = errors.New(rspMsg)
		return params, err1
	}

	fmt.Println("RuleParams get para is :", params)

	return params, err
}

//func: Get func of OLROController, to fetch Rule of self-defined screening detail info, and response it to PHP
//curl : keytype= [comnname|serialnum]
//       comnname=      , if keytype is 'comnname'
//       serialnum=    , if keytype is 'serialnum'
//       page=         , fetch ssl detail info of this page num
//       count=        , num of ssl info fetched per page
//       lastcount=    , the num of ssl info fetched regardless of page and count , its priority is higher than page and count
//
//w : Json format like   : {"code":10000,"msg":"success","data":{"total":4,"counts":4,"elements":[]}
//    Json elements like : {"s_cert_comnname":"www.161.com","s_cert_origname":"TenncetUnit","s_cert_unitname":"Tencent","s_cert_serialnum":"abcdef001","s_cert_notbefore":"2017-09-13 17:29:51","s_cert_notafter":"2017-09-25 07:16:31","s_cert_version":"3"}
func (this *OLROController) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	input := OLROSearchGetInput{}
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		public.Write(w, public.ErrParamentErr, "规则查询失败，内部消息解码失败。")
		return
	}

	//fmt.Println("DBY Rule Get intf received 1 is : ", r.URL.RawQuery)
	//fmt.Println("DBY Rule Get intf received 2 is : ", queryForm)
	input.Para.Type = "all"

	page, err := strconv.Atoi(queryForm.Get("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	input.Para.Page = int32(page)

	count := queryForm.Get("count")
	if count != "" {
		countnum, err := strconv.Atoi(count)
		if err != nil || countnum < 0 {
			public.Write(w, public.ErrParamentErr, "规则查询失败，内部参数count值非法。")
			return
		}
		input.Para.Count = int32(countnum)
	} else {
		input.Para.Count = 1000
	}
	if input.Para.Count == 0 {
		input.Para.Count = 1000
	}

	fmt.Println("Rule search Get parameter input is : ", input.Para)

	//lastCount, err := strconv.Atoi(queryForm.Get("lastcount"))
	//if err != nil || lastCount <= 0 {
	//	public.Write(w, public.ErrParamentErr, "para lastCount is illegal.")
	//	return
	//}
	//input.Para.LastCount = int32(lastCount)

	var list *offlineAssignment.TblRuleSdLstData
	err, list = offlineAssignment.GetRuleSdLst(&input.Para)
	if err != nil || list == nil {
		//這個異常沒有人 recover
		panic(fmt.Sprintf("GetRuleSdLst error:%s", err.Error()))
		return
	}

	//encode by base64
	//list = EncodeRuleListAsNeeded(list)

	//response
	fmt.Println("Get Rule List, response is : ", list)

	public.Write(w, public.ErrOkErr, list)
	return
}

//func EncodeRuleListAsNeeded(inputList *offlineAssignment.TblRuleSdLstData) (outputList *offlineAssignment.TblRuleSdLstData) {
//
//	for idx, _ := range inputList.Elements {
//		if len(inputList.Elements[idx].Alias) > 0 {
//			inputList.Elements[idx].Alias = base64.StdEncoding.EncodeToString([]byte(inputList.Elements[idx].Alias))
//		}
//
//		for numIdx, _ := range inputList.Elements[idx].VarSet {
//			if len(inputList.Elements[idx].VarSet[numIdx].VarInfo) > 0 {
//				inputList.Elements[idx].VarSet[numIdx].VarInfo =
//					base64.StdEncoding.EncodeToString([]byte(inputList.Elements[idx].VarSet[numIdx].VarInfo))
//			}
//		}
//
//		if len(inputList.Elements[idx].Alias) > 0 {
//			inputList.Elements[idx].OperInfo = base64.StdEncoding.EncodeToString([]byte(inputList.Elements[idx].OperInfo))
//		}
//	}
//
//	fmt.Println("EncodeRuleListAsNeeded result is :", inputList)
//
//	return inputList
//}

type OLROGetInput struct {
	Para offlineAssignment.TblRuleOperPara
}

type OLROSearchGetInput struct {
	Para offlineAssignment.TblRuleSdSearchPara
}
