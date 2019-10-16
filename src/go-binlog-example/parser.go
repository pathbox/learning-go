package main

import (
	"fmt"
	"github.com/json-iterator/go"
	"github.com/siddontang/go-mysql/canal"
	"github.com/siddontang/go-mysql/schema"
	"reflect"
	"strings"
	"time"
)

type BinlogParser struct{}

func (m *BinlogParser) GetBinLogData(element interface{}, e *canal.RowsEvent, n int) error {
	var columnName string
	var ok bool
	v := reflect.ValueOf(element)
	s := reflect.Indirect(v)
	t := s.Type()
	num := t.NumField()
	for k := 0; k < num; k++ {
		parsedTag := parseTagSetting(t.Field(k).Tag)
		name := s.Filed(k).Type().Name()

		if columnName, ok = parsedTag["COLUMN"]; !ok || columnName == "COLUMN" {
			continue
		}

		switch name {
		case "bool":
			s.Field(k).SetBool(m.boolHelper(e, n, columnName))
		case "int":
			s.Field(k).SetInt(m.intHelper(e, n, columnName))
		case "string":
			s.Field(k).SetString(m.stringHelper(e, n, columnName))
		case "Time":
			timeVal := m.dateTimeHelper(e, n, columnName)
			s.Field(k).Set(reflect.ValueOf(timeVal))
		case "float64":
			s.Field(k).SetFloat(m.floatHelper(e, n, columnName))
		default:
			if _, ok := parsedTag["FROMJSON"]; ok {

				newObject := reflect.New(s.Field(k).Type()).Interface()
				json := m.stringHelper(e, n, columnName)

				jsoniter.Unmarshal([]byte(json), &newObject)

				s.Field(k).Set(reflect.ValueOf(newObject).Elem().Convert(s.Field(k).Type()))
			}
		}
	}
	return nil
}

func (m *BinlogParser) dateTimeHelper(e *canal.RowsEvent, n int, columnName string) time.Time {

	columnId := m.getBinlogIdByName(e, columnName)
	if e.Table.Columns[columnId].Type != schema.TYPE_TIMESTAMP {
		panic("Not dateTime type")
	}
	t, _ := time.Parse("2006-01-02 15:04:05", e.Rows[n][columnId].(string))

	return t
}

func (m *BinlogParser) intHelper(e *canal.RowsEvent, n int, columnName string) int64 {

	columnId := m.getBinlogIdByName(e, columnName)
	if e.Table.Columns[columnId].Type != schema.TYPE_NUMBER {
		return 0
	}

	switch e.Rows[n][columnId].(type) {
	case int8:
		return int64(e.Rows[n][columnId].(int8))
	case int32:
		return int64(e.Rows[n][columnId].(int32))
	case int64:
		return e.Rows[n][columnId].(int64)
	case int:
		return int64(e.Rows[n][columnId].(int))
	case uint8:
		return int64(e.Rows[n][columnId].(uint8))
	case uint16:
		return int64(e.Rows[n][columnId].(uint16))
	case uint32:
		return int64(e.Rows[n][columnId].(uint32))
	case uint64:
		return int64(e.Rows[n][columnId].(uint64))
	case uint:
		return int64(e.Rows[n][columnId].(uint))
	}
	return 0
}

func (m *BinlogParser) floatHelper(e *canal.RowsEvent, n int, columnName string) float64 {

	columnId := m.getBinlogIdByName(e, columnName)
	if e.Table.Columns[columnId].Type != schema.TYPE_FLOAT {
		panic("Not float type")
	}

	switch e.Rows[n][columnId].(type) {
	case float32:
		return float64(e.Rows[n][columnId].(float32))
	case float64:
		return float64(e.Rows[n][columnId].(float64))
	}
	return float64(0)
}

func (m *BinlogParser) boolHelper(e *canal.RowsEvent, n int, columnName string) bool {

	val := m.intHelper(e, n, columnName)
	if val == 1 {
		return true
	}
	return false
}

func (m *BinlogParser) stringHelper(e *canal.RowsEvent, n int, columnName string) string {

	columnId := m.getBinlogIdByName(e, columnName)
	if e.Table.Columns[columnId].Type == schema.TYPE_ENUM {

		values := e.Table.Columns[columnId].EnumValues
		if len(values) == 0 {
			return ""
		}
		if e.Rows[n][columnId] == nil {
			//Если в енум лежит нуул ставим пустую строку
			return ""
		}

		return values[e.Rows[n][columnId].(int64)-1]
	}

	value := e.Rows[n][columnId]

	switch value := value.(type) {
	case []byte:
		return string(value)
	case string:
		return value
	}
	return ""
}

func (m *BinlogParser) getBinlogIdByName(e *canal.RowsEvent, name string) int {
	for id, value := range e.Table.Columns {
		if value.Name == name {
			return id
		}
	}
	panic(fmt.Sprintf("There is no column %s in table %s.%s", name, e.Table.Schema, e.Table.Name))
}

func parseTagSetting(tags reflect.StructTag) map[string]string {
	settings := map[string]string{}
	for _, str := range []string{tags.Get("sql"), tags.Get("gorm")} {
		tags := strings.Split(str, ";")
		for _, value := range tags {
			v := strings.Split(value, ":")
			k := strings.TrimSpace(strings.ToUpper(v[0]))
			if len(v) >= 2 {
				settings[k] = strings.Join(v[1:], ":")
			} else {
				settings[k] = k
			}
		}
	}
	return settings
}