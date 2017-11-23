package offlineAssignment

import (
	"apt-web-server_v2/modules/mconfig"
	"apt-web-server_v2/modules/mlog"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
)

type Conf struct {
	EngineReqPort int
	MaxCache      int
	Partition     map[string]int32
	Topic         []string
}

type StatusFromEtcd struct {
	ReceivedOfflineMsgOffset int64
	Status                   [3]map[string]*ETCDAgent
}

func EtcdCmd(cmd, key, value string) (rtn string, err error) {
	etcdError := fmt.Sprintf("etcd ", cmd, "error!")
	etcdCmdError := errors.New(etcdError)
	defer func() {
		if info := recover(); info != nil {
			mlog.Debug("etcd ", cmd, " error:", info)
			err = etcdCmdError
		}
	}()

	cfg := clientv3.Config{
		Endpoints:   EtcdIpPort,
		DialTimeout: 5 * time.Second,
	}
	cli, err := clientv3.New(cfg)
	if err != nil {
		panic(err)
		//return rtn, err
	}
	defer cli.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	switch cmd {
	case "put":
		_, err := cli.Put(ctx, key, value)
		cancel()
		if err != nil {
			panic(err)
			//return rtn, err
		}
		return rtn, err
	case "get":
		respGet, err := cli.Get(ctx, key)
		cancel()
		if err != nil {
			panic(err)
		}
		if respGet.Count == 0 {
			return "", err
		}
		rtn = string(respGet.Kvs[0].Value)
		return rtn, err
	case "delete":
		_, err := cli.Delete(ctx, key)
		cancel()
		if err != nil {
			panic(err)
			//return rtn, err
		}
		return rtn, err
	}
	cancel()

	return rtn, err
}

func GetEngenStatus(agentType string, agentstatus []byte, topicName string) (*ETCDAgent, error) {
	var s StatusFromEtcd
	var Waf = make(map[string]*ETCDAgent, 1000)
	var Vds = make(map[string]*ETCDAgent, 1000)
	var Rule = make(map[string]*ETCDAgent, 1000)

	err := json.Unmarshal(agentstatus, &s)
	if err != nil {
		mlog.Debug("GettEngenStatus json Unmarshal Err")
		return nil, err
	}
	Waf = s.Status[0]
	Vds = s.Status[1]
	Rule = s.Status[2]
	fmt.Println("Vds[topicName]=", Vds[topicName], "Waf[topicName]=", Waf[topicName])
	if agentType == "vds" {
		//return Vds[topicName].Engine + Vds[topicName].Err, Vds[topicName].Last, nil
		return Vds[topicName], nil
	}
	if agentType == "rule" {
		//return Rule[topicName].Engine + Rule[topicName].Err, Rule[topicName].Last, nil
		return Rule[topicName], nil
	}
	//return Waf[topicName].Engine + Waf[topicName].Err, Waf[topicName].Last, nil
	return Waf[topicName], nil
}

func GetEtcdAgent(agentType, topicName, key string) (*ETCDAgent, error) {
	agentEtcdStr, err := EtcdCmd("get", key, "")
	if err != nil {
		mlog.Debug("GetEtcdAgetn's EtcdCmd get from ", key, " error:", err)
		//return -1, -1, err
		return nil, err
	}
	//count, total, err := GetEngenStatus(agentType, []byte(agentEtcdStr), topicName)
	aStatus, err := GetEngenStatus(agentType, []byte(agentEtcdStr), topicName)
	if err != nil {
		mlog.Debug("GetEtcdAgent GetEngenStatus error:", err)
		//return -1, -1, err
		return nil, err
	}
	//return count, total, nil
	return aStatus, nil
}

func GetEtcdPicker(key string) (*ETCDPicker, error) {
	var pickerJson ETCDPicker
	pickerEtcdStr, err := EtcdCmd("get", key, "")
	if err != nil {
		mlog.Debug("GetEtcdPicker's EtcdCmd get from ", key, " error:", err)
		return nil, err
	}
	json.Unmarshal([]byte(pickerEtcdStr), &pickerJson)
	return &pickerJson, nil
}

func GetAgentCmd(key string) (*AgentPara, error) {
	var agentCmdJson AgentPara
	agentEtcdCmd, err := EtcdCmd("get", key, "")
	if nil != err {
		return nil, err
	}
	json.Unmarshal([]byte(agentEtcdCmd), &agentCmdJson)
	return &agentCmdJson, nil
}

func (this *TblOLA) WatchEtcdAgent(agentType, key, agentCmdKey string, chAgent chan int) {
	timeOutUnit, err := mconfig.Conf.Int("time", "Tout")
	if nil != err {
		timeOutUnit = 30
	}
	defer func() {
		if info := recover(); info != nil {
			chAgent <- 0
			mlog.Debug("WatchEtcdAgent ", EtcdIpPort, key, " error:", info)
		}
	}()
	cfg := clientv3.Config{
		Endpoints:   EtcdIpPort,
		DialTimeout: 5 * time.Second,
	}
	cli, err := clientv3.New(cfg)
	if err != nil {
		panic(err)
	}
	defer cli.Close()
	rch := cli.Watch(context.Background(), key, clientv3.WithPrefix())
	out := make(chan bool)

	finish := make(chan bool)
	timeOut := time.Duration(timeOutUnit) * time.Minute
	go func() {
		timed := time.NewTimer(timeOut)
	TAG:
		select {
		case <-finish:
			mlog.Debug("offline task ", this.Id, "complete")
		case <-timed.C:
			mlog.Debug("etcd time out")
			panic("timeout")
		case <-out:
			timed.Reset(timeOut)
			goto TAG
		}
	}()
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Println("ev.Kv.Value is:", string(ev.Kv.Value))
			out <- true
			etcdValue := string(ev.Kv.Value)
			//count, total, err := GetEngenStatus(agentType, []byte(etcdValue), this.Topic)
			aStatus, err := GetEngenStatus(agentType, []byte(etcdValue), this.Topic)
			if err != nil {
				mlog.Debug("WatchEtcdAgent's GetEngenStatus error!")
				chAgent <- 0
				goto BREAKTAG
			}
			count := aStatus.Engine + aStatus.Err
			total := aStatus.Last
			agentCmd, _ := GetAgentCmd(agentCmdKey)
			if count == total && total != -1 && "stop" == agentCmd.SignalType {
				finish <- true
				chAgent <- 1
				goto BREAKTAG
			} else if "shutdown" == agentCmd.SignalType || "" == agentCmd.SignalType {
				chAgent <- 2
				goto BREAKTAG
			}
		}
	}
BREAKTAG:
	fmt.Println("stop AgentETCD watcher")
}

func (this *TblOLA) WatchEtcdPicker(key, agentPar string, tableTag string) {
	defer func() {
		if info := recover(); info != nil {
			mlog.Debug("WatchEtcdPicker ", EtcdIpPort, key, " error:", info)
			err := this.UpgradeStatus("status", "error", this.Id, tableTag)
			if err != nil {
				mlog.Debug("WatchEtcdPicker's update status error:", err)
			}
			err = this.UpgradeStatus("details", "watchEtcdPicker error", this.Id, tableTag)
			if err != nil {
				mlog.Debug("WatchEtcdPicker's update details error:", err)
			}
		}
	}()
	cfg := clientv3.Config{
		Endpoints:   EtcdIpPort,
		DialTimeout: 5 * time.Second,
	}
	cli, err := clientv3.New(cfg)
	if err != nil {
		panic(err)
	}
	defer cli.Close()
	agentEtcdCmdKey := fmt.Sprintf(`%s/%d`, AgentETCDCmdKey, this.Id)
	rch := cli.Watch(context.Background(), key, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			//fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.key, ev.Kv.Value)
			fmt.Println("picker status:", string(ev.Kv.Value))
			etcdValue := string(ev.Kv.Value)
			result := WatchPicker(etcdValue, agentPar, this.Id)
			agentCmd, _ := GetAgentCmd(agentEtcdCmdKey)
			if result == 0 || agentCmd.SignalType == "shutdown" || agentCmd.SignalType == "" {
				goto BREAKTAG
			}
		}
	}
BREAKTAG:
	fmt.Println("etcd picker watch end!")
}

func WatchPicker(pickerValue, agentPar string, id int) int {
	var pickerStr ETCDPicker
	json.Unmarshal([]byte(pickerValue), &pickerStr)
	agentCmdKey := fmt.Sprintf(`%s/%d`, AgentETCDCmdKey, id)
	if pickerStr.State == "stop" && pickerStr.Offset == pickerStr.End && pickerStr.Total != 0 {
		//向agent发送"stop"消息
		err := SendOfflineMsg([]byte(agentPar))
		if nil != err {
			mlog.Debug(pickerStr.topic, "send `stop` msg to kafka failed")
		}
		_, err = EtcdCmd("put", agentCmdKey, agentPar)
		if nil != err {
			mlog.Debug(pickerStr.topic, "send stop msg to agent etcd error")
		}
		mlog.Debug("Task", pickerStr.topic, "WatchPicker finished!")
		return 0
	}
	return 1
}
