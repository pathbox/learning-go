https://github.com/SebastiaanKlippert/go-wkhtmltopdf/blob/master/options.go#L220:6

TODO:

解析 options

调用 ServerGO 接口 获取 html string

创建 html 文件

调用 wkhtmltopdf html 文件, 创建 PDF文件

将PDF 文件 返回

close html file 文件描述符

如果有报错, 则返回错误信息


其他

sid: 每个服务使用方有一个sid, 12位长度的唯一字符串
一个sid 对应两个secret, 一个为convert接口的secret, 一个为 url 对应鉴权的secret
存在两个对应表中

访问需要鉴权url的鉴权算法: token = SHA1(sid+timestamp+secret)
单独为ServerGo做鉴权方式

默认禁止WEB页面执行JavaScript --disable-javascript


第二期

支持word

新增 post方法,接收html字符串数据,进行转换