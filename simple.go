package QesyDb

import (
	"database/sql"
)

func ModelGet() Model {
	var m Model
	return m
}

func Select(TableName string, Cond map[string]string) ([]map[string]string, error) {
	var m Model
	return m.SetTable(TableName).SetWhere(Cond).ExecSelect()
}

func SelectOne(TableName string, Cond map[string]string) (map[string]string, error) {
	var m Model
	return m.SetTable(TableName).SetWhere(Cond).ExecSelectOne()
}

func Insert(TableName string, InsertArr map[string]string) (sql.Result, error) {
	var m Model
	return m.SetTable(TableName).SetInsert(InsertArr).ExecInsert()
}

func InsertBatch(TableName string, InsertArr []map[string]string) (sql.Result, error) {
	var m Model
	return m.SetTable(TableName).SetInsertArr(InsertArr).ExecInsertBatch()
}

func Replace(TableName string, InsertArr map[string]string) (sql.Result, error) {
	var m Model
	return m.SetTable(TableName).SetInsert(InsertArr).ExecReplace()
}

func Update(TableName string, UpdateArr map[string]string, Cond map[string]string) (sql.Result, error) {
	var m Model
	return m.SetTable(TableName).SetWhere(Cond).SetUpdate(UpdateArr).ExecUpdate()
}

func Delete(TableName string, Cond map[string]string) (sql.Result, error) {
	var m Model
	return m.SetTable(TableName).SetWhere(Cond).ExecDelete()
}
