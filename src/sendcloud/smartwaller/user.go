package sendcloud

import "net/url"

const (
	SEND_CLOUD_USERINFO_GET_API_URL = "http://api.sendcloud.net/apiv2/userinfo/get"
)

////////////////////////////////////////////////////////////////////////////////
// GetUserInfo 查询当前用户的相关信息
func GetUserInfo() (bool, error, string) {
	params := url.Values{}
	return doRequest(SEND_CLOUD_USERINFO_GET_API_URL, params)
}
