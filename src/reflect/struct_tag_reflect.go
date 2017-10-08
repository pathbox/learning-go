package main

import (
	"fmt"
	"reflect"
	"strings"
)

// params tag
type M map[string]interface{}

// 将m中的值赋给ptr指向的struct的相应字段
func (m M) AssignTo(ptr interface{}, tagName string) bool {
	v := reflect.ValueOf(ptr)
	if v.IsValid() == false {
		panic("not valid")
	}
	//找到最后指向的值，或者空指针，空指针是需要进行初始化的
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}

	tv := v
	if tv.Kind() == reflect.Ptr && tv.CanSet() {
		//对空指针进行初始化，暂时用临时变量保存
		tv.Set(reflect.New(tv.Type().Elem()))
		tv = tv.Elem()
	}

	if tv.Kind() != reflect.Struct {
		panic("not struct")
	}
	fmt.Println(tv)
	if assign(tv, m, tagName) { //赋值成功，将临时变量赋给原值
		if v.Kind() == reflect.Ptr {
			v.Set(tv.Addr())
		} else {
			v.Set(tv)
		}
		return true
	} else {
		return false
	}
}

//将src中的值填充到dstValue中
func assign(dstVal reflect.Value, src interface{}, tagName string) bool {
	sv := reflect.ValueOf(src)
	if !dstVal.IsValid() || !sv.IsValid() {
		return false
	}

	if dstVal.Kind() == reflect.Ptr {
		//初始化空指针
		if dstVal.IsNil() && dstVal.CanSet() {
			dstVal.Set(reflect.New(dstVal.Type().Elem()))
		}
		dstVal = dstVal.Elem()
	}

	// 判断可否赋值，小写字母开头的字段、常量等不可赋值
	if !dstVal.CanSet() {
		return false
	}

	switch dstVal.Kind() {
	case reflect.Bool: //TODO...
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64: //TODO...
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64: //TODO...
	case reflect.String: //TODO...
	case reflect.Slice: //TODO...
	case reflect.Map: //TODO...
	case reflect.Struct:
		if sv.Kind() != reflect.Map || sv.Type().Key().Kind() != reflect.String {
			return false
		}

		success := false
		for i := 0; i < dstVal.NumField(); i++ {
			fv := dstVal.Field(i)
			if fv.IsValid() == false || fv.CanSet() == false {
				continue
			}

			ft := dstVal.Type().Field(i)
			name := ft.Name
			strs := strings.Split(ft.Tag.Get(tagName), ",")
			if strs[0] == "-" { //处理ignore的标志
				continue
			}

			if len(strs[0]) > 0 {
				name = strs[0]
			}
			fmt.Println(name)
			fsv := sv.MapIndex(reflect.ValueOf(name))
			fmt.Println("fsv: ", fsv.IsValid())
			if fsv.IsValid() {
				fmt.Println("fv:", fv.Kind())
				if fv.Kind() == reflect.Ptr && fv.IsNil() {
					pv := reflect.New(fv.Type().Elem())
					if assign(pv, fsv.Interface(), tagName) {
						fmt.Println("set", fv, pv)
						fv.Set(pv)
						success = true
					}
				} else {
					if assign(fv, fsv.Interface(), tagName) {
						success = true
					}
				}
			} else if ft.Anonymous {
				//尝试对匿名字段进行递归赋值，跟JSON的处理原则保持一致
				if fv.Kind() == reflect.Ptr && fv.IsNil() {
					pv := reflect.New(fv.Type().Elem())
					if assign(pv, src, tagName) {
						fv.Set(pv)
						success = true
					}
				} else {
					if assign(fv, src, tagName) {
						success = true
					}
				}
			}
		}
		return success
	default:
		return false
	}

	return true
}

type Room struct {
	ID   int    `param:"-" json:"id"`
	Name string `param:"name" json:"name"`
}

type School struct {
	ID     int    `param:"-" json:"id"`
	Name   string `param:"name" json:"name"`
	RoomID int    `param:"room_id" json:"-"`
	Room   *Room  `param:"-" json:"room"`
}

func main() {
	params := M{"id": "123", "name": "Primary School", "room_id": "1"}
	var s *School
	params.AssignTo(&s, "param")
	fmt.Println(s)
}

// id字段其实没有赋值
