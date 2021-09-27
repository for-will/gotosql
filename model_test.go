package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"math/rand"
	"testing"
	"time"
)

var db *sql.DB

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
	db = openDb()
}

func openDb() *sql.DB {
	db, err := sql.Open("mysql", "puffer:puffer123@tcp(127.0.0.1:3306)/gotosql?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	return db
}

func TestTableModel_SelectSql(t *testing.T) {
	model := Model(&RewardTask{})
	//s := Model(&RewardTask{}).SelectSql()
	//fmt.Println(Model(&RewardTask{}).ScanRow("out"))
	//s = Model(&RewardTask{}).ScanRow("out")
	s := model.BuildFindOneFunc()
	fmt.Println(s)
}

func TestTableModel_InsertSql(t *testing.T) {
	t.Log(Model(&TestTableTask{}).InsertSql())
}

func TestTableModel_BuildSaveFunc(t *testing.T) {
	t.Log(Model(&TestTableTask{}).BuildSaveFunc())
}

func TestTableModel_CreateTableSql(t *testing.T) {
	println(Model(&TestTableTask{}).CreateTableSql())
}
