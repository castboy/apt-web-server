package offlineAssignment

import (
	"apt-web-server_v2/modules/mconfig"
	"apt-web-server_v2/modules/mlog"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	Status                   [3]map[string]ETCDAgent
}

func EtcdCmd(cmd, key, value, ipPort string) (rtn string, err error) {
	etcdError := fmt.Sprintf("etcd ", cmd, "error!")
	etcdCmdError := errors.New(etcdError)
	defer func() {
		if info := recover(); info != nil {
			//log.Fatal("etcd ", cmd, " panic:", info, "context.backgroud=", context.Background(), "key=", key)
			mlog.Debug("etcd ", cmd, " error:", info)
			err = etcdCmdError
		}
	}()

	cfg := clientv3.Config{
		Endpoints:   []string{ipPort},
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

func GetEngenStatus(agentType string, agentstatus []byte, topicName string) (int, int, error) {
	var s StatusFromEtcd
	var Waf = make(map[string]ETCDAgent, 1000)
	var Vds = make(map[string]ETCDAgent, 1000)
	var Rule = make(map[string]ETCDAgent, 1000)

	err := json.Unmarshal(agentstatus, &s)
	if err != nil {
		mlog.Debug("GettEngenStatus json Unmarshal Err")
		return -1, -1, err
	}

	Waf = s.Status[0]
	Vds = s.Status[1]
	Rule = s.Status[2]

	fmt.Println("Vds[topicName]=", Vds[topicName], "Waf[topicName]=", Waf[topicName])
	if agentType == "vds" {
		return Vds[topicName].Engine + Vds[topicName].Err, Vds[topicName].Last, nil
	}
	if agentType == "rule" {
		return Rule[topicName].Engine + Rule[topicName].Err, Rule[topicName].Last, nil
	}
	return Waf[topicName].Engine + Waf[topicName].Err, Waf[topicName].Last, nil
}

/*
func EtcdPut(key, cmdPara, ipPort string) error {
	defer func() {
		if info := recover(); info != nil {
			log.Fatal("etcdput get panic:", info, "context.backgroud=", context.Background(), "key=", key)
		}
	}()
	cfg := clientv3.Config{
		Endpoints:   []string{ipPort},
		DialTimeout: 5 * time.Second,
	}
	cli, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	fmt.Println(cfg)
	//byte,_ := json.Marshal(cmdPara)
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	_, err = cli.Put(ctx, key, cmdPara)
	cancel()

	if err != nil {

	}
	return err
}

func EtcdGet(key, ipPort string) (string, error) {
	defer func() {
		if info := recover(); info != nil {
			log.Fatal("etcdget get panic:", info, "context.backgroud=", context.Background(), "key=", key)
		}
	}()
	cfg := clientv3.Config{
		Endpoints:   []string{ipPort},
		DialTimeout: 5 * time.Second,
	}
	cli, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	resp, err := cli.Get(ctx, key)
	cancel()
	return string(resp.Kvs[0].Value), err
}
*/
func GetEtcdAgent(agentType, topicName, key, ipPort string) (int, int, error) {
	agentEtcdStr, err := EtcdCmd("get", key, "", ipPort)
	if err != nil {
		mlog.Debug("GetEtcdAgetn's EtcdCmd get from ", ipPort, key, " error:", err)
		return -1, -1, err
	}
	count, total, err := GetEngenStatus(agentType, []byte(agentEtcdStr), topicName)
	if err != nil {
		mlog.Debug("GetEtcdAgent GetEngenStatus error:", err)
		return -1, -1, err
	}
	return count, total, nil
}

func GetEtcdPicker(key, ipPort string) (int, int, error) {
	var pickerJson ETCDPicker
	pickerEtcdStr, err := EtcdCmd("get", key, "", ipPort)
	if err != nil {
		mlog.Debug("GetEtcdPicker's EtcdCmd get from ", ipPort, key, " error:", err)
		return -1, -1, err
	}
	json.Unmarshal([]byte(pickerEtcdStr), &pickerJson)
	return pickerJson.Count, pickerJson.Total, nil
}

func GetEtcdAgentCmd(key, ipPort string) string {
	var agentCmdJson AgentPara
	agentEtcdCmd, _ := EtcdCmd("get", key, "", ipPort)
	json.Unmarshal([]byte(agentEtcdCmd), &agentCmdJson)
	return agentCmdJson.SignalType
}

func (this *TblOLA) WatchEtcdAgent(agentType, key, ipPort, topicName, agentEtcdCmdKey string, taskId int, chAgent chan int) {
	timeOutUnit, _ := mconfig.Conf.Int("time", "Tout")
	defer func() {
		if info := recover(); info != nil {
			//log.Fatal("watchetcdagent get panic:", info, "context.backgroud=", context.Background(), "key=", key, "perfix=", clientv3.WithPrefix())
			mlog.Debug("WatchEtcdAgent ", ipPort, key, " error:", info)
			err := this.UpgradeStatus("status", "error", taskId, agentType)
			if err != nil {
				mlog.Debug("WatchEtcdAgent update status error")
			}
			err = this.UpgradeStatus("details", "watchEtcdAgent error", taskId, agentType)
			if err != nil {
				mlog.Debug("WatchEtcdAgent update details error")
			}
		}
	}()

	fmt.Println("task ", taskId, " start watch etcd agent")
	cfg := clientv3.Config{
		Endpoints:   []string{ipPort},
		DialTimeout: 5 * time.Second,
	}
	cli, err := clientv3.New(cfg)
	if err != nil {
		panic(err)
	}
	defer cli.Close()
	rch := cli.Watch(context.Background(), key, clientv3.WithPrefix())
	out := make(chan bool)
	tout := make(chan bool)
	timeOut := time.Duration(timeOutUnit) * time.Minute
	go func() {
		timed := time.NewTimer(timeOut)
	TAG:
		select {
		case <-timed.C:
			mlog.Debug("etcd time out")
			chAgent <- 0
			tout <- true
			break
		case <-out:
			mlog.Debug("time reset",time.Now())
			timed.Reset(timeOut)
			goto TAG
		}
	}()
	for wresp := range rch {
		for _, ev := range wresp.Events {
			out <- true
			fmt.Println("ev.Kv.Value is:", string(ev.Kv.Value))
			etcdValue := string(ev.Kv.Value)
			count, total, err := GetEngenStatus(agentType, []byte(etcdValue), topicName)
			if err != nil {
				mlog.Debug("WatchEtcdAgent's GetEngenStatus error!")
				chAgent <- 0
				goto BREAKTAG
			}
			cmdStatus := GetEtcdAgentCmd(agentEtcdCmdKey, AgentETCDCmdIpPort)
			if count == total && total != -1 && "stop" == cmdStatus {
				fmt.Println("channel check is stop?", cmdStatus, count, total)
				chAgent <- 1
				goto BREAKTAG
			} else if "shutdown" == cmdStatus || cmdStatus == "" {
				fmt.Println("send shutdown to agent etcd!")
				chAgent <- 2
				fmt.Println("task shutdown!!!!!!!")
				goto BREAKTAG
			}
		}
		select {
		case <-tout:
			goto BREAKTAG
		}
	}
BREAKTAG:
	close(out)
	close(tout)
	fmt.Println("stop AgentETCD watcher")
}

func (this *TblOLA) WatchEtcdPicker(key, ipPort, agentPar string, id int, tableTag string) {
	defer func() {
		if info := recover(); info != nil {
			//log.Fatal("watchetcdpicer get panic:", info, "context.backgroud=", context.Background(), "key=", key, "perfix=", clientv3.WithPrefix())
			mlog.Debug("WatchEtcdPicker ", ipPort, key, " error:", info)
			err := this.UpgradeStatus("status", "error", id, tableTag)
			if err != nil {
				mlog.Debug("WatchEtcdPicker's update status error:", err)
			}
			err = this.UpgradeStatus("details", "watchEtcdPicker error", id, tableTag)
			if err != nil {
				mlog.Debug("WatchEtcdPicker's update details error:", err)
			}
		}
	}()

	cfg := clientv3.Config{
		Endpoints:   []string{ipPort},
		DialTimeout: 5 * time.Second,
	}
	cli, err := clientv3.New(cfg)
	if err != nil {
		panic(err)
		log.Fatal(err)
	}
	defer cli.Close()
	agentEtcdCmdKey := fmt.Sprintf(`%s/%d`, AgentETCDCmdKey, id)
	rch := cli.Watch(context.Background(), key, clientv3.WithPrefix())

	for wresp := range rch {
		for _, ev := range wresp.Events {
			//fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.key, ev.Kv.Value)
			fmt.Println(string(ev.Kv.Value))
			etcdValue := string(ev.Kv.Value)
			result := WatchPicker(etcdValue, agentPar, id)
			cmdStatus := GetEtcdAgentCmd(agentEtcdCmdKey, AgentETCDCmdIpPort)
			if result == 0 || cmdStatus == "shutdown" || cmdStatus == "" {
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
	/*
		query := fmt.Sprintf(`update %s set status='%s' where id=%d;`,
			tableName,
			"running",
			id)
		rows, err := db.Query(query)
		if err != nil {
			mlog.Debug(query, "WatchPicker update running to status faild!")
			//return err
		}
		defer rows.Close()
	*/
	if pickerStr.State == "stop" && pickerStr.Offset == pickerStr.End && pickerStr.Total != 0 {
		//向agent发送"stop"消息
		err := SendOfflineMsg([]byte(agentPar))
		if nil != err {
			mlog.Debug("send `stop` msg failed")
		}
		fmt.Println("WatchPicker finished！")
		return 0
	}
	return 1
}