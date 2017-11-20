package db

import (
	"apt-web-server_v2/modules/mconfig"
	"apt-web-server_v2/modules/mlog"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func init() {
	defer func() {
		if r := recover(); r != nil {
			mlog.Error(r)
			panic(r)
		}
	}()
	//init mysql
	dbinit()
	//CreateTables(new(TblUgcD))
}

func dbinit() error {
	var err error

	conf_string := func(section string, option string) (string, error) {
		if err != nil {
			return "", err
		}
		return mconfig.Conf.String(section, option)
	}
	dbPort, err := conf_string("db", "DbPort")
	dbHost, err := conf_string("db", "DbHost")
	dbName, err := conf_string("db", "DbName")
	dbUser, err := conf_string("db", "DbUser")
	dbPassword, err := conf_string("db", "DbPassword")
	if err != nil {
		panic(err)
	}

	dbUrl := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8&&loc=Asia%2FShanghai"
	//dbUrl := "autelan:Autelan1202@tcp(rdsrenv7vrenv7v.mysql.rds.aliyuncs.com:3306)/umsdb?charset=utf8&&loc=Asia%2FShanghai"
	mlog.Debug("dburl=", dbUrl)
	DB, err = sql.Open("mysql", dbUrl)
	if err != nil {
		panic(err)
	}
	err = DB.Ping()
	if err != nil {
		panic(err)
	}
	return nil
}
