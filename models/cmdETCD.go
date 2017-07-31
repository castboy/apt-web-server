package models

import (
	"apt-web-server/modules/mlog"
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

func EtcdCmd(cmd, key, value, ipPort string) (rtn string, err error) {
	fmt.Println("key", key)
	fmt.Println("val", value)
	fmt.Println("ipPort", ipPort)
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
	var EngenType [2]map[string]ETCDAgent
	var Waf map[string]ETCDAgent
	var Vds map[string]ETCDAgent

	Waf = make(map[string]ETCDAgent, 1000)
	Vds = make(map[string]ETCDAgent, 1000)
	err := json.Unmarshal(agentstatus, &EngenType)
	if err != nil {
		mlog.Debug("GettEngenStatus json Unmarshal Err")
		return -1, -1, err
	}

	Waf = EngenType[0]
	Vds = EngenType[1]
	fmt.Println("Vds[topicName]=", Vds[topicName], "Waf[topicName]=", Waf[topicName])
	if agentType == "vds" {
		return Vds[topicName].Engine + Vds[topicName].Err, Vds[topicName].Last, nil
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
	defer func() {
		if info := recover(); info != nil {
			//log.Fatal("watchetcdagent get panic:", info, "context.backgroud=", context.Background(), "key=", key, "perfix=", clientv3.WithPrefix())
			mlog.Debug("WatchEtcdAgent ", ipPort, key, " error:", info)
			err := this.UpgradeStatus("status", "error", taskId)
			if err != nil {
				mlog.Debug("WatchEtcdAgent update status error")
			}
			err = this.UpgradeStatus("details", "watchEtcdAgent error", taskId)
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
	for wresp := range rch {
		for _, ev := range wresp.Events {
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
	}
BREAKTAG:
	fmt.Println("stop AgentETCD watcher")
}

func (this *TblOLA) WatchEtcdPicker(key, ipPort, agentPar string, id int) {
	defer func() {
		if info := recover(); info != nil {
			//log.Fatal("watchetcdpicer get panic:", info, "context.backgroud=", context.Background(), "key=", key, "perfix=", clientv3.WithPrefix())
			mlog.Debug("WatchEtcdPicker ", ipPort, key, " error:", info)
			err := this.UpgradeStatus("status", "error", id)
			if err != nil {
				mlog.Debug("WatchEtcdPicker's update status error:", err)
			}
			err = this.UpgradeStatus("details", "watchEtcdPicker error", id)
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
	agentEtcdCmdKey := fmt.Sprintf(`%s/%d`, AgentETCDCmdKey, id)
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
		_, err := EtcdCmd("put", agentEtcdCmdKey, agentPar, AgentETCDCmdIpPort)
		if err != nil {
			mlog.Debug("WhtchPicker send stop cmd to agent etcd error!")
		}
		fmt.Println("WatchPicker finished！")
		return 0
	}
	return 1
}
