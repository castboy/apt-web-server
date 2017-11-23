package offlineAssignment

import (
	"apt-web-server_v2/models/db"
	"apt-web-server_v2/modules/mlog"
	"encoding/json"
	"fmt"
)

func (this *TblOLA) GetTaskStatus(para *TaskListPara) (error, *[]TaskStatusList) {
	var s []TaskStatusRequestPara
	var list []TaskStatusList
	json.Unmarshal([]byte(para.TaskList), &s)
	for out := range s {
		statusElement := new(TaskStatusList)
		_, statusElement = this.GetStatus(s[out].Name, int64(s[out].Time), para.OfflineTag)

		list = append(list, TaskStatusList{statusElement.TaskStatusContent})
	}
	return nil, &list
}

func (this *TblOLA) GetStatus(paraName string, paraTime int64, tag string) (error, *TaskStatusList) {
	var taskID int
	var taskType, topicName string
	var statusReq TaskStatusList
	var statusCount TaskEtcdResData
	query := fmt.Sprintf(`SELECT id,type,topic,status FROM %s 
	WHERE name IN ('%s') AND time IN ('%d');`,
		this.TableName(tag),
		paraName,
		paraTime)
	rows, err := db.DB.Query(string(query))
	if err != nil {
		mlog.Debug(query, err)
		return err, nil
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(
			&taskID,
			&taskType,
			&topicName,
			&statusReq.Status)
		if err != nil {
			return err, nil
		}
	}
	statusReq.Name = paraName
	statusReq.Time = paraTime
	if statusReq.Status == "ready" || statusReq.Status == "wait" {
		return nil, &statusReq
	}
	pickerKey := fmt.Sprintf(`%s/%d`, PickerStatusKey, taskID)

	//pickerCount, pickerTotal, err := GetEtcdPicker(pickerKey, EtcdIpPort)
	pickerStat, err := GetEtcdPicker(pickerKey)
	if err != nil {
		mlog.Debug("OfflineTaskStatus GetEtcdPicker error:", err)
	}
	statusCount.PickerCount = float32(pickerStat.Count)
	statusCount.PickerTotal = float32(pickerStat.Total)

	for index, agentKey := range AgentStatKey {
		fmt.Println(index)
		//count, total, err := GetEtcdAgent(taskType, topicName, agentKey)
		aStatus, err := GetEtcdAgent(taskType, topicName, agentKey)
		if err != nil {
			mlog.Debug("OfflineTaskStatus GetEtcdAgent error:", err)
			statusReq.Status = "error"
			return err, &statusReq
		}
		statusCount.AgentCount += float32(aStatus.Engine + aStatus.Err)
		statusCount.AgentTotal += float32(aStatus.Last)
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
		err := this.UpgradeStatus("status", "complete", taskID, tag)
		if err != nil {
			mlog.Debug("set task ", taskID, " status to complete error!")
			return nil, &statusReq
		}
		statusReq.Status = "complete"
	}
	return nil, &statusReq
}
