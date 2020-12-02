package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Kong/go-pdk"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/redis.v3"
)

//1.build
//go build -buildmode plugin  custom-rate-limiting.go
//2.将生成的.so文件放到go_plugins_dir定义的目录中
//cp -f custom-rate-limiting.so dir_to/plugins/
//3.不停止kong更新插件
//kong prepare && kong reload
//开发环境调试一句话命令
//go build -buildmode plugin  custom-rate-limiting.go && cp -f custom-rate-limiting.so ../plugins/ && kong prepare && kong reload

/*
json格式
[{
"type": "header,query,body",
"key": "keyName",
"value": "value1,value2,value3"
}, {
"type": "body",
"key": "orderId",
"value": "1,2,3"
}]
*/
//kong限流前缀
const rateLimitPrefix = "kong:customratelimit:"

//限流类型
const rateLimitType = "qps"

//匹配条件:or
const matchConditionOr = "or"

//匹配条件:and
const matchConditionAnd = "and"

//版本号
const version = "v0.1.1"

var ctx = context.Background()

//限流资源列表
var limitResourceList []limitResource

//Config kong 插件配置
type Config struct {
	QPS                 int    `json:"QPS" validate:"required,gte=0"` //请求限制的QPS值
	Log                 bool   `json:"Log" validate:"omitempty"`      //是否记录日志
	Path                string `json:"Path"`                          //资源路径
	LimitResourcesJson  string `json:"LimitResourcesJson"`            //流控规则选项，使用json配置，然后解析
	RedisHost           string `json:"RedisHost" validate:"required"`
	RedisPort           int    `json:"RedisPort" validate:"required,gte=1,lte=65535"`
	RedisAuth           string `json:"RedisAuth" validate:"omitempty"`
	RedisTimeoutSecond  int    `json:"RedisTimeoutSecond" validate:"required,gt=0"`
	RedisDB             int    `json:"RedisDB" validate:"omitempty,gte=0"`
	RedisLimitKeyPrefix string `json:"RedisLimitKeyPrefix" validate:"omitempty"`         //Redis限流key前缀
	HideClientHeader    bool   `json:"HideClientHeader" validate:"omitempty"`            //隐藏response header
	MatchCondition      string `json:"MatchCondition" validate:"omitempty,oneof=and or"` //流控规则匹配条件，and：所有规则都需要匹配到则成功，or: 匹配到一个则成功, 为空时默认为and
}

//限流资源
type limitResource struct {
	Type  string `json:"type"`  //限流类型，使用英文逗号分隔,如：header,query,body
	Key   string `json:"key"`   //限流key
	Value string `json:"value"` //限流值，使用英文逗号分隔，如：value1,value2,orderId1
}

func New() interface{} {
	return &Config{}
}

// kong Access phase
func (conf Config) Access(kong *pdk.PDK) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("kong plugin panic at: %v, err: %v", time.Now(), err)
			if kong == nil {
				log.Printf("kong fatal err ===> kong is nil at: %v", time.Now())
			} else {
				_ = kong.Log.Err(fmt.Sprint(err))
			}
		}
	}()
	_ = kong.Response.SetHeader("X-Rate-Limiting-Plugin-Version", version)
	unix := time.Now().Unix()

	if err := conf.checkConfig(); err != nil {
		_ = kong.Log.Err("[checkConfig] ", err.Error())
		return
	}

	//检查当前请求是否需要限流
	limitKey, matched := conf.checkNeedRateLimit(kong)
	if !matched {
		return
	}
	//获取限制标识identifier
	identifier, err := conf.getIdentifier(kong, limitKey)
	if err != nil {
		_ = kong.Log.Err("[getIdentifier] ", err.Error())
		return
	}
	remaining, stop, err := conf.getRemainingAndIncr(kong, identifier, unix)
	if err != nil {
		//出错只记录日志，不处理
		_ = kong.Log.Err("[getUsage] ", err.Error())
		return
	}
	//如果设置不隐藏header,则输出到header
	if !conf.HideClientHeader {
		_ = kong.Response.SetHeader("X-Rate-Limiting-Limit-QPS", strconv.Itoa(conf.QPS))
		_ = kong.Response.SetHeader("X-Rate-Limiting-Remaining", strconv.Itoa(remaining))
	}
	if stop {
		kong.Response.Exit(429, "API rate limit exceeded", nil)
		return
	}
}

//进入此插件，说明kong中已经启用插件
func (conf Config) checkConfig() error {
	validate := validator.New()
	err := validate.Struct(conf)
	if err != nil {
		return err
	}
	//如果MatchCondition为空，设置默认值为and
	if conf.MatchCondition == "" {
		conf.MatchCondition = matchConditionAnd
	}

	//允许流控规则为空
	if conf.LimitResourcesJson != "" {
		err = json.Unmarshal([]byte(conf.LimitResourcesJson), &limitResourceList)
		//json格式错误
		if err != nil {
			return errors.New(fmt.Sprintf("LimitResourcesJson with incorrect json format,%s", err.Error()))
		}
		//如果有值为空，则提示错误
		for _, item := range limitResourceList {
			if item.Type == "" || item.Key == "" || item.Value == "" {
				return errors.New("LimitResourcesJson with empty value")
			}
		}
	}
	if conf.Path != "" {
		//将QueryPath组装成一个limitResource类型，放入到limitResourceList统一处理
		queryPathLimitResource := limitResource{
			Type:  "Path",
			Key:   "path",
			Value: conf.Path,
		}
		limitResourceList = append(limitResourceList, queryPathLimitResource)
	}
	return nil
}

//获取剩余数量的同时加1
func (conf Config) getRemainingAndIncr(kong *pdk.PDK, identifier string, unix int64) (remaining int, stop bool, err error) {
	stop = false
	remaining = 0
	limitKey := conf.getRateLimitKey(identifier, unix)
	if conf.Log {
		_ = kong.Log.Err("[rateLimitKey] ", limitKey)
	}
	//第一次执行才设置有效期，如果过了有效期，则为下一时间段,使用lua保证原子性
	luaScript := `
		local key, value, expiration = KEYS[1], tonumber(ARGV[1]), ARGV[2]
		local newVal = redis.call("incrby", key, value)
		if newVal == value then
			redis.call("expire", key, expiration)
		end
		return newVal - 1
`
	redisClient := conf.newRedisClient()
	defer redisClient.Close()
	result, err := redisClient.Eval(ctx, luaScript, []string{limitKey}, 1, 1).Result()
	if err == redis.Nil {
		return remaining, stop, nil
	} else if err != nil {
		return remaining, stop, err
	} else {
		int64Usage := result.(int64)
		usageStr := strconv.FormatInt(int64Usage, 10)
		intUsage, err := strconv.Atoi(usageStr)
		if err != nil {
			return remaining, stop, err
		}
		remaining = conf.QPS - intUsage
		if remaining <= 0 {
			stop = true
			remaining = 0
		} else {
			//friendly show
			remaining -= 1
		}
		return remaining, stop, nil
	}
}

//获取限流key
func (conf Config) getRateLimitKey(identifier string, unix int64) string {
	return conf.getPrefix() + identifier + ":" + rateLimitType + ":" + strconv.FormatInt(unix, 10)
}

//获取限流标识符
func (conf Config) getIdentifier(kong *pdk.PDK, limitKey string) (string, error) {
	var identifier string
	consumer, err := kong.Client.GetConsumer()
	if err != nil {
		return "", err
	}
	service, err := kong.Router.GetService()
	if err != nil {
		return "", err
	}
	route, err := kong.Router.GetRoute()
	if err != nil {
		return "", err
	}
	if consumer.Id != "" {
		identifier += ":consumer:" + consumer.Id
	}
	if service.Id != "" {
		identifier += ":service:" + service.Id
	}
	if route.Id != "" {
		identifier += ":route:" + route.Id
	}
	identifier += ":" + limitKey
	return identifier, nil
}

//获取redis rate limit key prefix
func (conf Config) getPrefix() string {
	var prefix string
	//如果配置的RedisLimitKeyPrefix有：，则不处理，没有：则添加
	if conf.RedisLimitKeyPrefix == "" {
		return prefix + rateLimitPrefix
	}
	if strings.Contains(conf.RedisLimitKeyPrefix, ":") {
		prefix = conf.RedisLimitKeyPrefix
	} else {
		prefix = conf.RedisLimitKeyPrefix + ":"
	}
	return prefix + rateLimitPrefix
}

//redis客户端
func (conf Config) newRedisClient() *redis.Client {
	options := &redis.Options{
		Addr:        conf.RedisHost + ":" + strconv.Itoa(conf.RedisPort),
		Password:    conf.RedisAuth,
		DB:          conf.RedisDB,
		DialTimeout: time.Duration(conf.RedisTimeoutSecond) * time.Second,
	}
	return redis.NewClient(options)
}

//检查并返回是否需要限流的key
func (conf Config) checkNeedRateLimit(kong *pdk.PDK) (limitKey string, matched bool) {
	var matchedKey []string
	for _, limitResource := range limitResourceList {
		typeList := strings.Split(limitResource.Type, ",")
		valueList := strings.Split(limitResource.Value, ",")
		rateLimitValue, matched := conf.matchRateLimitValue(kong, limitResource.Key, typeList, valueList)
		//如果匹配到了是or关系，返回匹配成功(如果没有配置MatchCondition，默认会为空字符串，默认匹配条件为and)
		if matchConditionOr == conf.MatchCondition {
			if matched {
				return rateLimitValue, true
			}
		} else {
			//否则是and的关系，没有匹配到，返回匹配失败，否则加入到数组中
			if !matched {
				return "", false
			} else {
				matchedKey = append(matchedKey, rateLimitValue)
			}
		}
	}
	//如果limitResourceList为空(没有配置Path和LimitResourcesJson)，则返回匹配成功
	//如果全匹配，则转为字符串返回
	if len(limitResourceList) == len(matchedKey) {
		return strings.Join(matchedKey, ":"), true
	}
	return "", false
}

//match rate limit key
func (conf Config) matchRateLimitValue(kong *pdk.PDK, key string, typeList, valueList []string) (limitKey string, matched bool) {
	for _, limitType := range typeList {
		limitType = strings.ToLower(limitType)
		switch limitType {
		case "header":
			find, err := kong.Request.GetHeader(key)
			//获取失败，跳过
			if err != nil {
				continue
			}
			//如果请求头中存在被限制的列表，则返回
			if inSlice(find, valueList) {
				return find, true
			}
		case "query":
			find, err := kong.Request.GetQueryArg(key)
			//获取失败，跳过
			if err != nil {
				continue
			}
			//如果请求头中存在被限制的列表，则返回
			if inSlice(find, valueList) {
				return find, true
			}
		case "body":
			rawBody, err := kong.Request.GetRawBody()
			//获取失败，跳过
			if err != nil {
				continue
			}
			//TODO if json format or other raw format, maybe use contain judge or use equal after decode to key value pairs.
			if !strings.Contains(rawBody, key) {
				continue
			}
			bodySlice := strings.Split(rawBody, "&")
			for _, value := range valueList {
				limitValue := key + "=" + value
				if inSlice(limitValue, bodySlice) {
					return value, true
				}
			}
		case "path":
			find, err := kong.Request.GetPath()
			//获取失败，跳过
			if err != nil {
				continue
			}
			//如果在被限制的列表，则返回
			if inSlice(find, valueList) {
				return find, true
			}
		case "cookie":
			//not support
			continue
		case "ip":
			//next iteration will support
			continue
		default:
			continue
		}
	}
	return "", false
}

//是否在slice中
func inSlice(search string, slice []string) bool {
	for _, value := range slice {
		if value == search {
			return true
		}
	}
	return false
}
