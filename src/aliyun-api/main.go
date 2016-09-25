package main
import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    // "encoding/json"  //用于解析获取到的json
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    // "reflect"  //用反射的方式获取解析到的json数据里的值
    "net/url"
    "sort"
)

var HTTPMETHOD string = "" // POST或者GET，请求的方法
var X_Ca_Key string = "X_Ca_Key" // 参与加签的头部

//阿里云服务器在创建应用的时候会分配appkey和appsecret
var SECRET string = "" // appkey
var KEY string = "" // appsecret
var LF string = "\n"

// url 由uri+path+params构成
var URI string = ""
var PATh string = ""

func main() {
  var apidata = map[string]map[string] string{
    "params":{
      "httpMethod": HTTPMETHOD,
      "path": PATH,
      "X-Ca-Signature-Headers": X_Ca_Key,
      "X-Ca-Signature":         "",
      "secret":                 SECRET,
      "LF":                     LF,
      "header_data":            "",
      "path_value":             "",
      "path_url":               "",
      "uri":                    URI,
      "url":                    "",
    },
    "RequestHeader": {
        "X-Ca-Key":               KEY,
        "X-Ca-Signature-Headers": X_Ca_Key,
        "X-Ca-Signature":         "",
        // "X-Ca-Request-Mode":      "debug", //用于调试
    },
    "HeaderValue": {
        "请求的参数名key": "请求的参数值value",
    },
    "HeaderData": {
        "X-Ca-Key": KEY,//为什么这里又有个KEY，是因为这个是设置加签的头部信息的数据
    },
  }
  var data string = GetData(&apidata)
  fmt.Println(data)
}

//请求api服务

func GetData(params *map[string]map[string]string) string{
  GetSign(params)
  client := http.Client{}
  request, err := http.NewRequest(
    (*params)["params"]["httpMethod"],
    (*params)["params"]["url"], nil
  )
  if err != nil {
      return err
  }

  for key, value := range (*params)["RequestHeader"]{
    request.Header.Add(key, value)
  }
  response, err_response := client.Do(request)
  if err_response != nil {
    return err_response
  }
  defer response.Body.Close()
  body, err_body := ioutil.ReadAll(response.Body)
  if err_body != nil{
    return err_body
  }
  // for key, value := range response.Header {
    //     fmt.Println(key, "  ", value)
    // }
    return string(body)

}
//拼接url 针对参数值中有中文的情况，url进行urlencode转码。由于刚开始写的时候没有注意到中文的问题，所以这个函数是后来仓促加的。
//参与url拼接的数据在api["HeaderValue"]中设置

func SetPath(params *map[string]map[string]string){
  var valueStr string = (*params)["params"]["path"] + "?"
  var reslist = make([]string, len((*params)["HeaderValue"]))
    var index int = 0
    for key, _ := range (*params)["HeaderValue"] {
        reslist[index] = key
        index++
    }
    sort.Strings(reslist)
    for _, value := range reslist {
        if "" == (*params)["HeaderValue"][value] {
            valueStr += value + "&"
        } else {
            valueStr += value + "=" + url.QueryEscape((*params)["HeaderValue"][value]) + "&"
        }
    }
    (*params)["params"]["path_url"] = valueStr[:len(valueStr)-1]
}

/拼接url 参数值中有文中的情况下参与加签不进行转码。
//参与url拼接的数据在api["HeaderValue"]中设置
func SetPathValue(params *map[string]map[string]string) {
    var valueStr string = (*params)["params"]["path"] + "?"
    var reslist = make([]string, len((*params)["HeaderValue"]))
    var index int = 0
    for key, _ := range (*params)["HeaderValue"] {
        reslist[index] = key
        index++
    }
    sort.Strings(reslist)
    for _, value := range reslist {
        if "" == (*params)["HeaderValue"][value] {
            valueStr += value + "&"
        } else {
            valueStr += value + "=" + (*params)["HeaderValue"][value] + "&"
        }
    }
    (*params)["params"]["path_value"] = valueStr[:len(valueStr)-1]
}
//拼接参与加签的path的字符串，参与的加签的字段在api["HeaderData"]中设置
func SetHeaderData(params *map[string]map[string]string) {
    for key, value := range (*params)["HeaderData"] {
        (*params)["params"]["header_data"] += key + ":" + value + (*params)["params"]["LF"]
    }
}

//对参与加签的字符串进行编码和加密
func GetSign(params *map[string]map[string]string){
  key := []byte((*params)["params"]["secret"])
  h := hmac.New(sha256.New, key)
  GetString(params)
  h.Write([]byte((*params)["params"]["X-Ca-Signature"]))
  (*params)["RequestHeader"]["X-Ca-Signature"] = base64.StdEncoding.EncodeToString(h.Sum(nil))
}


//拼接参与加签的字符串，根据阿里云提供的文档进行拼接
func GetString(params *map[string]map[string]string) {
    var signStr string = ""
    signStr += (*params)["params"]["httpMethod"] + (*params)["params"]["LF"]
    signStr += (*params)["params"]["LF"]
    signStr += (*params)["params"]["LF"]
    signStr += (*params)["params"]["LF"]
    signStr += (*params)["params"]["LF"]
    SetHeaderData(params)
    signStr += (*params)["params"]["header_data"]
    SetPathValue(params)
    SetPath(params)
    signStr += (*params)["params"]["path_value"]
    (*params)["params"]["X-Ca-Signature"] = signStr
    (*params)["params"]["url"] = (*params)["params"]["uri"] + (*params)["params"]["path_url"]
}
