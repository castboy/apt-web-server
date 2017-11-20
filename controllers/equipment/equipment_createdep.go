// equipment_createdep.go
package equipment

import (
	"apt-web-server_v2/controllers/public"
	"apt-web-server_v2/models/equipment"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

//"os"

func Getcreatebody(r *http.Request, input EquipmentcreatedepInput) (EquipmentcreatedepInput, error) {
	Body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return input, err
	}

	//fmt.Println(bytes.NewBuffer(Body).String())

	body_s := []byte(bytes.NewBuffer(Body).String())

	err = json.Unmarshal(body_s, &(input.Para))
	if nil != err {
		return input, err
		fmt.Println("your err:", err)
	}
	name, err := base64.StdEncoding.DecodeString(input.Para.Departmentname)
	if nil != err {
		return input, err
	}
	input.Para.Departmentname = string(name)
	//fmt.Println(input)
	return input, err
}

type EquipmentcreatedepController struct{}

type EquipmentcreatedepInput struct {
	Para equipment.Departmentip_st
}

var EquipmentcreatedepObj = new(EquipmentcreatedepController)

func (this *EquipmentcreatedepController) Post(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	//fmt.Println(r)
	//fmt.Println("string", r.Body)

	input := EquipmentcreatedepInput{}

	input, err := Getcreatebody(r, input)
	if err != nil {
		fmt.Println(err.Error())
		public.Write(w, public.ErrOkErr, "其他失败原因")
		return
	}
	//fmt.Println(input.Para)
	rst := equipment.Createdepartment(input.Para)
	public.Write(w, public.ErrOkErr, rst)
	return
}
