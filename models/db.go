package models

import (
	"apt-web-server/modules/mlog"
	//"errors"
	"strings"
	"unicode"
)

type CreatTabler interface {
	CreateSql() string
	TableName() string
}

func CreateTables(models ...CreatTabler) {
	for _, model := range models {
		do_create_table(model.TableName(), model.CreateSql())
	}
}

func MysqlErrorNum(ret string) string {
	//return strings.Fields(ret)[1]
	return strings.FieldsFunc(ret, func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	})[1]
}

func do_create_table(tableName string, sql string) error {
	_, err := db.Exec(sql)

	if err != nil && MysqlErrorNum(err.Error()) == "1050" {
		mlog.Notice("Table \"", tableName, "\" already exists")
		return err
	} else if err != nil {
		mlog.Debug("err:", err)
		return err
	} else {
		mlog.Notice("Create table", tableName)
	}
	return nil
}

