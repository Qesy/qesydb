package QesyDb

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var Db *sql.DB

type Model struct {
	Cond   interface{}
	Insert map[string]string
	Update map[string]string
	Field  string
	Table  string
	Index  string
	Limit  interface{}
	Sort   string
	IsDeug int
	Tx     *sql.Tx
}

// Connect  is a method with a sql.
func Connect(connStr string) {
	sqlDb, err := sql.Open("mysql", connStr)
	//defer sqlDb.Close()
	sqlDb.SetConnMaxLifetime(1800)
	sqlDb.SetMaxIdleConns(0)
	sqlDb.SetMaxOpenConns(600)
	if err != nil {
		log.Fatal("mysql connect error")
	}
	err = sqlDb.Ping()
	if err != nil {
		log.Fatal("mysql ping error")
	}
	fmt.Println("mysql connect sueccss")
	Db = sqlDb
}

func Begin() (*sql.Tx, error) {
	return Db.Begin()
}

func Rollback(tx *sql.Tx) error {
	return tx.Rollback()
}

func Commit(tx *sql.Tx) error {
	return tx.Commit()
}

// ExecSelectIndex  is a method with a sql.
func (m *Model) ExecSelectIndex() (map[string]map[string]string, error) {
	resultsSlice, _ := m.execSelect()
	retArr := map[string]map[string]string{}
	for _, v := range resultsSlice {
		if v[m.Index] == "" {
			continue
		}
		retArr[v[m.Index]] = v
	}
	m.Clean()
	return retArr, nil
}

func (m *Model) ExecSelect() ([]map[string]string, error) {
	ret, err := m.execSelect()
	m.Clean()
	return ret, err
}

// ExecSelect is a method with a sql.
func (m *Model) execSelect() ([]map[string]string, error) {
	cond := m.getSQLCond()
	field := m.getSQLField()
	sort := m.getSort()
	limit := m.getSQLLimite()
	sqlStr := "SELECT " + field + " FROM " + m.Table + cond + sort + limit + ";"
	m.Debug(sqlStr)
	var err error
	var stmt *sql.Stmt
	resultsSlice := []map[string]string{}
	if m.Tx == nil {
		stmt, err = Db.Prepare(sqlStr)
		if err != nil {
			return resultsSlice, err
		}
		defer stmt.Close()
	} else {
		stmt, err = m.Tx.Prepare(sqlStr)
	}
	//defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		fmt.Println("DBERR:", rows, err, sqlStr)
		return resultsSlice, err
	}
	defer rows.Close()
	fields, err := rows.Columns()
	if err != nil {
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
					break
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

func (m *Model) ExecSelectOne() (map[string]string, error) {
	resultsSlice, err := m.ExecSelect()
	if len(resultsSlice) == 0 {
		return map[string]string{}, err
	}
	return resultsSlice[0], nil
}

func (m *Model) ExecUpdate() (sql.Result, error) {
	updateStr := m.getSQLUpdate()
	condStr := m.getSQLCond()
	sqlStr := "UPDATE " + m.Table + " SET " + updateStr + condStr + ";"
	m.Debug(sqlStr)
	var err error
	var stmt *sql.Stmt
	if m.Tx == nil {
		stmt, err = Db.Prepare(sqlStr)
		defer stmt.Close()
	} else {
		stmt, err = m.Tx.Prepare(sqlStr)
	}
	if err != nil {
		fmt.Println("DBERR:", stmt, err, sqlStr)
		return nil, err
	}
	result, err := stmt.Exec()
	m.Clean()
	return result, err
}

func (m *Model) ExecInsert() (sql.Result, error) {
	insert := m.getSQLInsert()
	sqlStr := "INSERT INTO " + m.Table + " " + insert + ";"
	m.Debug(sqlStr)
	var err error
	var stmt *sql.Stmt
	if m.Tx == nil {
		stmt, err = Db.Prepare(sqlStr)
		defer stmt.Close()
	} else {
		stmt, err = m.Tx.Prepare(sqlStr)
	}
	if err != nil {
		fmt.Println("DBERR:", stmt, err, sqlStr)
		return nil, err
	}
	result, err := stmt.Exec()
	m.Clean()
	return result, err
}

func (m *Model) ExecReplace() (sql.Result, error) {
	insert := m.getSQLInsert()
	sqlStr := "REPLACE INTO " + m.Table + " " + insert + ";"
	m.Debug(sqlStr)
	var err error
	var stmt *sql.Stmt
	if m.Tx == nil {
		stmt, err = Db.Prepare(sqlStr)
		defer stmt.Close()
	} else {
		stmt, err = m.Tx.Prepare(sqlStr)
	}
	if err != nil {
		fmt.Println("DBERR:", stmt, err, sqlStr)
		return nil, err
	}
	result, err := stmt.Exec()
	m.Clean()
	return result, err
}

func (m *Model) ExecDelete() (sql.Result, error) {
	condStr := m.getSQLCond()
	sqlStr := "DELETE FROM " + m.Table + condStr + ";"
	m.Debug(sqlStr)
	var err error
	var stmt *sql.Stmt
	if m.Tx == nil {
		stmt, err = Db.Prepare(sqlStr)
		defer stmt.Close()
	} else {
		stmt, err = m.Tx.Prepare(sqlStr)
	}
	if err != nil {
		fmt.Println("DBERR:", stmt, err, sqlStr)
		return nil, err
	}
	result, err := stmt.Exec()
	m.Clean()
	return result, err
}

func (m *Model) Exec(sqlStr string) (sql.Result, error) {
	var err error
	var stmt *sql.Stmt
	if m.Tx == nil {
		stmt, err = Db.Prepare(sqlStr)
		defer stmt.Close()
	} else {
		stmt, err = m.Tx.Prepare(sqlStr)
	}
	if err != nil {
		fmt.Println("DBERR:", stmt, err, sqlStr)
		return nil, err
	}
	result, err := stmt.Exec()
	m.Clean()
	return result, err
}

func GetLastInsertId(result sql.Result) (int64, error) {
	return result.LastInsertId()
}

func GetRowsAffected(result sql.Result) (int64, error) {
	return result.RowsAffected()
}

func (m *Model) getSQLCond() string {
	if str, ok := m.Cond.(string); ok {
		return str
	}
	var strArr []string
	if arr, ok := m.Cond.(map[string]string); ok {
		var strArr []string
		for k, v := range arr {
			strArr = append(strArr, k+"='"+v+"'")
		}
		return " WHERE " + strings.Join(strArr, " && ")
	}
	if arr, ok := m.Cond.(map[string]interface{}); ok {
		for k, v := range arr {
			if isStr, ok := v.(string); ok {
				strArr = append(strArr, k+"='"+isStr+"'")
			}
			if isStrArr, ok := v.([]string); ok {
				for k, v := range isStrArr {
					isStrArr[k] = "'" + v + "'"
				}
				strArr = append(strArr, k+" in ("+strings.Join(isStrArr, ",")+")")
			}
		}
		return " WHERE " + strings.Join(strArr, " && ")
	}
	return ""
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

func (m *Model) getSQLLimite() string {
	if strArr, ok := m.Limit.([2]int); ok {
		return " LIMIT " + fmt.Sprintf("%d", strArr[0]) + ", " + fmt.Sprintf("%d", strArr[1])
	}
	return ""
}

func (m *Model) Clean() {
	m.Cond = nil
	m.Insert = nil
	m.Update = nil
	m.Field = ""
	m.Table = ""
	m.Index = ""
	m.Limit = nil
	m.Sort = ""
	m.IsDeug = 0
}

func (m *Model) Debug(sql string) {
	if m.IsDeug == 1 {
		fmt.Println(sql)
	}
}
