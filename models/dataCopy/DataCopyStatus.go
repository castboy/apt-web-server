package dataCopy

import (
	"apt-web-server_v2/models/db"
	"apt-web-server_v2/modules/mlog"
	"fmt"
	"strings"
	"time"
)

var TaskStatusCache DCTSList

func (this *TblDCT) UpdateTaskStatus(para *DCTSPara) (error, *DCTSList) {
	var list DCTSList
	var taskId int
	var taskStatus string
	r := strings.NewReplacer("~", " ")
	/******获取对应任务的id******/
	query := fmt.Sprintf(`SELECT id,status,diskpath FROM %s 
	    WHERE location IN ('%s') AND time IN ('%d');`,
		this.TableName(), para.Location, para.Date)
	rows, err := db.DB.Query(query)
	if err != nil {
		mlog.Debug(query, "Get datacopy task status error in start programe!")
		return err, &list
	}
	defer rows.Close()

	for rows.Next() {
		ugc := new(TblDCT)
		err = rows.Scan(
			&taskId,
			&taskStatus,
			&ugc.DiskPath)
		if err != nil {
			return err, nil
		}
	}
	if para.Rate == 100 {
		TaskStatusCache.Status = "complete"
		err = this.UpgradeDCTStatus("status", "complete", taskId)
		if err != nil {
			mlog.Debug("update status to running error in UpdateTaskStatus")
		}
	}
	if TaskStatusCache.Location == para.Location &&
		TaskStatusCache.Date == para.Date {
		TaskStatusCache.RateOfProgress = para.Rate
	} else {
		TaskStatusCache.Location = para.Location
		TaskStatusCache.Date = para.Date
		TaskStatusCache.RateOfProgress = para.Rate
	}

	switch para.Status {
	case "error":
		TaskStatusCache.Status = "error"
		TaskStatusCache.Details = r.Replace(para.Details)
		err = this.UpgradeDCTStatus("status", "error", taskId)
		err = this.UpgradeDCTStatus("details", TaskStatusCache.Details, taskId)
		if err != nil {
			mlog.Debug("update status error in UpdateTaskStatus")
		}
	default:
		if para.Rate != 100 {
			TaskStatusCache.Status = "running"
			if taskStatus != "running" {
				err = this.UpgradeDCTStatus("status", "running", taskId)
				if err != nil {
					mlog.Debug("update status to running error in UpdateTaskStatus")
				}
			}
		}
	}

	list = TaskStatusCache
	fmt.Println("update task status", "taskId=", taskId, "taskStatus=", taskStatus)
	return nil, &list
}

func (this *TblDCT) GetTaskStatus(para *DCTSPara) (error, *DCTSList) {
	var list DCTSList

	if (para.Location == TaskStatusCache.Location &&
		para.Date == TaskStatusCache.Date) ||
		(para.Location == "" && para.Date == 0) {
		list = TaskStatusCache
	} else {
		list.Location = para.Location
		list.Date = para.Date
		list.Status = "error"
		list.Details = "this task is not running"
	}

	return nil, &list
}

func (this *TblDCT) GetLastTime(para *DCTSPara) (error, *DCTSLastTime) {
	var lastTime DCTSLastTime
	r := strings.NewReplacer("/", "-")
	query := fmt.Sprintf(`SELECT MAX(endday) FROM %s WHERE status IN ('complete');`,
		this.TableName())
	rows, err := db.DB.Query(query)
	if err != nil {
		lastTime.LastTime = "faild"
		mlog.Debug(query, "Get datacopy task enday error in getlasttime programe!")
		return err, &lastTime
	}
	defer rows.Close()

	for rows.Next() {
		ugc := new(TblDCT)
		err = rows.Scan(
			&ugc.DataEnd)
		if err != nil {
			lastTime.LastTime = "00-00-00"
			return nil, &lastTime
		}
		endtime, err := time.Parse("2006-01-02", r.Replace(ugc.DataEnd))
		if err != nil {
			lastTime.LastTime = "00-00-00"
			return err, &lastTime
		}
		lastTime.LastTime = endtime.Add(24 * time.Hour).Format("2006-01-02")
	}
	fmt.Println("get last time")
	return nil, &lastTime
}
