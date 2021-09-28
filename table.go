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
			ColName:    fieldName(field),
			Type:       fieldType(field),
			PrimaryKey: isPrimaryKey(field),
		}
		if col.Type == "unknown" {
			continue
		}
		if len(col.ColName) > nameLen {
			nameLen = len(col.ColName)
		}
		if len(col.Type) > typeLen {
			typeLen = len(col.Type)
		}

		columns = append(columns, col)

		// 索引&唯一索引
		indies := haveIndex(tableName(typ), field)
		for _, v := range indies {
			if idx, ok := indiesMap[v.Index]; ok {
				idx.Keys = append(idx.Keys, v.Keys...)
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
		sb.WriteString(fmt.Sprintf(format, column.ColName, column.Type))
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
		for i, col := range v.Keys {
			if i != 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(col.KeyName)
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

func GenModelAutoFile(file string, pkg string, models ...interface{}) {
	var sb strings.Builder
	sb.WriteString("//\n")
	sb.WriteString("// Code generated auto. DO NOT EDIT.\n\n")
	sb.WriteString("package " + pkg)
	sb.WriteString(`

import (
	"database/sql"
	"strings"
)`)

	for _, m := range models {
		model := Model(m)
		sb.WriteString("\n")
		sb.WriteString(model.BuildCreateTableFunc())
		sb.WriteString("\n")
		sb.WriteString(model.BuildSaveFunc())
		sb.WriteString("\n")
		sb.WriteString(model.BuildFindOneFunc())
		sb.WriteString("\n")
		sb.WriteString(model.BuildFindFunc())
	}
	ioutil.WriteFile(file, []byte(sb.String()), 0664)
}
