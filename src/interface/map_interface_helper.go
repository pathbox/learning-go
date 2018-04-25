import (
	"fmt"
	"strconv"
)

// map[string]interface{}  将值转为string返回的helper方法
func MapValueGet(key string) string {
	v, ok := fm.Payload[key]
	if ok {
		switch v.(type) {
		case string:
			return v.(string)
		case int:
			return strconv.Itoa(v.(int))
		case int32:
			return strconv.Itoa(int32(v.(int)))
		case int64:
			return strconv.Itoa(int(v.(int64)))
		case float32:
			s := fmt.Sprintf("%v", v.(float32))
			return s
		case float64:
			s := fmt.Sprintf("%v", v.(float64))
			return s
		}
		return v.(string)
	} else {
		return code.BlankString
	}
}