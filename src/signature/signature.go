package signature

import (
	"crypto/sha1"
	"encoding/hex"
	"sort"
	"strconv"
)

type SignParamsPair struct {
	Key   string
	Value interface{}
}

type SignParams []SignParamsPair

// 实现sort接口
func (p SignParams) Len() int {
	return len(p)
}

func (p SignParams) Less(i, j int) bool {
	if p[i].Key < p[k].Key {
		return true
	}
	return false
}

func (p SignParams) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func VerifySign(params map[string]interface{}, priKey, signature string) bool {
	return Sign(params, priKey) == signature
}

// Sign 计算签名
func Sign(params map[string]interface{}, priKey string) string {
	str := map2String(map2Sort(params)) + priKey
	hashed := sha1.Sum([]byte(str))

	return hex.EncodeToString(hashed[:])
}

func map2String(params SignParams) (sortStr string) {
	for _, v := range params {
		switch v.Value.(type) {
		case string, int, int8, int16, int32, int64, float64, float32:
			sortStr += v.Key + simple2String(v.Value)
			break
		case []interface{}:
			sortStr += v.Key + array2String(v.Value.([]interface{}))
			break
		case map[string]interface{}:
			paramsSort := map2Sort(v.Value.(map[string]interface{}))
			sortStr += v.Key + map2String(paramsSort) // 嵌套key进行递归
			break
		}
	}
	return sortStr
}

func simple2String(v interface{}) string {
	switch v.(type) {
	case string:
		return v.(string)
	case int:
		return strconv.FormatInt(int64(v.(int)), 10)
	case int8:
		return strconv.FormatInt(int64(v.(int8)), 10)
	case int16:
		return strconv.FormatInt(int64(v.(int16)), 10)
	case int32:
		return strconv.FormatInt(int64(v.(int32)), 10)
	case int64:
		return strconv.FormatInt(v.(int64), 10)
	case float32:
		return strconv.FormatFloat(float64(v.(float32)), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(v.(float64), 'f', -1, 64)
	}
	return ""
}

// 把嵌套的map 转为SignParamsPair，然后根据key排序
func map2Sort(params map[string]interface{}) SignParams {
	paramsSort := make(SignParams, 0)
	for k, v := range params {
		paramsSort = append(paramsSort, SignParamsPair{k, v})
	}
	sort.Sort(paramsSort)
	return paramsSort
}

func array2String(params []interface{}) (sortStr string) {
	for _, v := range params {
		switch v.(type) {
		case string, int, int8, int16, int32, int64, float64, float32:
			sortStr += simple2String(v)
			break
		case []interface{}:
			sortStr += array2String(v.([]interface{}))
			break
		case map[string]interface{}:
			sortStr += map2String(map2Sort(v.(map[string]interface{})))
			break
		}
	}

	return sortStr
}
