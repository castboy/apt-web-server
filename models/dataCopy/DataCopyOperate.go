/********数据迁移********/
package dataCopy

import (
	"apt-web-server_v2/models/db"
	"apt-web-server_v2/models/modelsPublic"
	"apt-web-server_v2/modules/mlog"
	"fmt"
	"time"
)

/************default config************/

func (this *TblDCT) TableName() string {
	return "datacopy"
}

func (this *TblDCT) CreateDCT(para *DCTOperatePara) (error, *DCTOperateResult) {
	var res_t DCTOperateResult
	para.Date = time.Now().Unix()
	query := fmt.Sprintf(`INSERT INTO %s(location,startday,endday,time,diskpath,status) 
						VALUE('%s','%s','%s',%d,'%s','%s');`,
		this.TableName(), para.Location, para.Start, para.End, para.Date,
		para.DiskPath, "ready")
	rows, err := db.DB.Query(query)
	fmt.Println(query)
	if err != nil {
		mlog.Debug(query, "Creat task ", para.Location, para.Date, " error")
		res_t.Result = "faild"
		return err, &res_t
	}
	defer rows.Close()
	TaskStatusCache.Location = para.Location
	TaskStatusCache.Date = para.Date
	TaskStatusCache.Status = "ready"

	res_t.Result = "ok"
	return nil, &res_t
}

func (this *TblDCT) DeleteDCT(para *DCTOperatePara) (error, *DCTOperateResult) {
	var res_t DCTOperateResult
	query := fmt.Sprintf(`DELETE FROM %s WHERE location IN ('%s') AND time IN ('%d');`,
		this.TableName(), para.Location, para.Date)
	rows, err := db.DB.Query(query)
	fmt.Println(query)
	if err != nil {
		mlog.Debug(query, "Delete task error")
		res_t.Result = "faild"
		return err, &res_t
	}
	defer rows.Close()

	res_t.Result = "ok"
	return nil, &res_t
}

func (this *TblDCT) StartDCT(para *DCTOperatePara) (error, *DCTOperateResult) {
	var res_t DCTOperateResult
	var taskId int
	query := fmt.Sprintf(`SELECT id,status,diskpath,startday,endday 
	    FROM %s WHERE location IN ('%s') AND time IN ('%d');`,
		this.TableName(), para.Location, para.Date)
	rows, err := db.DB.Query(query)
	if err != nil {
		res_t.Result = "faild"
		mlog.Debug(query, "Get datacopy task status error in start programe!")
		return err, &res_t
	}
	defer rows.Close()

	for rows.Next() {
		ugc := new(TblDCT)
		err = rows.Scan(
			&ugc.Id,
			&ugc.Status,
			&ugc.DiskPath,
			&ugc.DataStart,
			&ugc.DataEnd)
		if err != nil {
			return err, nil
		}
		if ugc.Status == "running" {
			res_t.Result = "this task is already running"
			return err, &res_t
		}
		taskId = ugc.Id
		if taskId == 0 {
			res_t.Result = "no such task"
			return nil, &res_t
		} else {
			res_t.Result = "ok"
		}
		err = this.UpgradeDCTStatus("status", "running", taskId)
		err = this.UpgradeDCTStatus("details", "", taskId)
		if err != nil {
			mlog.Debug("task set running error")
		}
		TaskStatusCache.Status = "running"
		TaskStatusCache.Details = ""
		go func() {
			fmt.Println("start copy data")
			shellpath := fmt.Sprintf("/etc/diskcopy/diskTransfer.sh add %d %s %s", para.Date, ugc.DataStart, ugc.DataEnd)
			err, _ = modelsPublic.DoShell(shellpath)
			if err != nil {
				res_t.Result = "run shell error"
				err = this.UpgradeDCTStatus("status", "error", taskId)
				err = this.UpgradeDCTStatus("details", "script running error", taskId)
				if err != nil {
					mlog.Debug("task set running shell error")
				}
				TaskStatusCache.Status = "error"
				TaskStatusCache.Details = "script running error"
			}
		}()
	}

	return nil, &res_t
}
func (this *TblDCT) RemoveDiskDCT(para *DCTOperatePara) (error, *DCTOperateResult) {
	var res_t DCTOperateResult
	var taskId int
	query := fmt.Sprintf(`SELECT id,status FROM %s WHERE location IN ('%s') AND 
	    diskpath IN ('%s') AND status IN ('ready','running');`,
		this.TableName(), para.Location, para.DiskPath)
	rows, err := db.DB.Query(query)
	if err != nil {
		res_t.Result = "faild"
		mlog.Debug(query, "Get datacopy task status error in start programe!")
		return err, &res_t
	}
	defer rows.Close()

	for rows.Next() {
		ugc := new(TblDCT)
		err = rows.Scan(
			&ugc.Id,
			&ugc.Status)
		if err != nil {
			return err, nil
		}
		taskId = ugc.Id
		switch ugc.Status {
		case "running":
			_ = this.UpgradeDCTStatus("status", "error", taskId)
			_ = this.UpgradeDCTStatus("details", "disk unmounted in task running", taskId)
		case "ready":
			_ = this.UpgradeDCTStatus("status", "error", taskId)
			_ = this.UpgradeDCTStatus("details", "disk unmounted before task running", taskId)
		}
	}
	TaskStatusCache.Location = ""
	TaskStatusCache.Date = 0
	TaskStatusCache.Status = ""
	TaskStatusCache.RateOfProgress = 0
	TaskStatusCache.Details = ""

	res_t.Result = "ok"
	return nil, &res_t
}
func (this *TblDCT) GetDiskDCT(para *DCTOperatePara) (error, *DCTOperateResult) {
	var res_t DCTOperateResult
	var listcount int
	query := fmt.Sprintf(`SELECT COUNT(id) FROM %s WHERE status IN ('ready','running');`,
		this.TableName())
	rows, err := db.DB.Query(query)
	if err != nil {
		res_t.Result = "faild"
		mlog.Debug(query, "Get datacopy task status error in getdisk programe!")
		return err, &res_t
	}
	defer rows.Close()
	for rows.Next() {
		fmt.Println("in for loop")
		err = rows.Scan(
			&listcount)
		if err != nil {
			res_t.Result = "no"
			return nil, &res_t
		}
		if listcount == 0 {
			res_t.Result = "no"
			return err, &res_t
		}
	}
	res_t.Result = "yes"
	return nil, &res_t
}

func (this *TblDCT) UpgradeDCTStatus(column, value string, taskID int) error {
	query := fmt.Sprintf(`UPDATE %s SET %s='%s' WHERE id=%d;`,
		this.TableName(),
		column,
		value,
		taskID)
	rows, err := db.DB.Query(query)
	fmt.Println(query)
	if err != nil {
		mlog.Debug("error!upgrade status faild!")
		return err
	}
	defer rows.Close()
	return nil
}
