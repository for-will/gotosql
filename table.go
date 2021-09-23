package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
)

func CreateTableSql(model interface{}) []string {
	typ := reflect.TypeOf(model)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	var nameLen, typeLen int
	var columns []*ColumnDesc
	var indiesMap = map[string]*IndexDesc{}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		col := &ColumnDesc{
			Name:       fieldName(field),
			Type:       fieldType(field),
			PrimaryKey: isPrimaryKey(field),
		}
		if col.Type == "unknown" {
			continue
		}
		if len(col.Name) > nameLen {
			nameLen = len(col.Name)
		}
		if len(col.Type) > typeLen {
			typeLen = len(col.Type)
		}

		columns = append(columns, col)

		// 索引&唯一索引
		indies := haveIndex(tableName(typ), field)
		for _, v := range indies {
			if idx, ok := indiesMap[v.Index]; ok {
				idx.Fields = append(idx.Fields, v.Fields...)
			} else {
				indiesMap[v.Index] = v
			}
		}
	}

	var sb strings.Builder

	var sqlList []string

	// Sql for drop table if exist
	sb.WriteString("drop table if exists ")
	sb.WriteString(tableName(typ))
	sb.WriteString(";")
	sqlList = append(sqlList, sb.String())
	sb.Reset()

	// Sql for create table
	sb.WriteString(fmt.Sprintf("create table %s\n(\n", tableName(typ)))
	format := fmt.Sprintf("    %%-%ds %%-%ds", nameLen, typeLen)
	for i, column := range columns {
		sb.WriteString(fmt.Sprintf(format, column.Name, column.Type))
		if column.PrimaryKey {
			sb.WriteString(" not null auto_increment\n")
			sb.WriteString("        primary key")
		} else {
			sb.WriteString(" not null")
		}
		if i != len(columns)-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}
	sb.WriteString(");")
	sqlList = append(sqlList, sb.String())
	sb.Reset()

	// Generate sql for create index
	for _, v := range indiesMap {
		if v.Unique {
			sb.WriteString("create unique index ")
		} else {
			sb.WriteString("create index ")
		}
		sb.WriteString(v.Index)
		sb.WriteString(" on ")
		sb.WriteString(tableName(typ))
		sb.WriteString(" (")
		for i, col := range v.Fields {
			if i != 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(col)
		}
		sb.WriteString(");")
		sqlList = append(sqlList, sb.String())
		sb.Reset()
	}

	return sqlList
}

func RecreateTable(db *sql.DB, model interface{}) {

	queries := CreateTableSql(model)
	for _, s := range queries {
		fmt.Println("Exec Sql:", s)
		_, err := db.Exec(s)
		if err != nil {
			panic(err)
		}
	}
}

func GenRecreateTableFunc(model interface{}) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("func NewTable%s(db *sql.DB) {\n", ModelName(model)))
	sb.WriteString("    sqlList := []string{\n")
	sqlList := CreateTableSql(model)
	for _, s := range sqlList {
		sb.WriteString("\t\t`")
		sb.WriteString(s)
		sb.WriteString("`,\n")
	}
	sb.WriteString("}\n")
	sb.WriteString(`
		for _, s := range sqlList {
		fmt.Println("Exec Sql:", s)
		_, err := db.Exec(s)
		if err != nil {
			fmt.Printf("db error: %+v", err)
		}
	}
}`)
	return sb.String()
}

func GenModelAutoFile(file string, models ...interface{}) {
	var sb strings.Builder
	sb.WriteString(`//
// Code generated auto. DO NOT EDIT.

package db

import (
	"database/sql"
	"fmt"
)`)

	for _, m := range models {
		sb.WriteString("\n\n")
		sb.WriteString(GenRecreateTableFunc(m))
	}
	ioutil.WriteFile(file, []byte(sb.String()), 0664)
}
