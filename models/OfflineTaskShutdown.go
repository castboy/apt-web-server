package models

import (
	"apt-web-server/modules/mlog"
	"encoding/json"
	"fmt"
	//	"time"
)

func (this *TblOLA) ShutDownAssignment(para *TblOLASearchPara) (error, *CMDResult) {
	var res_t CMDResult
	var taskId int
	var agentPar AgentPara
	query := fmt.Sprintf(`select id,start,end,type,weight,topic,status from %s where name='%s' and time=%d;`,
		this.TableName(para.OfflineTag),
		para.Name,
		para.Time)
	rows, err := db.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	for rows.Next() {
		ugc := new(TblOLA)
		err = rows.Scan(
			&ugc.Id,
			&ugc.Start,
			&ugc.End,
			&ugc.Type,
			&ugc.Weight,
			&ugc.Topic,
			&ugc.Status)
		if err != nil {
			return err, nil
		}
		taskId = ugc.Id
		agentPar.Engine = ugc.Type
		agentPar.Weight = ugc.Weight
		agentPar.Topic = ugc.Topic
	}
	/******获取picker状态并向picker发送shutdown消息******/
	pickerKey := fmt.Sprintf("%s/%d", PickerETCDStatusKey, taskId)
	shutdownPickerCmd := fmt.Sprintf("%s %s", ShutdownPicker, pickerKey)
	pickerCount, pickerTotal, err := GetEtcdPicker(pickerKey, PickerETCDIpPort)
	if err != nil {
		mlog.Debug("OfflineTaskShutdown's GetEtcdPicker error")
	}
	if pickerTotal != 0 && pickerCount != pickerTotal {
		err = SSHCmd(PickerSSHUser, PickerSSHPass, PickerSSHIP, shutdownPickerCmd, SSHPort)
		if err != nil {
			mlog.Debug("picker shutdown fail!")
		}
	}

	//向agent发送"shutdown"消息
	agentPar.SignalType = "shutdown"
	agentCmdShutdown, err := json.Marshal(agentPar)
	if nil != err {
		mlog.Debug("`shutdown` json.Marshal err")
	}
	err = SendOfflineMsg(agentCmdShutdown)
	if nil != err {
		mlog.Debug("send `shutdown` msg failed")
	}

	/******删除etcd和topic ******/
	_, err = EtcdCmd("delete", pickerKey, "", PickerETCDIpPort)
	if err != nil {
		fmt.Println("delete picker etcd fail!")
	}
	/*	_, err = EtcdCmd("delete", agentCmdKey, "", AgentETCDCmdIpPort)
		if err != nil {
			fmt.Println("delete agent etcd fail!")
		}
	*/
	deleteTopicCmd := fmt.Sprintf(`kafka-topics --zookeeper %s --topic %s --delete`, KafkaTopicIpPort, agentPar.Topic)
	err = SSHCmd(TopicSSHUser, TopicSSHPass, TopicSSHIP, deleteTopicCmd, SSHPort)
	if err != nil {
		fmt.Println("delete topic error!", err)
	}
	/***************************/
	/******设置任务状态******/
	err = this.UpgradeStatus("status", "shutdown", taskId, para.OfflineTag)
	if err != nil {
		fmt.Println("set task ", taskId, " error!")
	}
	/*********************/
	res_t.Result = "ok"
	return nil, &res_t
}
