/********获取紧急事件详情********/
package controllers

import (
	"apt-web-server/models"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (this *OLPOController) Post(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var list *models.CMDResult
	var err error

	input := Params(r)
	ruleBytes, err := base64.StdEncoding.DecodeString(input.Para.Rule)
	if nil != err {
		//TODO log
	}

	switch input.Para.Cmmand {
	case "creat":
		file := fmt.Sprintf("%s%s%s.conf", input.Para.Type,
			strconv.FormatInt(input.Para.Time, 10), input.Para.Name)
		WriteFile("/tmp/rules", file, ruleBytes)

		err, list = new(models.TblOLA).CreatAssignment(&input.Para)
	case "delete":
		err, list = new(models.TblOLA).DeleteAssignment(&input.Para)
	case "start":
		err, list = new(models.TblOLA).StartAssignment(&input.Para)

	case "shutdown":
		err, list = new(models.TblOLA).ShutDownAssignment(&input.Para)

	}

	if err != nil {
		panic(fmt.Sprintf("TblUrgency GetLast error:%s", err.Error()))
		return
	}

	Write(w, ErrOkErr, list)
}

func Params(r *http.Request) (input OLPOGetInput) {
	var params models.TblOLASearchPara

	bytes := []byte(GetDataString(r))

	json.Unmarshal(bytes, &params)

	input.Para = models.TblOLASearchPara{
		Cmmand:     params.Cmmand,
		Name:       params.Name,
		Time:       params.Time,
		Type:       "rule",
		Start:      params.Start,
		End:        params.End,
		Weight:     params.Weight,
		OfflineTag: "rule",
		Rule:       params.Rule,
		Details:    fmt.Sprintf("%s offline dispatch", input.Para.Type),
	}
	fmt.Println(input)

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
