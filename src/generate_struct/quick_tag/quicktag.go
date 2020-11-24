package quicktag

import (
	_ "encoding/json"
	"fmt"
	. "reflect"
	"sort"
	"strconv"
	"strings"
)

var (
	StyleConvert    func(string) string = PascalToUnderline
	MaxSelfRefLevel                     = 5
	TagNames                            = []string{"json", "bson"}
	byteType                            = TypeOf(uint8(0))
	dynTypeMap                          = make(map[Type]Type)
)

func Q(x interface{}) interface{} {
	t := TypeOf(x)
	nt := dynTypeMap[t]
	if nt == nil {
		if t.Kind() == Ptr { // 如果是指针 Kind 取Elem
			t = t.Elem()
		}
		if t.Kind() != Struct {
			return x
		}
		nt = DynamicType(t)
	}
	return TypeCast(x, nt)
}

func PascalToUnderline(s string) string {
	if len(s) <= 1 {
		return strings.ToLower(s)
	}
	isUpper := func(c byte) bool {
		return c >= 'A' && c <= 'Z'
	}
	toLower := func(c byte) byte {
		if isUpper(c) {
			return c - 'A' + 'a'
		}
		return c
	}
	r := strings.Builder{}
	r.WriteByte(toLower(s[0]))
	i, l := 1, len(s)
	for i < l {
		if isUpper(s[i]) {
			if !isUpper(s[i-1]) {
				r.WriteByte('_')
			} else if i < l-1 && !isUpper(s[i+1]) {
				r.WriteByte('_')
			}
			r.WriteByte(toLower(s[i]))
		} else {
			r.WriteByte(s[i])
		}
		i += 1
	}
	return r.String()
}

func DynamicType(t Type) Type {
	if t.Kind() != Struct {
		panic("must struct type")
	}
	nt := dynTypeMap[t]
	if nt == nil {
		nt = genType(t, &genRecords{
			tmap: make(map[Type]int),
			max:  MaxSelfRefLevel,
		})
		dynTypeMap[t] = nt
		dynTypeMap[PtrTo(t)] = nt
	}
	return nt
}

type genRecords struct {
	tmap map[Type]int
	max  int
}

func genType(t Type, r *genRecords) Type {
	switch t.Kind() {
	case Ptr:
		gt := genType(t.Elem(), r)
		if gt == nil {
			return nil
		}
		return PtrTo(gt)
	case Struct:
		break
	case Map:
		if t.Key().Kind() != String {
			panic("map key must be string")
		}
		gt := genType(t.Elem(), r)
		if gt == nil {
			return nil
		}
		return MapOf(t.Key(), gt)
	case Slice:
		gt := genType(t.Elem(), r)
		if gt == nil {
			return nil
		}
		return SliceOf(gt)
	case Array:
		gt := genType(t.Elem(), r)
		if gt == nil {
			return nil
		}
		return ArrayOf(t.Len(), gt)
	default:
		return t
	}
	if r.tmap[t] > r.max {
		return nil
	}
	r.tmap[t] += 1
	fs := make([]StructField, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tm := parseTag(f.Tag)
		customJSON := false
		if tm["json"] != "" {
			customJSON = true
		}
		autoName := StyleConvert(f.Name)
		for _, tname := range TagNames {
			if tm[tname] == "" {
				tm[tname] = autoName
			}
		}
		var gt Type
		if tm["quicktag"] == "-" {
			gt = f.Type
		} else {
			gt = genType(f.Type, r)
		}
		if gt == nil {
			blankTm := make(map[string]string)
			for _, tname := range TagNames {
				blankTm[tname] = "-"
			}
			fs = append(fs, StructField{
				Name: f.Name,
				Type: ArrayOf(int(f.Type.Size()), byteType),
				Tag:  makeTag(blankTm),
			})
		} else {
			if f.Anonymous && !customJSON {
				delete(tm, "json")
			}
			fs = append(fs, StructField{
				Name:      f.Name,
				Type:      gt,
				Tag:       makeTag(tm),
				Anonymous: f.Anonymous,
			})
		}
	}
	r.tmap[t] -= 1

	nt := StructOf(fs)

	return nt
}

//copied from go reflect
func parseTag(stag StructTag) map[string]string {
	tag := string(stag)
	m := make(map[string]string)
	for tag != "" {
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i += 1
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		value, err := strconv.Unquote(qvalue)
		if err != nil {
			continue
			m[name] = value
		}
		return m
	}
	return m
}
func makeTag(m map[string]string) StructTag {
	if len(m) == 0 {
		return ""
	}
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	s := ""
	for _, k := range ks {
		s += fmt.Sprintf(`%s:"%s" `, k, m[k])
	}
	return StructTag(s)
}
