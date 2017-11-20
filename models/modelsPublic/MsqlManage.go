package modelsPublic

import (
	"apt-web-server_v2/models/db"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func Insert_to_mysql(insert_cmd string) bool {
	/*
		db := Init_mysql(Sqlcfg)
		defer db.Close()
	*/
	//result, _ := db.Exec("insert into user values(?,?,?)", "test", 2, "test")
	//c, _ := result.RowsAffected()
	//fmt.Println("add affected rows:", c)
	_, err := db.DB.Exec(insert_cmd)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}
func Select_mysql(select_cmd string) *sql.Rows {
	var rows *sql.Rows
	rows, err1 := db.DB.Query(select_cmd)
	if err1 != nil {
		fmt.Println(err1.Error())
		return rows
	}
	//fmt.Println(rows)
	//fmt.Println("rows type:", reflect.TypeOf(db))
	return rows
}
func Update_mysql(update_cmd string) bool {
	_, err := db.DB.Exec(update_cmd)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}
func Delete_mysql(delete_cmd string) bool {
	_, err := db.DB.Exec(delete_cmd)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}
