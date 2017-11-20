// Equipment_import.go
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

func Getbody(r *http.Request, input EquipmentimportInput) (EquipmentimportInput, error) {
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
	for i, _ := range input.Para.Equipment {
		alias, err := base64.StdEncoding.DecodeString(input.Para.Equipment[i].Alias)
		if nil != err {
			return input, err
		}
		input.Para.Equipment[i].Alias = string(alias)
	}
	fmt.Println(input)

	return input, err
}

type EquipmentimportController struct{}

type EquipmentimportInput struct {
	Para equipment.Importdate_st
}

var EquipmentimportObj = new(EquipmentimportController)

func (this *EquipmentimportController) Post(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	//fmt.Println(r)
	//fmt.Println("string", r.Body)

	input := EquipmentimportInput{}

	input, err := Getbody(r, input)
	if err != nil {
		public.Write(w, public.ErrOkErr, "err")
		return
	}
	//fmt.Println(input.Para)
	rst := equipment.Import_mysql(input.Para)
	if rst == false {
		public.Write(w, public.ErrOkErr, "err")
		return
	}

	public.Write(w, public.ErrOkErr, "ok")
}
