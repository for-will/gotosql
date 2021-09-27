package db

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"text/template"
	"time"
)

type ColumnDesc struct {
	FieldName  string
	FieldType  string
	ColName    string
	Type       string
	Default    string
	PrimaryKey bool
}

type IndexDesc struct {
	Index  string
	Fields []string
	Unique bool
}

func Model(m interface{}) *TableModel {

	typ := reflect.TypeOf(m)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	model := &TableModel{
		ModName:  ModelName(m),
		TblName:  tableName(typ),
		Columns:  modelFields(typ),
		LogSql:   "LogSql",
		LogError: "LogError",
		Indies:   modelIndies(typ),
	}
	for _, col := range model.Columns {
		if col.PrimaryKey {
			model.Pk = col
		}
	}
	return model
}

type TableModel struct {
	ModName      string
	TblName      string
	QuerySql     string
	QuerySqlFmt  string
	QueryArgs    string
	ScanReceiver string
	LogError     string
	LogSql       string
	Columns      []*ColumnDesc
	Pk           *ColumnDesc
	Indies       map[string]*IndexDesc
}

func (tm *TableModel) CreateTableSql() string {

	var nameLen, typeLen int
	for _, col := range tm.Columns {
		if len(col.ColName) > nameLen {
			nameLen = len(col.ColName)
		}
		if len(col.Type) > typeLen {
			typeLen = len(col.Type)
		}
	}

	var sb strings.Builder

	// Sql for create table
	sb.WriteString(fmt.Sprintf("create or replace table %s\n(\n", tm.TblName))
	format := fmt.Sprintf("    %%-%ds %%-%ds", nameLen, typeLen)
	for i, col := range tm.Columns {
		sb.WriteString(fmt.Sprintf(format, col.ColName, col.Type))
		if col.PrimaryKey {
			sb.WriteString(" not null auto_increment\n")
			sb.WriteString("        primary key")
		} else {
			sb.WriteString(" not null")
		}
		if col.Default != "" {
			sb.WriteString(" default ")
			sb.WriteString(col.Default)
		}
		if i != len(tm.Columns)-1 {
			sb.WriteString(",\n")
		}
	}

	if len(tm.Indies) == 0 {
		sb.WriteString("\n")
	}

	// Generate sql for create index
	for _, v := range tm.Indies {
		sb.WriteString(",\n")
		if v.Unique {
			sb.WriteString("\tunique index ")
		} else {
			sb.WriteString("\tindex ")
		}
		sb.WriteString(v.Index)
		sb.WriteString("(")
		for i, col := range v.Fields {
			if i != 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(col)
		}
		sb.WriteString(") using btree")
	}

	sb.WriteString("\n);")

	return sb.String()
}

func (tm *TableModel) SelectSql() string {
	var sb strings.Builder
	sb.WriteString("SELECT ")
	for i, f := range tm.Columns {
		if i != 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(f.ColName)
	}
	sb.WriteString(" FROM ")
	sb.WriteString(tm.TblName)
	return sb.String()
}

func (tm *TableModel) InsertSql() string {
	var sb strings.Builder
	sb.WriteString("INSERT INTO ")
	sb.WriteString(tm.TblName)
	sb.WriteString("(")

	var first = true
	for _, col := range tm.Columns {
		if col.PrimaryKey {
			continue
		}
		if !first {
			sb.WriteString(", ")
		}
		sb.WriteString(col.ColName)
		first = false
	}
	sb.WriteString(") VALUE (")
	first = true
	for _, col := range tm.Columns {
		if col.PrimaryKey {
			continue
		}
		if !first {
			sb.WriteString(", ")
		}
		sb.WriteString("?")
		first = false
	}
	sb.WriteString(")")
	return sb.String()
}

func (tm *TableModel) ScanRow(receiver string) string {
	var sb strings.Builder
	addr := "&" + receiver + "."
	for i, f := range tm.Columns {
		if i != 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(addr)
		sb.WriteString(f.FieldName)
	}
	return sb.String()
}

func (tm *TableModel) InsertArgs(obj string) string {
	var sb strings.Builder
	addr := obj + "."
	var first = true
	for _, f := range tm.Columns {
		if f.PrimaryKey {
			continue
		}

		if !first {
			sb.WriteString(", ")
		}
		sb.WriteString(addr)
		sb.WriteString(f.FieldName)
		first = false
	}
	return sb.String()
}

func (tm *TableModel) BuildCreateTableFunc() string {

	tm.QuerySql = tm.CreateTableSql()
	tpl := template.New(tm.TblName)
	tpl.Parse(`
func NewTbl{{.ModName}}(db *sql.DB) {

	querySql := ` + "`{{.QuerySql}}`" + `
	{{.LogSql}}(querySql)
	_, err := db.Exec(querySql)
	if err != nil {
		{{.LogError}}("%+v", err)
	}
}`)
	buf := bytes.NewBuffer(make([]byte, 0, 1024*16))

	template.Must(tpl, tpl.Execute(buf, tm))
	return buf.String()
}

func (tm *TableModel) BuildFindOneFunc() string {
	tm.QuerySql = tm.SelectSql()
	tm.ScanReceiver = tm.ScanRow("out")

	tpl := template.New(tm.TblName)
	tpl.Parse(`
func First{{.ModName}}(db *sql.DB, out *{{.ModName}}, cond string, args ...interface{}) error {

	sfmt := strings.Replace(cond, "?", "'%+v'", -1)
	{{.LogSql}}("{{.QuerySql}} WHERE "+sfmt+" LIMIT 1", args...)
	rows, err := db.Query("{{.QuerySql}} WHERE "+cond+" LIMIT 1", args...)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		{{.LogError}}("db query error: %+v", err)
		return err
	}

	if rows.Next() {
		if err = rows.Scan({{.ScanReceiver}}); err != nil {
			{{.LogError}}("db query Scan error: %+v", err)
			return err
		}
	} else {
		return sql.ErrNoRows
	}
	return nil
}`)
	buf := bytes.NewBuffer(make([]byte, 0, 1024*16))

	template.Must(tpl, tpl.Execute(buf, tm))
	return buf.String()
}

func (tm *TableModel) BuildFindFunc() string {
	tm.QuerySql = tm.SelectSql()
	tm.ScanReceiver = tm.ScanRow("obj")

	tpl := template.New(tm.TblName)
	tpl.Parse(`
func Find{{.ModName}}(db *sql.DB, out *[]*{{.ModName}}, cond string, args ...interface{}) error {

	sfmt := strings.Replace(cond, "?", "'%+v'", -1)
	{{.LogSql}}("{{.QuerySql}} WHERE "+sfmt, args...)
	rows, err := db.Query("{{.QuerySql}} WHERE "+cond, args...)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		{{.LogError}}("db query error: %+v", err)
		return err
	}

	for rows.Next() {
		obj := &{{.ModName}}{}
		if err = rows.Scan({{.ScanReceiver}}); err != nil {
			{{.LogError}}("db query Scan error: %+v", err)
			return err
		}
		*out = append(*out, obj)
	}
	return nil
}`)
	buf := bytes.NewBuffer(make([]byte, 0, 1024*16))

	template.Must(tpl, tpl.Execute(buf, tm))
	return buf.String()
}

func (tm *TableModel) BuildSaveFunc() string {
	tm.QuerySql = tm.InsertSql()
	tm.QuerySqlFmt = strings.Replace(tm.QuerySql, "?", "'%v'", -1)
	tm.QueryArgs = tm.InsertArgs("obj")

	tpl := template.New(tm.TblName)
	tpl.Parse(`
func Save{{.ModName}}(db *sql.DB, obj *{{.ModName}}) error {

	{{.LogSql}}("{{.QuerySqlFmt}}",
		{{.QueryArgs}})
	result, err := db.Exec("{{.QuerySql}}",
		{{.QueryArgs}})

	if err != nil {
		{{.LogError}}("db query error: %+v", err)
		return err
	}

	if id, err := result.LastInsertId(); err != nil {
		{{.LogError}}("get last insert id error: %+v", err)
		return err
	} else {
		obj.{{.Pk.FieldName}} = {{.Pk.FieldType}}(id)
	}
	return nil
}`)
	buf := bytes.NewBuffer(make([]byte, 0, 1024*16))

	template.Must(tpl, tpl.Execute(buf, tm))
	return buf.String()
}

func tableName(typ reflect.Type) string {
	words := strings.Split(typ.String(), ".")
	return SnakeCase(words[len(words)-1])
}

func fieldName(f reflect.StructField) string {
	return SnakeCase(f.Name)
}

func fieldType(f reflect.StructField) string {
	t := f.Type
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	tag := f.Tag.Get("db")

	words := strings.Split(tag, ",")
	for _, word := range words {
		kv := strings.Split(word, ":")
		if kv[0] == "type" {
			return kv[1]
		}
	}

	if t == reflect.TypeOf(time.Time{}) {
		return "datetime"
	}

	switch t.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64:
		return "int"

	case reflect.Int8:
		return "tinyint"

	case reflect.String:
		return "varchar(64)"
	}
	return "unknown"
}

func columnDefault(f reflect.StructField) string {

	tag := f.Tag.Get("db")

	words := strings.Split(tag, ",")
	for _, word := range words {
		kv := strings.Split(word, ":")
		if kv[0] == "default" && len(kv) == 2 {
			return kv[1]
		}
	}
	return ""
}

func isPrimaryKey(f reflect.StructField) bool {
	tag := f.Tag.Get("db")

	words := strings.Split(tag, ",")
	for _, word := range words {
		if word == "primary_key" {
			return true
		}
	}
	return false
}

func haveIndex(tbn string, f reflect.StructField) (indies []*IndexDesc) {
	tag := f.Tag.Get("db")

	words := strings.Split(tag, ",")
	for _, word := range words {
		kv := strings.Split(word, ":")
		if kv[0] == "index" {
			index := &IndexDesc{
				Unique: false,
			}
			index.Fields = append(index.Fields, fieldName(f))
			if len(kv) == 1 {
				index.Index = "idx_" + tbn + "_" + fieldName(f)
			} else {
				index.Index = kv[1]
			}
			indies = append(indies, index)

		} else if kv[0] == "unique" {
			unique := &IndexDesc{
				Unique: true,
			}
			unique.Fields = append(unique.Fields, fieldName(f))
			if len(kv) == 1 {
				unique.Index = "uni_" + tbn + "_" + fieldName(f)
			} else {
				unique.Index = kv[1]
			}
			indies = append(indies, unique)
		}
	}
	return indies
}

func modelFields(typ reflect.Type) []*ColumnDesc {
	var nameLen, typeLen int
	var columns []*ColumnDesc
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		col := &ColumnDesc{
			FieldName:  field.Name,
			FieldType:  field.Type.Kind().String(),
			ColName:    fieldName(field),
			Type:       fieldType(field),
			Default:    columnDefault(field),
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
	}
	return columns
}

func modelIndies(typ reflect.Type) map[string]*IndexDesc {

	var indiesMap = map[string]*IndexDesc{}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if fieldType(field) == "unknown" {
			continue
		}

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
	return indiesMap
}

/*Test Model*/

type RewardTask struct {
	ReturnCode  int32 `db:"index"`
	Id          int32 `db:"primary_key"`
	ExpAdd      int32
	Gold        int32
	Honor       int32
	Achievement int32
	Diamond     int32
}

type TestTableTask struct {
	Id         int32     `db:"name:sn,type:int,primary_key"`
	PlayerSn   int32     `db:"unique:uni_player_mission"`
	Mission    int32     `db:"unique:uni_player_mission"`
	State      int8      `db:"index:idx_state_progress"`
	Progress   int32     `db:"index:idx_state_progress"`
	RewardedAt time.Time `db:"type:timestamp,default:from_unixtime(1)"`
	PsX        int32     `db:"default:-1"`
	PsY        int32     `db:"default:-1"`
}
