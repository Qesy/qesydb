package qesydb

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" //mysql 包
)

// Db 指针
var Db *sql.DB

// OpenLog 是否记录日志
var OpenLog int = 0

// Model 结构
type Model struct {
	Cond      interface{}
	Insert    map[string]string
	InsertArr []map[string]string
	Update    map[string]string
	Field     string
	Table     string
	Index     string
	Limit     interface{}
	Sort      string
	GroupBy   string
	IsDeug    int
	Tx        *sql.Tx
}

// Connect  is a method with a sql.
func Connect(connStr string) error {
	if sqlDb, err := sql.Open("mysql", connStr); err == nil {
		//defer sqlDb.Close()
		sqlDb.SetMaxOpenConns(600)
		sqlDb.SetMaxIdleConns(600)
		sqlDb.SetConnMaxLifetime(time.Minute * 3)
		//sqlDb.SetConnMaxIdleTime(4 * time.Second)

		if err = sqlDb.Ping(); err != nil {
			return err
		}
		Db = sqlDb
	} else {
		return err
	}
	return nil
}

// Begin 开始事务
func Begin() (*sql.Tx, error) {
	return Db.Begin()
}

// Rollback 事务回滚
func Rollback(tx *sql.Tx) error {
	return tx.Rollback()
}

// Commit 事务提交
func Commit(tx *sql.Tx) error {
	return tx.Commit()
}

// ExecSelectIndex 返回一个MAP
func (m *Model) ExecSelectIndex() (map[string]map[string]string, error) {
	resultsSlice, err := m.execSelect()
	retArr := map[string]map[string]string{}
	for _, v := range resultsSlice {
		if v[m.Index] == "" {
			continue
		}
		retArr[v[m.Index]] = v
	}
	m.Clean()
	return retArr, err
}

// Query 查询SQL,返回一个 切片MAP;
// SqlStr : SQL语句
func (m *Model) Query(sqlStr string) ([]map[string]string, error) {
	ret, err := m.query(sqlStr)
	m.Clean()
	return ret, err
}

// ExecSelect 执行查询 返回一个 切片MAP
func (m *Model) ExecSelect() ([]map[string]string, error) {
	ret, err := m.execSelect()
	m.Clean()
	return ret, err
}

// ExecSelect 拼装SQL语句
func (m *Model) execSelect() ([]map[string]string, error) {
	cond := m.getSQLCond()
	field := m.getSQLField()
	sort := m.getSort()
	groupby := m.getGroupBy()
	limit := m.getSQLLimite()
	sqlStr := "SELECT " + field + " FROM " + m.Table + cond + groupby + sort + limit + ";"
	return m.query(sqlStr)
}

// ExecSelectOne 只查询一条
func (m *Model) ExecSelectOne() (map[string]string, error) {
	m.SetLimit([2]int{0, 1})
	resultsSlice, err := m.ExecSelect()
	if len(resultsSlice) == 0 {
		return map[string]string{}, err
	}
	return resultsSlice[0], nil
}

// ExecUpdate 修改
func (m *Model) ExecUpdate() (sql.Result, error) {
	updateStr := m.getSQLUpdate()
	condStr := m.getSQLCond()
	sqlStr := "UPDATE " + m.Table + " SET " + updateStr + condStr + ";"
	m.Debug(sqlStr)
	var err error
	var stmt *sql.Stmt
	if m.Tx == nil {
		stmt, err = Db.Prepare(sqlStr)

	} else {
		stmt, err = m.Tx.Prepare(sqlStr)
	}
	if err != nil {
		logRecord("ERR:" + err.Error() + "SQL:" + sqlStr)
		return nil, err
	}
	defer stmt.Close()
	result, err := stmt.Exec()
	m.Clean()
	return result, err
}

// ExecInsert 添加
func (m *Model) ExecInsert() (sql.Result, error) {
	insert := m.getSQLInsert()
	sqlStr := "INSERT INTO " + m.Table + " " + insert + ";"
	m.Debug(sqlStr)
	var err error
	var stmt *sql.Stmt
	if m.Tx == nil {
		stmt, err = Db.Prepare(sqlStr)
	} else {
		stmt, err = m.Tx.Prepare(sqlStr)
	}
	if err != nil {
		logRecord("ERR:" + err.Error() + "SQL:" + sqlStr)
		return nil, err
	}
	defer stmt.Close()
	result, err := stmt.Exec()
	m.Clean()
	return result, err
}

// ExecInsertBatch 批量添加 （预计要删除）
func (m *Model) ExecInsertBatch() (sql.Result, error) {
	insert := m.getSQLInsertArr()
	sqlStr := "INSERT INTO " + m.Table + " " + insert + ";"
	m.Debug(sqlStr)
	var err error
	var stmt *sql.Stmt
	if m.Tx == nil {
		stmt, err = Db.Prepare(sqlStr)
	} else {
		stmt, err = m.Tx.Prepare(sqlStr)
	}
	if err != nil {
		logRecord("ERR:" + err.Error() + "SQL:" + sqlStr)
		return nil, err
	}
	defer stmt.Close()
	result, err := stmt.Exec()
	m.Clean()
	return result, err
}

// ExecReplace 替换
func (m *Model) ExecReplace() (sql.Result, error) {
	insert := m.getSQLInsert()
	sqlStr := "REPLACE INTO " + m.Table + " " + insert + ";"
	m.Debug(sqlStr)
	var err error
	var stmt *sql.Stmt
	if m.Tx == nil {
		stmt, err = Db.Prepare(sqlStr)
	} else {
		stmt, err = m.Tx.Prepare(sqlStr)
	}
	if err != nil {
		logRecord("ERR:" + err.Error() + "SQL:" + sqlStr)
		return nil, err
	}
	defer stmt.Close()
	result, err := stmt.Exec()
	m.Clean()
	return result, err
}

// ExecDelete 删除
func (m *Model) ExecDelete() (sql.Result, error) {
	condStr := m.getSQLCond()
	sqlStr := "DELETE FROM " + m.Table + condStr + ";"
	m.Debug(sqlStr)
	var err error
	var stmt *sql.Stmt
	if m.Tx == nil {
		stmt, err = Db.Prepare(sqlStr)
	} else {
		stmt, err = m.Tx.Prepare(sqlStr)
	}
	if err != nil {
		logRecord("ERR:" + err.Error() + "SQL:" + sqlStr)
		return nil, err
	}
	defer stmt.Close()
	result, err := stmt.Exec()
	m.Clean()
	return result, err
}

// Exec 执行SQL语句
func (m *Model) Exec(sqlStr string) (sql.Result, error) {
	var err error
	var stmt *sql.Stmt
	if m.Tx == nil {
		stmt, err = Db.Prepare(sqlStr)
	} else {
		stmt, err = m.Tx.Prepare(sqlStr)
	}
	if err != nil {
		logRecord("ERR:" + err.Error() + "SQL:" + sqlStr)
		return nil, err
	}
	defer stmt.Close()
	result, err := stmt.Exec()
	m.Clean()
	return result, err
}

// GetLastInsertID 获取最后插入的ID
func GetLastInsertID(result sql.Result) (int64, error) {
	return result.LastInsertId()
}

// GetRowsAffected 获取受影响行数
func GetRowsAffected(result sql.Result) (int64, error) {
	return result.RowsAffected()
}

func (m *Model) getSQLCond() string {
	if str, ok := m.Cond.(string); ok || m.Cond == nil {
		return str
	}
	var strArr []string
	if arr, ok := m.Cond.(map[string]string); ok {
		if len(arr) == 0 {
			return ""
		}
		var strArr []string
		for k, v := range arr {
			if strings.Contains(k, "LIKE") {
				strArr = append(strArr, k+" '%"+v+"%'")
			} else if strings.Contains(k, ">") || strings.Contains(k, "<") {
				strArr = append(strArr, k+" "+v)
			} else {
				strArr = append(strArr, k+"='"+v+"'")
			}

		}
		if len(strArr) == 0 {
			return ""
		}
		return " WHERE " + strings.Join(strArr, " && ")
	}
	if arr, ok := m.Cond.(map[string]interface{}); ok {
		if len(arr) == 0 {
			return ""
		}
		for k, v := range arr {
			if isStr, ok := v.(string); ok {
				if strings.Contains(k, "LIKE") {
					strArr = append(strArr, k+" '%"+isStr+"%'")
				} else if strings.Contains(k, ">") || strings.Contains(k, "<") {
					strArr = append(strArr, k+" "+isStr)
				} else {
					strArr = append(strArr, k+"='"+isStr+"'")
				}

			} else if isStrArrTmp, ok := v.([]string); ok {
				if len(isStrArrTmp) == 0 {
					strArr = append(strArr, k+"=''")
				} else {
					strArr = append(strArr, k+" in ('"+strings.Join(isStrArrTmp, "', '")+"')")
				}

			} else {
				return "WHERE "
			}
		}
		return " WHERE " + strings.Join(strArr, " && ")
	}
	return " WHERE "
}

func (m *Model) getSQLField() string {
	if m.Field != "" {
		return m.Field
	}
	return "*"
}

func (m *Model) getSort() string {
	if m.Sort != "" {
		return " ORDER BY " + m.Sort + " "
	}
	return ""
}

func (m *Model) getGroupBy() string {
	if m.GroupBy != "" {
		return " GROUP BY " + m.GroupBy + " "
	}
	return ""
}

func (m *Model) getSQLUpdate() string {
	var strArr []string
	for k, v := range m.Update {
		strArr = append(strArr, k+"='"+v+"'")
	}
	return strings.Join(strArr, ",")
}

func (m *Model) getSQLInsert() string {
	var fieldArr, valueArr []string
	for k, v := range m.Insert {
		fieldArr = append(fieldArr, k)
		valueArr = append(valueArr, "'"+v+"'")
	}
	return "(" + strings.Join(fieldArr, ",") + ") values (" + strings.Join(valueArr, ",") + ")"
}

func (m *Model) getSQLInsertArr() string {
	fieldArr, valuesArr := []string{}, []string{}
	for k := range m.InsertArr[0] {
		fieldArr = append(fieldArr, k)
	}
	for _, value := range m.InsertArr {
		var valueArr []string
		for _, v := range fieldArr {
			valueArr = append(valueArr, "'"+value[v]+"'")
		}
		valuesArr = append(valuesArr, "("+strings.Join(valueArr, ",")+")")
	}
	return "(" + strings.Join(fieldArr, ",") + ")" + " values " + strings.Join(valuesArr, ",")
}

func (m *Model) getSQLLimite() string {
	if strArr, ok := m.Limit.([2]int); ok {
		return " LIMIT " + fmt.Sprintf("%d", strArr[0]) + ", " + fmt.Sprintf("%d", strArr[1])
	}
	if strArr, ok := m.Limit.([]int); ok {
		return " LIMIT " + fmt.Sprintf("%d", strArr[0]) + ", " + fmt.Sprintf("%d", strArr[1])
	}
	return ""
}

// Clean 清楚orm
func (m *Model) Clean() {
	m.Cond = nil
	m.Insert = nil
	m.InsertArr = nil
	m.Update = nil
	m.Field = ""
	m.Table = ""
	m.Index = ""
	m.Limit = nil
	m.Sort = ""
	m.GroupBy = ""
	m.IsDeug = 0
}

// Debug 打印调试
func (m *Model) Debug(sql string) {
	if m.IsDeug == 1 {
		fmt.Println(sql)
	}
}

func logRecord(str string) {
	if OpenLog == 0 {
		return
	}
	log.Println(str)
}

func (m *Model) query(sqlStr string) ([]map[string]string, error) {
	m.Debug(sqlStr)
	var err error
	var stmt *sql.Stmt
	resultsSlice := []map[string]string{}
	if m.Tx == nil {
		stmt, err = Db.Prepare(sqlStr)
		if err != nil {
			logRecord("ERR:" + err.Error() + "SQL:" + sqlStr)
			return resultsSlice, err
		}
		defer stmt.Close()
	} else {
		stmt, err = m.Tx.Prepare(sqlStr)
		if err != nil {
			logRecord("ERR:" + err.Error() + "SQL:" + sqlStr)
			return resultsSlice, err
		}
		defer stmt.Close()
	}

	rows, err := stmt.Query()
	if err != nil {
		logRecord("ERR:" + err.Error() + "SQL:" + sqlStr)
		return resultsSlice, err
	}
	defer rows.Close()
	fields, err := rows.Columns()
	if err != nil {
		logRecord("ERR:" + err.Error() + "SQL:" + sqlStr)
		return resultsSlice, err
	}

	for rows.Next() {
		result := make(map[string]string)
		var scanResultContainers []interface{}
		for i := 0; i < len(fields); i++ {
			var scanResultContainer interface{}
			scanResultContainers = append(scanResultContainers, &scanResultContainer)
		}
		if err := rows.Scan(scanResultContainers...); err != nil {
			return resultsSlice, err
		}
		for k, v := range fields {
			rawValue := reflect.Indirect(reflect.ValueOf(scanResultContainers[k]))
			if rawValue.Interface() == nil {
				continue
			}
			rawType := reflect.TypeOf(rawValue.Interface())
			rawVal := reflect.ValueOf(rawValue.Interface())
			var str string
			switch rawType.Kind() {
			case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				str = strconv.FormatInt(rawVal.Int(), 10)
				result[v] = str
			case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				str = strconv.FormatUint(rawVal.Uint(), 10)
				result[v] = str
			case reflect.Float32, reflect.Float64:
				str = strconv.FormatFloat(rawVal.Float(), 'f', -1, 64)
				result[v] = str
			case reflect.Slice:
				if rawType.Elem().Kind() == reflect.Uint8 {
					result[v] = string(rawVal.Interface().([]byte))
				}
			case reflect.String:
				str = rawVal.String()
				result[v] = str
			case reflect.Struct:
				str = rawVal.Interface().(time.Time).Format("2006-01-02 15:04:05.000 -0700")
				result[v] = str
			case reflect.Bool:
				if rawVal.Bool() {
					result[v] = "1"
				} else {
					result[v] = "0"
				}
			}
		}
		resultsSlice = append(resultsSlice, result)
	}
	return resultsSlice, nil
}
