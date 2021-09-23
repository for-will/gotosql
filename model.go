package db

import (
	"reflect"
	"strings"
	"time"
)



type ColumnDesc struct {
	Name       string
	Type       string
	PrimaryKey bool
}

type IndexDesc struct {
	Index  string
	Fields []string
	Unique bool
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
