package offlineAssignment

import (
	//"apt-web-server_v2/models/db"
	"apt-web-server_v2/modules/mlog"
	"encoding/json"
	"fmt"
	//	"time"
)

func (this *TblOLA) ShutDownAssignment(para *TblOLASearchPara) (error, *CMDResult) {
	var res_t CMDResult
	err := this.GetTaskMsg(para)
	if err != nil {
		res_t.Result = "faild"
		mlog.Debug("Get task message error")
		return err, &res_t
	}
	var agentPar AgentPara
	agentPar.Engine = this.Type
	agentPar.Weight = this.Weight
	agentPar.Topic = this.Topic
	/******获取picker状态并向picker发送shutdown消息******/
	pickerKey := fmt.Sprintf("%s/%d", PickerStatusKey, this.Id)
	stopPickerCmd := fmt.Sprintf("%s %s", ShutdownPicker, pickerKey)
	pickerStat, err := GetEtcdPicker(pickerKey)
	if err != nil {
		mlog.Debug("OfflineTaskShutdown's GetEtcdPicker error")
	}
	if pickerStat.Total != 0 && pickerStat.Count != pickerStat.Total {
		err = SSHCmd(PickerSSHUser, PickerSSHPass, PickerSSHIP, stopPickerCmd, SSHPort)
		if err != nil {
			mlog.Debug("picker shutdown fail!")
		}
	}
	/******获取agent状态并向agent发送shutdown消息******/
	agentCmdKey := fmt.Sprintf(`%s/%d`, AgentETCDCmdKey, this.Id)
	agentPar.SignalType = "shutdown"
	agentCmdShutdown, err := json.Marshal(agentPar)
	if nil != err {
		mlog.Debug("`shutdown` json.Marshal err")
	}
	err = SendOfflineMsg(agentCmdShutdown)
	if nil != err {
		mlog.Debug("send `shutdown` msg failed")
	}
	_, err = EtcdCmd("put", agentCmdKey, string(agentCmdShutdown))
	if nil != err {
		mlog.Debug("put shutdown to agent etcd error:", err)
	}
	/******删除etcd和topic ******/
	_, err = EtcdCmd("delete", pickerKey, "")
	if err != nil {
		mlog.Debug("delete picker etcd fail!", err)
	}
	_, err = EtcdCmd("delete", agentCmdKey, "")
	if err != nil {
		fmt.Println("delete agent etcd fail!")
	}
	DelTopic(this.Topic)
	/******设置任务状态******/
	err = this.UpgradeStatus("status", "shutdown", this.Id, para.OfflineTag)
	if err != nil {
		mlog.Debug("set task ", this.Id, " error!")
	}
	/*********************/
	res_t.Result = "ok"
	return nil, &res_t
}
