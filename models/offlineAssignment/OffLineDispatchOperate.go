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

/************config************/
var (
	Reg                = regexp.MustCompile(`[\S]+`)
	AgentIpList, _     = mconfig.Conf.RawString("agent", "Ip")
	AgentIp            = Reg.FindAllString(AgentIpList, -1)
	AgentStatList, _   = mconfig.Conf.RawString("agent", "StatusKey")
	AgentStatKey       = Reg.FindAllString(AgentStatList, -1)
	EtcdIpList, _      = mconfig.Conf.RawString("etcd", "Ip")
	EtcdIp             = Reg.FindAllString(EtcdIpList, -1)
	EtcdPort, _        = mconfig.Conf.RawString("etcd", "Port")
	EtcdIpPortList, _  = mconfig.Conf.RawString("etcd", "IpPort")
	EtcdIpPort         = Reg.FindAllString(EtcdIpPortList, -1)
	KafkaIpList, _     = mconfig.Conf.RawString("kafka", "Ip")
	KafkaIp            = Reg.FindAllString(KafkaIpList, -1)
	KafkaPort, _       = mconfig.Conf.RawString("kafka", "Port")
	KafkaIpPortList, _ = mconfig.Conf.RawString("kafka", "IpPort")
	KafkaIpPort        = Reg.FindAllString(KafkaIpPortList, -1)
)

var (
	//Abs:"Absolute"
	PyAbsPath, _        = mconfig.Conf.String("picker", "PyAbsPath")
	SparkAbsPath, _     = mconfig.Conf.String("picker", "SparkAbsPath")
	SparkParaList, _    = mconfig.Conf.String("picker", "SparkParaList")
	PickerScrAbsPath, _ = mconfig.Conf.String("picker", "PickerScrAbsPath")
	AgentETCDCmdKey, _  = mconfig.Conf.String("agent", "CmdKey")
	KafkaCmd, _         = mconfig.Conf.String("kafka", "KafkaCmd")
	RFNum, _            = mconfig.Conf.Int("kafka", "RFNum")
	PartitionNum, _     = mconfig.Conf.Int("kafka", "PNum")
	ParaTopic, _        = mconfig.Conf.String("kafka", "PTopic")
	PickerStatusKey, _  = mconfig.Conf.String("picker", "ETCDStatusKey")
	ShutdownPicker, _   = mconfig.Conf.String("picker", "ShutdownPicker")
	TopicSSHUser, _     = mconfig.Conf.String("ssh", "TopicSSHUser")
	TopicSSHPass, _     = mconfig.Conf.String("ssh", "TopicSSHPass")
	TopicSSHIP, _       = mconfig.Conf.String("ssh", "TopicSSHIP")
	PickerSSHUser, _    = mconfig.Conf.String("ssh", "PickerSSHUser")
	PickerSSHPass, _    = mconfig.Conf.String("ssh", "PickerSSHPass")
	PickerSSHIP, _      = mconfig.Conf.String("ssh", "PickerSSHIP")
	SSHPort, _          = mconfig.Conf.Int("ssh", "SSHPort")
	TOPIC, _            = mconfig.Conf.String("agent", "Topic")
	PARTITION, _        = mconfig.Conf.Int("agent", "PartitionNum")
	e                   = mconfig.Conf.AddOption("a", "b", "c")
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
	paraTopic := fmt.Sprintf(`%s%d%s`, para.Type, para.Time, para.Name)
	query := fmt.Sprintf(`INSERT INTO %s(name,ruleset,time,type,start,end,weight,topic,status,details) 
						value('%s','%s',%d,'%s','%s','%s',%d,'%s','%s','%s');`,
		this.TableName(para.OfflineTag),
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
	query := fmt.Sprintf(`DELETE FROM %s WHERE name='%s' AND time=%d;`,
		this.TableName(para.OfflineTag), para.Name, para.Time)
	rows, err := db.DB.Query(query)
	if err != nil {
		mlog.Debug(query, "DeleteAssignment error")
		res_t.Result = "faild"
		return err, &res_t
	}
	defer rows.Close()

	res_t.Result = "ok"
	return nil, &res_t
}
func (this *TblOLA) GetTaskMsg(para *TblOLASearchPara) error {
	query := fmt.Sprintf(`SELECT id,start,end,type,weight,topic,status 
	    FROM %s WHERE name='%s' AND time=%d;`,
		this.TableName(para.OfflineTag), para.Name, para.Time)
	rows, err := db.DB.Query(query)
	if err != nil {
		mlog.Debug(query, "StartAssignment get task status error")
		return err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(
			&this.Id,
			&this.Start,
			&this.End,
			&this.Type,
			&this.Weight,
			&this.Topic,
			&this.Status)
		if err != nil {
			return err
		}
	}
	return nil
}
func GetFileType(taskType string) string {
	var fileType string
	switch taskType {
	case "vds":
		fileType = "file"
	case "waf":
		fileType = "http"
	case "rule":
		fileType = "http"
	}
	return fileType
}
func (this *TblOLA) CheckTaskStatus(oflTag string) (string, error) {
	if this.Status == "running" {
		return "this task is running!", nil
	} else if this.Status == "ready" || this.Status == "error" {
		err := this.UpgradeStatus("status", "running", this.Id, oflTag)
		if err != nil {
			mlog.Debug("Update status task ", this.Id, "to running error")
			return "Update task to running error", err
		}
	}
	return "", nil
}
func (this *TblOLA) StartAssignment(para *TblOLASearchPara) (error, *CMDResult) {
	var res_t CMDResult
	var agentPar AgentPara

	err := this.GetTaskMsg(para)
	if err != nil {
		res_t.Result = "faild"
		mlog.Debug("Get task message error")
		return err, &res_t
	}
	agentPar.Weight = this.Weight
	agentPar.Engine = this.Type
	agentPar.Topic = this.Topic
	fileType := GetFileType(this.Type)
	/******判断任务运行状态，避免重复运行******/
	res_t.Result, err = this.CheckTaskStatus(para.OfflineTag)
	if res_t.Result != "" || err != nil {
		return err, &res_t
	}
	/************************************/
	pickerKey := fmt.Sprintf("%s/%d", PickerStatusKey, this.Id)
	mlog.Debug(pickerKey)
	agentCmdKey := fmt.Sprintf("%s/%d", AgentETCDCmdKey, this.Id)
	mlog.Debug(agentCmdKey)
	/******创建topic******/
	err, res_t.Result = this.TopicCreate()
	if nil != err || "" != res_t.Result {
		fmt.Println("create topic ", this.Topic, " error:", err)
	}
	/******创建picker etcd******/
	err, res_t.Result = this.PickerEtcdCreate(pickerKey)
	if nil != err || "" != res_t.Result {
		DelTopic(this.Topic)
		return err, &res_t
	}
	/******调度agent******/
	err, res_t.Result = this.AgentTell(agentPar, agentCmdKey, para.OfflineTag)
	if "" != res_t.Result || nil != err {
		DelTopic(this.Topic)
		DelETCD(pickerKey, agentCmdKey)
		return err, &res_t
	}
	/******监控 picker_etcd******/
	go this.PickerWatch(agentPar, pickerKey, para.OfflineTag)
	/******启动 picker******/
	go this.PickerStart(pickerKey, fileType, para.OfflineTag)
	/******监控agent_etcd******/
	go this.AgentWatch(para, agentPar, pickerKey, agentCmdKey)

	res_t.Result = "ok"
	return nil, &res_t
}
func (this *TblOLA) TopicCreate() (error, string) {
	/*创建topic命令
	kafka-topics --create          指定创建动作
	--zookeeper KafkaIpPort        指定kafka连接zookeeper的连接url
	--replication-factor RFNum     指定每个分区的复制因子个数
	--partitions len(AgentIp)      指定当前创建的kafka分区数量
	--topic topicName              指定要创建的topic的名称*/
	creatTopicCmd := fmt.Sprintf(`%s --create --zookeeper %s --replication-factor %d --partitions %d --topic %s`,
		KafkaCmd, KafkaIpPort[0], RFNum, len(AgentIp), this.Topic)
	mlog.Debug("######创建topic######", creatTopicCmd)
	err := SSHCmd(TopicSSHUser, TopicSSHPass, TopicSSHIP, creatTopicCmd, SSHPort)
	if err != nil {
		mlog.Debug("creat topic error is :", err)
		return err, "create topic error"
	}
	return nil, ""
}
func (this *TblOLA) PickerEtcdCreate(pickerKey string) (error, string) {
	mlog.Debug("######创建pickerEtcd######")
	_, err := EtcdCmd("put", pickerKey, "")
	if err != nil {
		mlog.Debug("pickeretcd create error:", err)
		return err, "create picker etcd error"
	}
	return nil, ""
}
func (this *TblOLA) AgentTell(agentPar AgentPara, agentCmdKey, oflTag string) (error, string) {
	mlog.Debug("######调度 agent start######")
	agentPar.SignalType = "start"
	paraAgent, err := json.Marshal(agentPar)
	if err != nil {
		mlog.Debug("task ", this.Id, " jsonMarshal agentPar.SignalType=stop error:", err)
		return err, "make json to start agent error"
	}
	_, err = EtcdCmd("put", agentCmdKey, string(paraAgent))
	if err != nil {
		mlog.Debug("task", this.Id, "put start to agentetcd error")
		statuserr := this.UpgradeStatus("status", "error", this.Id, oflTag)
		if statuserr != nil {
			mlog.Debug("agent start error and update status error to task ", this.Id, " error")
		}
		detailserr := this.UpgradeStatus("details", "start agent faild", this.Id, oflTag)
		if detailserr != nil {
			mlog.Debug("agent start error and update details error to task ", this.Id, " error")
		}
		return err, "put start json to agentetcd error"
	}
	/******向agent发送start消息******/
	mlog.Debug("######向agentKafka发送start消息######")
	err = SendOfflineMsg(paraAgent)
	if nil != err {
		mlog.Debug("send `start` msg failed")
		return err, "send start to agentKafka error"
	}
	return nil, ""
}
func (this *TblOLA) PickerStart(pickerKey, fileType, oflTag string) {
	mlog.Debug("######启动 picker######")
	startPickerCmd := fmt.Sprintf(`%s %s %s %s %s -x %s -s %s -e %s -i %d -k %s -t %s 2>/dev/null &`,
		PyAbsPath, SparkAbsPath, SparkParaList, pickerKey, PickerScrAbsPath,
		fileType, this.Start, this.End, this.Id, this.Topic, oflTag)
	mlog.Debug(startPickerCmd)
	err := SSHCmd(PickerSSHUser, PickerSSHPass, PickerSSHIP, startPickerCmd, SSHPort)
	if err != nil {
		mlog.Debug("start picker error is:", err)
	}
}
func (this *TblOLA) PickerWatch(agentPar AgentPara, pickerKey, oflTag string) {
	mlog.Debug("######监控picker_etcd start######")
	agentPar.SignalType = "stop"
	agentStop, err := json.Marshal(agentPar)
	if err != nil {
		mlog.Debug("task ", this.Id, " make agent stop json error:", err)
	}
	this.WatchEtcdPicker(pickerKey, string(agentStop), oflTag)
}
func (this *TblOLA) AgentWatch(para *TblOLASearchPara, agentPar AgentPara, pickerKey, agentCmdKey string) {
	agentListNum := len(AgentIp)
	chAgent := make(chan int, agentListNum)
	for index, StatKey := range AgentStatKey {
		mlog.Debug("######监控agent_etcd start######", index)
		go this.WatchEtcdAgent(this.Type, StatKey, agentCmdKey, chAgent)
	}
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
				OflDup(this.Type, para)
				err := this.UpgradeStatus("status", "complete", this.Id, para.OfflineTag)
				if err != nil {
					mlog.Debug("task ", this.Id, " complete but update status error!")
				}
				mlog.Debug(this.Name, "finished")
				agentPar.SignalType = "complete"
			} else if chCount > agentListNum {
				mlog.Debug("task ", this.Id, " shutdown!")
				goto WATCHAGENTOUT
			} else if chCount < agentListNum {
				agentPar.SignalType = "error"
				mlog.Debug("task ", this.Id, " WatchEtcdAgent error!")
				upStaErr := this.UpgradeStatus("status", "error", this.Id, para.OfflineTag)
				if upStaErr != nil {
					mlog.Debug("task ", this.Id, " WatchEtcdAgent error and update status error!")
				}
				upDetErr := this.UpgradeStatus("details", "WatchEtcdAgent Error!", this.Id, para.OfflineTag)
				if upDetErr != nil {
					mlog.Debug("task ", this.Id, " WatchEtcdAgent error and update status error!")
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
			DelTopic(this.Topic)
			DelETCD(pickerKey, agentCmdKey)
			/*********************/
		}
	}
WATCHAGENTOUT:
}
func OflDup(taskType string, para *TblOLASearchPara) error {
	var olerr, oferr error
	switch taskType {
	case "vds":
		olerr = vds.DuplicateVds("alert_vds", para.Name, para.Time)
		if olerr != nil {
			mlog.Debug("duplicate alert_vds error")
		}
		oferr = vds.DuplicateVds("alert_vds_offline", para.Name, para.Time)
		if oferr != nil {
			mlog.Debug("duplicate alert_vds_offline error")
		}
	case "waf":
		olerr = waf.DuplicateWaf("alert_waf", para.Name, para.Time)
		if olerr != nil {
			mlog.Debug("duplicate alert_waf error")
		}
		oferr = waf.DuplicateWaf("alert_waf_offline", para.Name, para.Time)
		if oferr != nil {
			mlog.Debug("duplicate alert_waf_offline error")
		}
	}
	if nil != olerr {
		return olerr
	} else if nil != oferr {
		return oferr
	}
	return nil
}
func DelTopic(topicName string) error {
	/*kafka-topics --zookeeper KafkaTopicIpPort --topic topicName --delete*/
	delTopicCmd := fmt.Sprintf(`%s --zookeeper %s --topic %s --delete`,
		KafkaCmd, KafkaIpPort[0], topicName)
	mlog.Debug(delTopicCmd)
	err := SSHCmd(TopicSSHUser, TopicSSHPass, TopicSSHIP, delTopicCmd, SSHPort)
	if err != nil {
		mlog.Debug("delete topic error!")
		return err
	} else {
		mlog.Debug("delete topic ", topicName, " OK!")
	}
	return nil
}

func DelETCD(pickerKey, agentCmdKey string) error {
	err := SSHCmd(PickerSSHUser, PickerSSHPass, PickerSSHIP, ShutdownPicker+" "+pickerKey, SSHPort)
	if err != nil {
		mlog.Debug("stop ", pickerKey, " error!", err)
	} else {
		mlog.Debug("stop ", pickerKey, " OK!", err)
	}
	_, err = EtcdCmd("delete", pickerKey, "")
	if err != nil {
		mlog.Debug("delete picker etcd error!")
	}
	_, err = EtcdCmd("delete", agentCmdKey, "")
	if err != nil {
		mlog.Debug("delete agent cmd etcd error!")
	}
	return nil
}

func (this *TblOLA) UpgradeStatus(column, value string, taskID int, tag string) error {
	query := fmt.Sprintf(`UPDATE %s SET %s='%s' WHERE id=%d;`,
		this.TableName(tag),
		column,
		value,
		taskID)
	rows, err := db.DB.Query(query)
	//fmt.Println(query)
	if err != nil {
		mlog.Debug("error!agent upgrade status faild!")
		return err
	}
	defer rows.Close()
	return nil
}
