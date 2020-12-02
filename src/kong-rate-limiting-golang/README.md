https://github.com/lampnick/kong-rate-limiting-golang

部署方式一：使用Docker(可直接clone仓库，执行make相关命令构建kong镜像）【如果clone了代码，可以在项目根目录使用make run-kong-konga-pg运行kong、postgres及konga全套环境】
拉镜像
docker pull lampnick/kong-rate-limiting-plugin-golang:latest
运行docker
docker run --rm --name kong-rate-limiting-plugin-golang \
    -e "KONG_LOG_LEVEL=info" \
    -e "KONG_NGINX_USER=root root" \
    -p 8000:8000 \
    -p 8443:8443 \
    -p 8001:8001 \
    -p 8444:8444 \
    lampnick/kong-rate-limiting-plugin-golang:latest
测试插件是否加载成功
curl http://localhost:8001/ |grep --color custom-rate-limiting
部署方式二：服务器编译部署（牛刀小试,只需简单几步即可体验本插件）
clone本项目到/etc/kong/
mkdir /etc/kong
cd /etc/kong
git clone https://github.com/lampnick/kong-rate-limiting-golang.git
修改kong配置文件
plugins = bundled,custom-rate-limiting
go_plugins_dir = /etc/kong/plugins
go_pluginserver_exe = /usr/local/bin/go-pluginserver
构建go-pluginserver
在go-pluginserver中执行go build github.com/Kong/go-pluginserver
会生成 go-pluginserver文件，复制到/usr/local/bin目录
编译go插件
go build -buildmode plugin custom-rate-limiting.go
将生成的.so文件放到go_plugins_dir(上面配置为/etc/kong/plugins)定义的目录中
cp custom-rate-limiting.so /etc/kong/plugins/
重启kong
kong prepare && kong reload



插件开发流程
定义一个结构体类型保存配置文件
用lua写的插件通过schema来指定怎样读取和验证来自数据库和Admin API中的配置数据。由于GO是静态类型语言，都需要用配置结构体定义
type MyConfig struct {
    Path   string //这里配置的会在konga添加插件时显示出来
    Reopen bool
}
公有属性将会被配置数据填充，如果希望在数据库中使用不同的名称，可以使用encoding/json加tag的方式
type MyConfig struct {
    Path   string `json:my_file_path`
    Reopen bool   `json:reopen`
}
使用New()创建一个实例
你的go插件必须定义一个名叫New的函数来创建这个类型的实例并返回一个interface{}类型
func New() interface{} {
    return &MyConfig{}
}
添加处理阶段方法
你可以在请求的生命周期的各个阶段实现自定义的逻辑。如在"access"阶段,定义一个名为Access的方法
func (conf *MyConfig) Access (kong *pdk.PDK) {
  ...
}
你可以实现自定义逻辑的阶段方法有如下几种
Certificate
Rewrite
Access
Preread
Log
编译go插件
go build -buildmode plugin  custom-rate-limiting.go
将生成的.so文件放到go_plugins_dir定义的目录中
cp custom-rate-limiting.so ../plugins/
重启kong（平滑重启）
kong prepare && kong reload