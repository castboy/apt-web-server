/********获取紧急事件详情********/
package offlineAssignment

import (
	"apt-web-server_v2/controllers/public"
	"apt-web-server_v2/models/offlineAssignment"
	"bytes"
	//"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (this *OLPOController) Post(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var list *offlineAssignment.CMDResult
	var err error

	input := Params(r)
	//ruleBytes, err := base64.StdEncoding.DecodeString(input.Para.Rule)
	//if nil != err {
	//	//TODO log
	//}

	switch input.Para.Cmmand {
	case "creat":
		if len(input.Para.RuleSet) == 0 {
			public.Write(w, public.ErrParamentErr, "任务建立失败，任务没有指定规则 。")
			return
		}
		ruleBytes, err1 := offlineAssignment.GetRulesConfText(&(input.Para))
		if nil != err1 {
			public.Write(w, public.ErrParamentErr, err1.Error())
			return
		}

		file := fmt.Sprintf("%s%s%s.conf", input.Para.Type,
			strconv.FormatInt(input.Para.Time, 10), input.Para.Name)
		WriteFile("/tmp/rules", file, ruleBytes)

		err, list = new(offlineAssignment.TblOLA).CreatAssignment(&input.Para)
	case "delete":
		err, list = new(offlineAssignment.TblOLA).DeleteAssignment(&input.Para)
	case "start":
		err, list = new(offlineAssignment.TblOLA).StartAssignment(&input.Para)

	case "shutdown":
		err, list = new(offlineAssignment.TblOLA).ShutDownAssignment(&input.Para)

	}

	if err != nil {
		panic(fmt.Sprintf("Offline self-defined screening, Last error:%s", err.Error()))
		return
	}

	public.Write(w, public.ErrOkErr, list)
}

func Params(r *http.Request) (input OLPOGetInput) {
	var params offlineAssignment.TblOLASearchPara

	bytes := []byte(GetDataString(r))

	json.Unmarshal(bytes, &params)

	//fmt.Println("Params Input r is : ", r)
	//fmt.Println("Params Json is : ", params)

	//input.Para = models.TblOLASearchPara{
	//	Cmmand:     params.Cmmand,
	//	Name:       params.Name,
	//	Time:       params.Time,
	//	Type:       "rule",
	//	Start:      params.Start,
	//	End:        params.End,
	//	Weight:     params.Weight,
	//	OfflineTag: "rule",
	//	Rule:       params.Rule,
	//	Details:    fmt.Sprintf("%s offline dispatch", input.Para.Type),
	//}
	input.Para.Cmmand = params.Cmmand
	input.Para.Name = params.Name
	input.Para.Time = params.Time
	input.Para.Type = "rule"
	input.Para.Start = params.Start
	input.Para.End = params.End
	input.Para.Weight = params.Weight
	input.Para.OfflineTag = "rule"
	input.Para.Rule = params.Rule
	input.Para.RuleSet = params.RuleSet
	//copy(input.Para.RuleSets, params.RuleSets)

	//fmt.Println("Params input.Para.RuleSets is ", input.Para.RuleSets)

	//var VarEle models.TblRuleSdSet
	//VarEle.Rule = ""
	//for len(input.Para.RuleSets) < 5 {
	//	input.Para.RuleSets = append(input.Para.RuleSets, VarEle)
	//}

	//for idx, _ := range params.RuleSets {
	////if len(para.RuleSets[idx].Rule) > 0 {
	////	singleRStr, err1 := GetSingleRuleConfText(&(para.RuleSets[idx]))
	////	if err1 != nil || singleRStr == "" {
	////		return RuleBytes, err1
	////	}
	////	qslice = append(qslice, singleRStr)
	////	qslice = append(qslice, "\n")
	////}
	////if idx >= 5 {
	////	break
	////}
	//	input.Para.RuleSets = append(input.Para.RuleSets, params.RuleSets[idx])
	//}

	//if len(input.Para.RuleSets) < 5 {
	//	input.Para.RuleSets = append(input.Para.RuleSets, VarEle)
	//}

	input.Para.Details = fmt.Sprintf("%s offline dispatch", input.Para.Type)

	fmt.Println("Params input result is ", input)

	return input
}

func GetDataString(req *http.Request) string {
	result, err := ioutil.ReadAll(req.Body)
	if err != nil {
	} else {

	}
	fmt.Println(bytes.NewBuffer(result).String())
	return bytes.NewBuffer(result).String()
}

func WriteFile(dir string, file string, bytes []byte) bool {
	isExist, err := pathExists(dir)
	if !isExist {
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			fmt.Printf("%s", err)
		} else {
			fmt.Print("Create Directory OK!")
		}
	}

	f, err := os.Create(dir + "/" + file)
	if nil != err {
		fmt.Println(err.Error())
	}

	defer f.Close()

	ok := true
	err = ioutil.WriteFile(dir+"/"+file, bytes, 0644)
	if nil != err {
		ok = false
	}

	return ok
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
