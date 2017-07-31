package models

import (
	"apt-web-server/modules/mlog"
	"encoding/json"
	"fmt"
)

func (this *TblOLA) GetTaskStatus(para *TaskListPara) (error, *[]TaskStatusList) {
	var s []TaskStatusRequestPara
	var list []TaskStatusList
	json.Unmarshal([]byte(para.TaskList), &s)
	fmt.Println(para.TaskList)
	for out := range s {
		statusElement := new(TaskStatusList)
		_, statusElement = this.GetStatus(s[out].Name, int64(s[out].Time), para.OfflineTag)

		list = append(list, TaskStatusList{statusElement.TaskStatusContent})
	}
	return nil, &list
}

func (this *TblOLA) GetStatus(paraName string, paraTime int64, paraOfflineTag string) (error, *TaskStatusList) {
	var taskID int
	var taskType, topicName string
	var statusReq TaskStatusList
	var statusCount TaskEtcdResData
	query := fmt.Sprintf(`select id,type,topic,status from %s where name='%s' and time=%d;`,
		this.TableName(paraOfflineTag),
		paraName,
		paraTime)
	rows, err := db.Query(string(query))
	fmt.Println("query++++++++++++++++++++++++++++++++")
	fmt.Println(query)
	if err != nil {
		//return err, nil
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(
			&taskID,
			//&statusReq.Name,
			//&statusReq.Time,
			&taskType,
			&topicName,
			&statusReq.Status)
		if err != nil {
			//return err, nil
		}
	}
	statusReq.Name = paraName
	statusReq.Time = paraTime
	if statusReq.Status == "ready" || statusReq.Status == "wait" {
		return nil, &statusReq
	}
	pickerKey := fmt.Sprintf(`picker/%d`, taskID)

	pickerCount, pickerTotal, err := GetEtcdPicker(pickerKey, PickerETCDIpPort)
	if err != nil {
		mlog.Debug("OfflineTaskStatus GetEtcdPicker error:", err)
	}
	statusCount.PickerCount = float32(pickerCount)
	statusCount.PickerTotal = float32(pickerTotal)

	for index, agentKey := range AgentStatusETCDSlice {
		count, total, err := GetEtcdAgent(taskType, topicName, agentKey, AgentStatusETCDIpPort[index])
		if err != nil {
			mlog.Debug("OfflineTaskStatus' GetEtcdAgent error:", err)
			statusReq.Status = "error"
			return err, &statusReq
		}
		statusCount.AgentCount += float32(count)
		statusCount.AgentTotal += float32(total)
	}

	if statusCount.PickerTotal == 0 {
		statusReq.PickerPercent = 0
	} else {
		statusReq.PickerPercent = float32(int((statusCount.PickerCount/statusCount.PickerTotal)*10000)) / 100
	}

	if statusReq.PickerPercent == 100 && statusCount.AgentTotal == 0 {
		statusReq.AgentPercent = 100
	} else if statusCount.AgentTotal == -1 || statusCount.AgentTotal == 0 {
		statusReq.AgentPercent = 0
	} else {
		statusReq.AgentPercent = float32(int((statusCount.AgentCount/statusCount.AgentTotal)*10000)) / 100
	}

	if statusReq.PickerPercent == 100 && statusReq.AgentPercent == 100 && statusReq.Status != "complete" {
		err := this.UpgradeStatus("status", "complete", taskID, paraOfflineTag)
		if err != nil {
			fmt.Println("set task ", taskID, " status error!")
			return nil, &statusReq
		}
		statusReq.Status = "complete"
	}
	return nil, &statusReq
}
