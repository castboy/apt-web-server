package public

import (
	//"apt-web-server/modules/mlog"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var ErrSsoMap map[int32]error
var ErrCodeMap map[error]int32

var ErrApOkErr error
var ErrOkErr error
var ErrPasswdErr error
var ErrCookieErr error
var ErrParamentErr error
var ErrLogoutErr error
var ErrNoPermErr error
var ErrPartFailErr error
var ErrNotExistErr error
var ErrNameRepeatErr error
var ErrErrorKeyErr error
var ErrHasExistErr error
var ErrCodeErr error
var ErrOperErr error

type RetCode struct {
	Code        int32
	Description string
}

type RespData struct {
	Code    int32       `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func init() {
	init_err_vars()
	init_local()
	init_sso()
}

func Write(w http.ResponseWriter, err error, data interface{}) {
	code, ok := ErrCodeMap[err]
	if !ok {
		panic("Write() parament error")
	}
	rsp := RespData{
		Code:    code,
		Message: err.Error(),
		Data:    data,
	}
	writeContent, _ := json.Marshal(rsp)
	fmt.Println(string(writeContent))
	//mlog.Debug(string(writeContent))
	w.Write(writeContent)
}

func init_err_vars() {
	ErrApOkErr = errors.New("success")
	ErrOkErr = errors.New("success")
	ErrPasswdErr = errors.New("password or username error")
	ErrCookieErr = errors.New("cookie error")
	ErrParamentErr = errors.New("parament error")
	ErrLogoutErr = errors.New("logout error")
	ErrNoPermErr = errors.New("no perm")
	ErrPartFailErr = errors.New("part error")
	ErrNotExistErr = errors.New("not exist")
	ErrNameRepeatErr = errors.New("name repeate")
	ErrErrorKeyErr = errors.New("key or secret error")
	ErrHasExistErr = errors.New("has exist")
	ErrCodeErr = errors.New("code error")
	ErrOperErr = errors.New("operate error")
}
func init_sso() {
	ErrSsoMap = make(map[int32]error)
	ErrSsoMap[10000] = ErrOkErr
	ErrSsoMap[20000] = ErrPasswdErr
	ErrSsoMap[20001] = ErrCookieErr
	ErrSsoMap[20002] = ErrParamentErr
	ErrSsoMap[20003] = ErrLogoutErr
	ErrSsoMap[20004] = ErrNoPermErr
	ErrSsoMap[20007] = ErrNameRepeatErr
	ErrSsoMap[20009] = ErrErrorKeyErr
	ErrSsoMap[20010] = ErrHasExistErr
	ErrSsoMap[20011] = ErrCodeErr
	ErrSsoMap[20012] = ErrOperErr
}
func init_local() {
	ErrCodeMap = make(map[error]int32)
	ErrCodeMap[ErrApOkErr] = 0
	ErrCodeMap[ErrOkErr] = 10000
	ErrCodeMap[ErrPasswdErr] = 20000
	ErrCodeMap[ErrCookieErr] = 20001
	ErrCodeMap[ErrParamentErr] = 20002
	ErrCodeMap[ErrLogoutErr] = 20003
	ErrCodeMap[ErrNoPermErr] = 20004
	ErrCodeMap[ErrPartFailErr] = 20005
	ErrCodeMap[ErrNotExistErr] = 20006
	ErrCodeMap[ErrNameRepeatErr] = 20007
	ErrCodeMap[ErrErrorKeyErr] = 20009
	ErrCodeMap[ErrHasExistErr] = 20010
	ErrCodeMap[ErrCodeErr] = 20011
	ErrCodeMap[ErrOperErr] = 20012
}
