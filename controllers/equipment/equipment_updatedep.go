// equipment_updatedep.go
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

func Getupdatebody(r *http.Request, input EquipmentupdatedepInput) (EquipmentupdatedepInput, error) {
	Body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return input, err
	}

	//fmt.Println(bytes.NewBuffer(Body).String())

	body_s := []byte(bytes.NewBuffer(Body).String())

	err = json.Unmarshal(body_s, &(input.Para))
	if nil != err {
		fmt.Println("your err:", err)
		return input, err
	}
	name, err := base64.StdEncoding.DecodeString(input.Para.Departmentname)
	if err != nil {
		return input, err
	}
	input.Para.Departmentname = string(name)
	//fmt.Println(input)
	return input, err
}

type EquipmentupdatedepController struct{}

type EquipmentupdatedepInput struct {
	Para equipment.Departmentip_st
}

var EquipmentupdatedepObj = new(EquipmentupdatedepController)

func (this *EquipmentupdatedepController) Post(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	//fmt.Println(r)
	//fmt.Println("string", r.Body)

	input := EquipmentupdatedepInput{}

	input, err := Getupdatebody(r, input)
	if err != nil {
		fmt.Println(err.Error())
		public.Write(w, public.ErrOkErr, "其他失败原因")
		return
	}
	//fmt.Println(input.Para)
	rst := equipment.Updatedepartment(input.Para)
	public.Write(w, public.ErrOkErr, rst)
}
