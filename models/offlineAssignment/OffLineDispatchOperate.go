/********离线调度********/
package offlineAssignment

import (
	"apt-web-server_v2/models/db"
	"apt-web-server_v2/models/vds"
	"apt-web-server_v2/models/waf"
	"apt-web-server_v2/modules/mconfig"
	"apt-web-server_v2/modules/mlog"
	"encoding/json"
	"fmt"
	"regexp"
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
	Reg                   = regexp.MustCompile(`[\S]+`)
	keyslice, _           = mconfig.Conf.RawString("agent", "StatusKey")
	ipportslice, _        = mconfig.Conf.RawString("agent", "StatusIpPort")
	AgentStatusETCDSlice  = Reg.FindAllString(keyslice, -1)
	AgentStatusETCDIpPort = Reg.FindAllString(ipportslice, -1)
)

var (
	//Abs:"Absolute"
	PyAbsPath, _           = mconfig.Conf.String("picker", "PyAbsPath")
	SparkAbsPath, _        = mconfig.Conf.String("picker", "SparkAbsPath")
	SparkParaList, _       = mconfig.Conf.String("picker", "SparkParaList")
	PickerScrAbsPath, _    = mconfig.Conf.String("picker", "PickerScrAbsPath")
	AgentETCDCmdKey, _     = mconfig.Conf.String("agent", "CmdKey")
	AgentETCDCmdIpPort, _  = mconfig.Conf.String("agent", "CmdIpPort")
	KafkaCmd, _            = mconfig.Conf.String("topic", "KafkaCmd")
	KafkaCreate, _         = mconfig.Conf.String("topic", "Create")
	KafkaDelete, _         = mconfig.Conf.String("topic", "Delete")
	ParaZookeeper, _       = mconfig.Conf.String("topic", "ZKP")
	ParaRF, _              = mconfig.Conf.String("topic", "RF") //replication-factor
	RFNum, _               = mconfig.Conf.Int("topic", "RFNum")
	ParaPartition, _       = mconfig.Conf.String("topic", "Partition")
	PartitionNum, _        = mconfig.Conf.Int("topic", "PNum")
	ParaTopic, _           = mconfig.Conf.String("topic", "PTopic")
	KafkaTopicIpPort, _    = mconfig.Conf.String("topic", "IpPort")
	PickerETCDStatusKey, _ = mconfig.Conf.String("picker", "ETCDStatusKey")
	ShutdownPicker, _      = mconfig.Conf.String("picker", "ShutdownPicker")
	PickerETCDIpPort, _    = mconfig.Conf.String("picker", "ETCDIpPort")
	TopicSSHUser, _        = mconfig.Conf.String("ssh", "TopicSSHUser")
	TopicSSHPass, _        = mconfig.Conf.String("ssh", "TopicSSHPass")
	TopicSSHIP, _          = mconfig.Conf.String("ssh", "TopicSSHIP")
	PickerSSHUser, _       = mconfig.Conf.String("ssh", "PickerSSHUser")
	PickerSSHPass, _       = mconfig.Conf.String("ssh", "PickerSSHPass")
	PickerSSHIP, _         = mconfig.Conf.String("ssh", "PickerSSHIP")
	SSHPort, _             = mconfig.Conf.Int("ssh", "SSHPort")
	TOPIC, _               = mconfig.Conf.String("agent", "Topic")
	PARTITION, _           = mconfig.Conf.Int("agent", "PartitionNum")
	KAFKA, _               = mconfig.Conf.String("agent", "KafkaIP")
	e                      = mconfig.Conf.AddOption("a", "b", "c")
)

func (this *TblOLA) TableName(tag string) string {
	switch tag {
	case "rule":
		return "ofl_rule"
	default:
		return "ofl_task"
	}
}
func CheckTaskName(tblName, taskName string) string {
	var name string
	query := fmt.Sprintf(`SELECT name FROM %s WHERE name IN ('%s');`, tblName, taskName)
	rows, err := db.DB.Query(query)
	if err != nil {
		mlog.Debug(query, "Check repeated taskName error:", err)
		return "err"
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(
			&name)
		if err != nil {
			return "err"
		}
	}
	if name != "" {
		return "repeated"
	}
	return ""
}
func (this *TblOLA) CreatAssignment(para *TblOLASearchPara) (error, *CMDResult) {
	var res_t CMDResult
	checkName := CheckTaskName(this.TableName(para.OfflineTag), para.Name)
	if checkName != "" {
		switch checkName {
		case "err":
			res_t.Result = "check taskname error"
		case "repeated":
			res_t.Result = "exist"
		default:
			res_t.Result = "faild"
		}
		return nil, &res_t
	}
	if para.Weight == 0 {
		para.Weight = 5
	}

	fmt.Println("CreatAssignment Enter, para is :", para)

	//var ele [5]string
	//for idx, _ := range para.RuleSet {
	//	if len(ele) > 0 {
	//		ele[idx] = para.RuleSet[idx].Rule
	//	}
	//}
	//fmt.Println("CreatAssignment, ele is :", ele)

	paraTopic := fmt.Sprintf(`%s%d%s`, para.Type, para.Time, para.Name)
	//query := fmt.Sprintf(`insert into %s(name,rule,rule2,rule3,rule4,rule5,time,type,start,end,weight,topic,status,details)
	//					value('%s',%d,'%s','%s','%s',%d,'%s','%s','%s');`,
	//query := fmt.Sprintf(`insert into %s(name,rule,rule2,rule3,rule4,rule5,time,type,start,end,weight,topic,status,details)
	//					value('%s','%s','%s','%s','%s','%s',%d,'%s','%s','%s',%d,'%s','%s','%s');`,
	//	/*"offline_assignment_rule2",*/ this.TableName(para.OfflineTag),

	//rslice := make([]string, 0)
	//for idx, _ := range para.RuleSet {
	//	if len(para.RuleSet[idx].Rule) > 0 {
	//		singlrStr := fmt.Sprintf(`%d:%s`, para.RuleSet[idx].Id, para.RuleSet[idx].Rule)
	//		rslice = append(rslice, singlrStr)
	//	}
	//}
	//VarSetStr := strings.Join(rslice, "|")

	query := fmt.Sprintf(`insert into %s(name,ruleset,time,type,start,end,weight,topic,status,details) 
						value('%s','%s',%d,'%s','%s','%s',%d,'%s','%s','%s');`,
		this.TableName(para.OfflineTag), /*"offline_assignment_rule2",*/
		para.Name,
		para.RuleSet, /*VarSetStr*/
		para.Time,
		para.Type,
		para.Start,
		para.End,
		para.Weight,
		paraTopic,
		"ready",
		para.Details)

	//query := fmt.Sprintf(`insert into %s(name,time,type,start,end,weight,topic,status,details)
	//					value('%s',%d,'%s','%s','%s',%d,'%s','%s','%s');`,
	//	this.TableName(para.OfflineTag),
	//	para.Name,
	//	para.Time,
	//	para.Type,
	//	para.Start,
	//	para.End,
	//	para.Weight,
	//	paraTopic,
	//	"ready",
	//	para.Details)

	rows, err := db.DB.Query(query)
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
	query := fmt.Sprintf(`delete from %s where name='%s' and time=%d;`,
		this.TableName(para.OfflineTag), para.Name, para.Time)
	rows, err := db.DB.Query(query)
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
	var res_t CMDResult
	var agentPar AgentPara
	query := fmt.Sprintf(`SELECT id,start,end,type,weight,topic,status 
	    FROM %s WHERE name='%s' AND time=%d;`,
		this.TableName(para.OfflineTag), para.Name, para.Time)
	rows, err := db.DB.Query(query)
	if err != nil {
		res_t.Result = "faild"
		mlog.Debug(query, "StartAssignment get task status error")
		return err, &res_t
	}
	defer rows.Close()
	var pickerETCDKey, startPickerCmd, fileType, topicName, taskType string
	var taskID int

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
		taskID = ugc.Id
		taskType = ugc.Type
		topicName = ugc.Topic
		agentPar.Weight = ugc.Weight
		agentPar.Engine = taskType
		agentPar.Topic = topicName

		/******判断任务运行状态，避免重复运行******/
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
			PyAbsPath, SparkAbsPath, SparkParaList, pickerETCDKey, PickerScrAbsPath,
			fileType, ugc.Start, ugc.End, ugc.Id, ugc.Topic, para.OfflineTag)
	}
	fmt.Println(pickerETCDKey, startPickerCmd)

	/******构建：创建和删除topic命令******/
	//creatTopicCmd := fmt.Sprintf(`kafka-topics --create --zookeeper %s --replication-factor 3 --partitions 1 --topic %s`, KafkaTopicIpPort, topicName)
	//delTopicCmd := fmt.Sprintf(`kafka-topics --zookeeper %s --topic %s --delete`, KafkaTopicIpPort, topicName)
	creatTopicCmd := fmt.Sprintf(`%s %s %s %s %s %d %s %d %s %s`,
		KafkaCmd, KafkaCreate, ParaZookeeper, KafkaTopicIpPort, ParaRF, RFNum,
		ParaPartition, PartitionNum, ParaTopic, topicName)
	delTopicCmd := fmt.Sprintf(`%s %s %s %s %s %s`, KafkaCmd, ParaZookeeper,
		KafkaTopicIpPort, ParaTopic, topicName, KafkaDelete)
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
					if para.OfflineTag != "rule" {
						time.Sleep(10 * time.Minute)
					}
					OflDup(taskType, para)
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
				/******通知agent任务完成******/
				agentCmdComplete, err := json.Marshal(agentPar)
				if nil != err {
					mlog.Debug("`complete` json.Marshal err")
				}
				err = SendOfflineMsg(agentCmdComplete)
				if nil != err {
					mlog.Debug("send `complete` msg failed")
				}
				/******删除topic和etcd******/
				DelTopic(delTopicCmd, topicName)
				DelETCD(pickerETCDKey, agentEtcdCmdKey)
				/*********************/
			}
		}
	WATCHAGENTOUT:
	}()
	/*************************/
	res_t.Result = "ok"
	return nil, &res_t
}
func OflDup(taskType string, para *TblOLASearchPara) {
	switch taskType {
	case "vds":
		verr := vds.DuplicateVds("alert_vds", para.Name, para.Time)
		if verr != nil {
			mlog.Debug("duplicate alert_vds error")
		}
		voerr := vds.DuplicateVds("alert_vds_offline", para.Name, para.Time)
		if voerr != nil {
			mlog.Debug("duplicate alert_vds_offline error")
		}
	case "waf":
		werr := waf.DuplicateWaf("alert_waf", para.Name, para.Time)
		if werr != nil {
			mlog.Debug("duplicate alert_waf error")
		}
		woerr := waf.DuplicateWaf("alert_waf_offline", para.Name, para.Time)
		if woerr != nil {
			mlog.Debug("duplicate alert_waf_offline error")
		}
	}
}
func DelTopic(delTopicCmd, topicName string) {
	err := SSHCmd(TopicSSHUser, TopicSSHPass, TopicSSHIP, delTopicCmd, SSHPort)
	if err != nil {
		mlog.Debug("delete topic error!")
	} else {
		mlog.Debug("delete topic ", topicName, " OK!")
	}
}
func DelETCD(pickerETCDKey, agentEtcdCmdKey string) {
	_, err := EtcdCmd("delete", pickerETCDKey, "", PickerETCDIpPort)
	if err != nil {
		mlog.Debug("delete picker etcd error!")
	}
	_, err = EtcdCmd("delete", agentEtcdCmdKey, "", AgentETCDCmdIpPort)
	if err != nil {
		mlog.Debug("delete agent cmd etcd error!")
	}
}

func (this *TblOLA) UpgradeStatus(column, value string, taskID int, tag string) error {
	query := fmt.Sprintf(`update %s set %s='%s' where id=%d;`,
		this.TableName(tag),
		column,
		value,
		taskID)
	rows, err := db.DB.Query(query)
	fmt.Println(query)
	if err != nil {
		mlog.Debug("error!agent upgrade status faild!")
		return err
	}
	defer rows.Close()
	return nil
}
