package QesyDb

import (
	"database/sql"
)

func SelectTx(TableName string, Cond map[string]string, m Model) ([]map[string]string, error) {
	return m.SetTable(TableName).SetWhere(Cond).ExecSelect()
}

func SelectOneTx(TableName string, Cond map[string]string, m Model) (map[string]string, error) {
	return m.SetTable(TableName).SetWhere(Cond).ExecSelectOne()
}

func InsertTx(TableName string, InsertArr map[string]string, m Model) (sql.Result, error) {
	return m.SetTable(TableName).SetInsert(InsertArr).ExecInsert()
}

func InsertBatchTx(TableName string, InsertArr []map[string]string, m Model) (sql.Result, error) {
	return m.SetTable(TableName).SetInsertArr(InsertArr).ExecInsertBatch()
}

func ReplaceTx(TableName string, InsertArr map[string]string, m Model) (sql.Result, error) {
	return m.SetTable(TableName).SetInsert(InsertArr).ExecReplace()
}

func UpdateTx(TableName string, UpdateArr map[string]string, Cond map[string]string, m Model) (sql.Result, error) {
	return m.SetTable(TableName).SetWhere(Cond).SetUpdate(UpdateArr).ExecUpdate()
}

func DeleteTx(TableName string, Cond map[string]string, m Model) (sql.Result, error) {
	return m.SetTable(TableName).SetWhere(Cond).ExecDelete()
}
