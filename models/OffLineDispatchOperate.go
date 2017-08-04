/********离线调度********/
package models

import (
	"apt-web-server/modules/mlog"
	"encoding/json"
	"fmt"
	"time"
)

/************test config************
var (
	AgentStatusETCDSlice = []string{
		"apt/agent/status/192.168.1.203"}
	//"apt/agent/status/192.168.1.105"}
	AgentStatusETCDIpPort = []string{
		"http://192.168.1.203:2379"}
	//"http://10.88.1.105:2379"}
)

const (
	//Abs:"Absolute"
	PyAbsPath           = "PYSPARK_PYTHON=/usr/bin/python2.7"
	SparkAbsPath        = "/home/spark-2.1.1-bin-hadoop2.6/bin/spark-submit"
	SparkParaList       = "--master yarn --deploy-mode cluster --name"
	PickerScrAbsPath    = "/home/apt/picker.py"
	AgentETCDCmdKey     = "apt/agent/offlineReq"
	AgentETCDCmdIpPort  = "http://192.168.1.204:2379"
	KafkaTopicIpPort    = "192.168.1.203:2181"
	PickerETCDStatusKey = "picker"
	PickerETCDCmdKey    = "pickerCmd"
	ShutdownPicker      = "/bin/bash /home/apt/stopYarnApp.sh"
	PickerETCDIpPort    = "http://192.168.1.204:2379"
	TopicSSHUser        = "root"
	TopicSSHPass        = "aaaaaa"
	TopicSSHIP          = "192.168.1.203"
	PickerSSHUser       = "root"
	PickerSSHPass       = "aaaaaa"
	PickerSSHIP         = "192.168.1.204"
	SSHPort             = 22
)
*/
/************default config************/
var (
	AgentStatusETCDSlice = []string{
		"apt/agent/status/192.168.1.103"}
	//"apt/agent/status/192.168.1.105"}
	AgentStatusETCDIpPort = []string{
		"http://192.168.1.103:2379"}
	//"http://10.88.1.105:2379"}
)

const (
	IP                  = "192.168.1.103"
	PyAbsPath           = "PYSPARK_PYTHON=/usr/local/bin/python2.7"
	SparkAbsPath        = "/home/spark-2.1.1-bin-hadoop2.6/bin/spark-submit"
	SparkParaList       = "--master yarn --deploy-mode cluster --name"
	PickerScrAbsPath    = "/home/apt/picker.py"
	AgentETCDCmdKey     = "apt/agent/offlineReq"
	AgentETCDCmdIpPort  = "http://" + IP + ":2379"
	KafkaTopicIpPort    = IP + ":2181"
	PickerETCDStatusKey = "picker"
	ShutdownPicker      = "/bin/bash /home/apt/stopYarnApp.sh"
	PickerETCDIpPort    = "http://" + IP + ":2379"
	TopicSSHUser        = "root"
	TopicSSHPass        = "aaaaaa"
	TopicSSHIP          = IP
	PickerSSHUser       = "root"
	PickerSSHPass       = "aaaaaa"
	PickerSSHIP         = IP
	SSHPort             = 22

	TOPIC     = "offline_msg"
	PARTITION = 0
	KAFKA     = "192.168.1.103"
)

func (this *TblOLA) TableName(offlineTag string) string {
	var tbl string
	if "rule" == offlineTag {
		tbl = "offline_assignment_rule"
	} else {
		tbl = "offline_assignment"
	}
	return tbl
}

func (this *TblOLA) CreatAssignment(para *TblOLASearchPara) (error, *CMDResult) {
	var res_t CMDResult
	if para.Weight == 0 {
		para.Weight = 5
	}
	paraTopic := fmt.Sprintf(`%s%d%s`, para.Type, para.Time, para.Name)
	query := fmt.Sprintf(`insert into %s(name,time,type,start,end,weight,topic,status,details) 
						value('%s',%d,'%s','%s','%s',%d,'%s','%s','%s');`,
		this.TableName(para.OfflineTag),
		para.Name,
		para.Time,
		para.Type,
		para.Start,
		para.End,
		para.Weight,
		paraTopic,
		"ready",
		para.Details)
	rows, err := db.Query(query)
	fmt.Println(query)
	if err != nil {
		mlog.Debug(query, "CreatAssignment error")
		res_t.Result = "faild"
		return err, &res_t
	}
	defer rows.Close()

	res_t.Result = "ok"
	return nil, &res_t
}

func (this *TblOLA) DeleteAssignment(para *TblOLASearchPara) (error, *CMDResult) {
	var res_t CMDResult
	query := fmt.Sprintf(`delete from %s where name='%s' and time=%d;`, this.TableName(para.OfflineTag), para.Name, para.Time)
	rows, err := db.Query(query)
	fmt.Println(query)
	if err != nil {
		mlog.Debug(query, "DeleteAssignment error")
		res_t.Result = "faild"
		return err, &res_t
	}
	defer rows.Close()

	res_t.Result = "ok"
	return nil, &res_t
}

func (this *TblOLA) StartAssignment(para *TblOLASearchPara) (error, *CMDResult) {
	fmt.Println("开始离线任务")
	var res_t CMDResult
	var agentPar AgentPara
	query := fmt.Sprintf(`select id,start,end,type,weight,topic,status from %s where name='%s' and time=%d;`,
		this.TableName(para.OfflineTag),
		para.Name,
		para.Time)
	rows, err := db.Query(query)
	fmt.Println("sql 语句:", query)
	if err != nil {
		fmt.Println("执行sql语句失败")
		res_t.Result = "faild"
		mlog.Debug(query, "StartAssignment get task status error")
		return err, &res_t
	}
	defer rows.Close()
	var pickerETCDKey, startPickerCmd, fileType, topicName, taskType string
	var taskID int

	for rows.Next() {
		fmt.Println("进入next")
		ugc := new(TblOLA)
		err = rows.Scan(
			&ugc.Id,
			&ugc.Start,
			&ugc.End,
			&ugc.Type,
			&ugc.Weight,
			&ugc.Topic,
			&ugc.Status)

		fmt.Println("ugc.Id:", ugc.Id, "ugc:", ugc)
		if err != nil {
			return err, nil
		}
		taskID = ugc.Id
		taskType = ugc.Type
		topicName = ugc.Topic
		agentPar.Engine = taskType
		agentPar.Topic = topicName
		agentPar.Weight = ugc.Weight

		/******判断任务运行状态，避免重复运行******/

		fmt.Println("判断运行状态")
		if ugc.Status == "running" {
			res_t.Result = "this task is running!"
			return nil, &res_t
		} else if ugc.Status == "ready" || ugc.Status == "error" {
			err := this.UpgradeStatus("status", "running", taskID, para.OfflineTag)
			if err != nil {
				mlog.Debug("Update status task ", taskID, "to running error")
				res_t.Result = "Update task to running error"
				return nil, &res_t
			}
		}
		/************************************/

		switch ugc.Type {
		case "vds":
			fileType = "file"
		case "waf":
			fileType = "http"
		case "rule":
			fileType = "http"
		}
		pickerETCDKey = fmt.Sprintf("%s/%d", PickerETCDStatusKey, ugc.Id)
		startPickerCmd = fmt.Sprintf(`%s %s %s %s %s -x %s -s %s -e %s -i %d -k %s -t %s 2>/dev/null &`,
			PyAbsPath,
			SparkAbsPath,
			SparkParaList,
			pickerETCDKey,
			PickerScrAbsPath,
			fileType,
			ugc.Start,
			ugc.End,
			ugc.Id,
			ugc.Topic,
			"rule")
	}
	fmt.Println(pickerETCDKey, startPickerCmd)

	/******构建：创建和删除topic命令******/
	creatTopicCmd := fmt.Sprintf(`kafka-topics --create --zookeeper %s --replication-factor 3 --partitions 1 --topic %s`, KafkaTopicIpPort, topicName)
	delTopicCmd := fmt.Sprintf(`kafka-topics --zookeeper %s --topic %s --delete`, KafkaTopicIpPort, topicName)
	/*********************************/

	/******创建topic******/
	fmt.Println("######创建topic######")
	err = SSHCmd(TopicSSHUser, TopicSSHPass, TopicSSHIP, creatTopicCmd, SSHPort)
	if err != nil {
		fmt.Println("creat topic error is :", err)
	}
	/********************/

	/******创建picker etcd******/
	fmt.Println("######创建 picker_etcd!######")
	_, err = EtcdCmd("put", pickerETCDKey, "", PickerETCDIpPort)

	if err != nil {
		mlog.Debug("pickeretcd creat error:", err)
		res_t.Result = "faild"
		delerr := SSHCmd(TopicSSHUser, TopicSSHPass, TopicSSHIP, delTopicCmd, SSHPort)
		if delerr != nil {
			mlog.Debug("pickeretcd creat error and delete topic ", topicName, " error:", delerr)
		}
		return err, &res_t
	}
	fmt.Println("creat picker_etcd finish! err=", err)
	/****************************/

	/******调度agent******/
	fmt.Println("######调度 agent start######")
	agentEtcdCmdKey := fmt.Sprintf("%s/%d", AgentETCDCmdKey, taskID)
	agentPar.SignalType = "start"
	paraAgent, err := json.Marshal(agentPar)
	if err != nil {
		mlog.Debug("task ", taskID, " jsonMarshal agentPar.SignalType=stop error:", err)
	}
	_, err = EtcdCmd("put", agentEtcdCmdKey, string(paraAgent), AgentETCDCmdIpPort)
	if err != nil {
		res_t.Result = "faild"
		statuserr := this.UpgradeStatus("status", "error", taskID, para.OfflineTag)
		if statuserr != nil {
			mlog.Debug("agent start error and update status error to task ", taskID, " error")
		}
		detailserr := this.UpgradeStatus("details", "start agent faild", taskID, para.OfflineTag)
		if detailserr != nil {
			mlog.Debug("agent start error and update details error to task ", taskID, " error")
		}
		return err, &res_t
	}
	/*********************/

	//向agent发送"start"消息
	agentPar.SignalType = "start"
	paraAgent, err = json.Marshal(agentPar)
	if nil != err {
		mlog.Debug("`start` json.Marshal err")
	}
	err = SendOfflineMsg(paraAgent)
	if nil != err {
		mlog.Debug("send `start` msg failed")
	}

	/*******监控 picker_etcd*******/
	fmt.Println("######监控picker_etcd start######")
	go func() {
		agentPar.SignalType = "stop"
		agentStop, err := json.Marshal(agentPar)
		if err != nil {
			mlog.Debug("task ", taskID, " jsonMarshal agentPar.SignalType=stop error:", err)
		}
		this.WatchEtcdPicker(pickerETCDKey, PickerETCDIpPort, string(agentStop), taskID, para.OfflineTag)
	}()
	/*****************************/

	/******启动 picker******/
	fmt.Println("######启动 picker######")
	go func() {
		err = SSHCmd(PickerSSHUser, PickerSSHPass, PickerSSHIP, startPickerCmd, SSHPort)
		if err != nil {
			fmt.Println("start picker error is:", err)
		}
	}()
	/**********************/

	/******监控agent_etcd******/
	fmt.Println("######监控agent_etcd start######")
	agentListNum := len(AgentStatusETCDSlice)
	chAgent := make(chan int, agentListNum)
	for index, agentKey := range AgentStatusETCDSlice {
		go this.WatchEtcdAgent(taskType, agentKey, AgentStatusETCDIpPort[index], topicName, agentEtcdCmdKey, taskID, chAgent)
	}
	go func() {
		var chCount int
		i := 0
		for res := range chAgent {
			chCount += res
			i++
			if i == agentListNum {
				close(chAgent)
				if chCount == agentListNum {
					time.Sleep(10 * time.Minute)
					switch taskType {
					case "vds":
						verr := DuplicateVds("alert_vds", para.Name, para.Time)
						if verr != nil {
							mlog.Debug("duplicate alert_vds error")
						}
						voerr := DuplicateVds("alert_vds_offline", para.Name, para.Time)
						if voerr != nil {
							mlog.Debug("duplicate alert_vds_offline error")
						}
					case "waf":
						werr := DuplicateWaf("alert_waf", para.Name, para.Time)
						if werr != nil {
							mlog.Debug("duplicate alert_waf error")
						}
						woerr := DuplicateWaf("alert_waf_offline", para.Name, para.Time)
						if woerr != nil {
							mlog.Debug("duplicate alert_waf_offline error")
						}
					}
					err := this.UpgradeStatus("status", "complete", taskID, para.OfflineTag)
					if err != nil {
						mlog.Debug("task ", taskID, " complete but update status error!")
					}

					fmt.Println("finish")
					agentPar.SignalType = "complete"
				} else if chCount > agentListNum {
					fmt.Println("task ", taskID, " shutdown!")
					goto WATCHAGENTOUT
				} else if chCount < agentListNum {
					agentPar.SignalType = "shutdown"
					mlog.Debug("task ", taskID, " WatchEtcdAgent error!")
					upStaErr := this.UpgradeStatus("status", "error", taskID, para.OfflineTag)
					if upStaErr != nil {
						mlog.Debug("task ", taskID, " WatchEtcdAgent error and update status error!")
					}
					upDetErr := this.UpgradeStatus("details", "WatchEtcdAgent Error!", taskID, para.OfflineTag)
					if upDetErr != nil {
						mlog.Debug("task ", taskID, " WatchEtcdAgent error and update status error!")
					}
				}

				//向agent发送'complete'消息
				agentCmdComplete, err := json.Marshal(agentPar)
				if nil != err {
					mlog.Debug("`complete` json.Marshal err")
				}
				err = SendOfflineMsg(agentCmdComplete)
				if nil != err {
					mlog.Debug("send `complete` msg failed")
				}

				err = SSHCmd(TopicSSHUser, TopicSSHPass, TopicSSHIP, delTopicCmd, SSHPort)
				if err != nil {
					mlog.Debug("delete topic error!")
				} else {
					mlog.Debug("delete topic ", topicName, " OK!")
				}
				_, err = EtcdCmd("delete", pickerETCDKey, "", PickerETCDIpPort)
				if err != nil {
					mlog.Debug("delete picker etcd error!")
				}
				_, err = EtcdCmd("delete", agentEtcdCmdKey, "", AgentETCDCmdIpPort)
				if err != nil {
					mlog.Debug("delete agent cmd etcd error!")
				}
				/*********************/
			}
		}
	WATCHAGENTOUT:
	}()
	/*************************/
	res_t.Result = "ok"
	return nil, &res_t
}

func (this *TblOLA) UpgradeStatus(column, value string, taskID int, tableTag string) error {
	query := fmt.Sprintf(`update %s set %s='%s' where id=%d;`,
		this.TableName(tableTag),
		column,
		value,
		taskID)
	rows, err := db.Query(query)
	fmt.Println(query)
	if err != nil {
		mlog.Debug("error!agent upgrade status faild!")
		return err
	}
	defer rows.Close()
	return nil
}
