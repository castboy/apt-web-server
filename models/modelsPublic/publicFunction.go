package modelsPublic

import (
	"apt-web-server_v2/models/db"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os/exec"
)

func GetDateMold(unit string) string {
	var datemold string
	switch unit {
	case "second":
		datemold = "%Y-%m-%d %H:%i:%s"
	case "minute":
		datemold = "%Y-%m-%d %H:%i"
	case "hour":
		datemold = "%Y-%m-%d %H"
	case "day":
		datemold = "%Y-%m-%d"
	case "week":
		datemold = "%Y %u"
	case "month":
		datemold = "%Y-%m"
	default:
		datemold = "%Y-%m-%d"
	}
	return datemold
}

func DefaultParaCmd(cmdType string, tablename string, para *TblPublicPara) ([]string, int) {
	qslice := make([]string, 0)
	flag := 0
	var qslice_tmp string

	switch cmdType {
	case "getcounts":
		if tablename != "" {
			qslice_tmp = fmt.Sprintf(`SELECT count(id) FROM %s `, tablename)

		} else {
			qslice_tmp = fmt.Sprintf(`SELECT FOUND_ROWS() AS count`)
			qslice = append(qslice, qslice_tmp)
			return qslice, 0
		}
	case "getlist":
		qslice_tmp = fmt.Sprintf(`SELECT * FROM %s `, tablename)
		flag = 0
	}
	qslice = append(qslice, qslice_tmp)
	if para.Start != 0 && para.End != 0 {
		var temp_se string
		temp_se = fmt.Sprintf(" WHERE time BETWEEN %d AND %d", para.Start, para.End)
		flag = 1
		qslice = append(qslice, temp_se)
	}
	return qslice, flag
}
func GetOfflineTaskID(tblname, taskname string, time int64) (error, int) {
	var taskID int
	query := fmt.Sprintf(`SELECT id FROM %s WHERE name='%s' AND time=%d;`,
		tblname,
		taskname,
		time)
	fmt.Println(query)
	rows, err := db.DB.Query(query)
	if err != nil {
		return err, 0
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(
			&taskID)
		if err != nil {
			return err, 0
		}
	}
	return nil, taskID
}
func GetTimes(tblName string, taskId int) (int64, int64) {
	var maxtime, mintime int64
	query := fmt.Sprintf(`SELECT MAX(time),MIN(time) FROM %s WHERE taskid=%d;`,
		tblName,
		taskId)
	rows, err := db.DB.Query(query)
	if err != nil {
		return 0, 0
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(
			&maxtime,
			&mintime)
		if err != nil {
			return 0, 0
		}
	}
	return maxtime, mintime
}
func DoShell(shellpath string) (err error, output string) {
	cmd := exec.Command("/bin/bash", "-c", shellpath)
	var out bytes.Buffer
	defer func() {
		if info := recover(); info != nil {
			fmt.Println("cmd running error", info)
			err = errors.New("cmd running error")
		}
	}()
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		panic(err)
		log.Fatal(err)
	}
	output = out.String()
	fmt.Printf("%s", output)
	return err, output
}

/*
func GetOflTaskMsg(tblname, taskname string, time int64) (int, int64, int64) {
	var taskID int
	var startTime string
	var endTime string
	query := fmt.Sprintf(`select id,start,end from %s where name='%s' and time=%d;`,
		tblname,
		taskname,
		time)
	fmt.Println(query)
	rows, err := db.DB.Query(query)
	if err != nil {
		//return err, nil
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(
			&taskID,
			&startTime,
			&endTime)
		if err != nil {
			return 0, 0, 0
		}
	}
	return taskID, startTime, endTime
}
*/
