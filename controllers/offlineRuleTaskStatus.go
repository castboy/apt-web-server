/*离线任务状态列表*/
package controllers

import (
	"apt-web-server/models"
	"fmt"
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
)

func (this *TSLController) Rule(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//check ssn
	//TODO
	//parse parament
	input := TSLGetInput{}

	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		Write(w, ErrParamentErr, nil)
		return
	}

	tasklist := queryForm.Get("tasklist")
	if tasklist == "" {
		fmt.Println("the task list is null")
	} else {
		input.Para.TaskList = tasklist
		fmt.Println("tasklist:", tasklist)
	}

	input.Para.OfflineTag = "rule"

	err, list := new(models.TblOLA).GetTaskStatus(&input.Para)
	if err != nil {
		panic(fmt.Sprintf("TblFTD GetLast error:%s", err.Error()))
		return
	}
	//response
	Write(w, ErrOkErr, list)
}
