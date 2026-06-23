module github.com/pathbox/learning-go

go 1.14

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/PuerkitoBio/boom v0.0.0-20140219125548-fecdef1c97ca // indirect
	github.com/RoaringBitmap/roaring v0.4.21
	github.com/Shopify/sarama v1.26.4
	github.com/asim/go-micro/v3 v3.5.0
	github.com/astaxie/beego v1.10.1
	github.com/atotto/clipboard v0.1.4
	github.com/aws/aws-sdk-go v1.30.19
	github.com/benmanns/goworker v0.1.3
	github.com/bitly/go-simplejson v0.5.0
	github.com/bits-and-blooms/bitset v1.2.1
	github.com/boltdb/bolt v1.3.1
	github.com/boombuler/barcode v1.0.1-0.20190219062509-6c824513bacc
	github.com/bramvdbogaerde/go-scp v0.0.0-20200119201711-987556b8bdd7
	github.com/buger/jsonparser v0.0.0-20181115193947-bf1c66bbce23
	github.com/bwmarrin/snowflake v0.3.0
	github.com/c3mb0/go-do-work v0.0.0-20160309135746-a33fd02143e1
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575
	github.com/cirocosta/gupload v0.0.0-20180103143842-e6d8fa4fdf4c
	github.com/coreos/etcd v3.3.20+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/cornelk/hashmap v1.0.1
	github.com/dchest/siphash v1.2.2
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/djimenez/iconv-go v0.0.0-20160305225143-8960e66bd3da
	github.com/dlclark/regexp2 v1.2.0
	github.com/docker/go-connections v0.4.0
	github.com/dovejb/quicktag v0.0.0-20190829080553-340537080f34
	github.com/eapache/channels v1.1.0
	github.com/emirpasic/gods v1.12.0
	github.com/facebookgo/clock v0.0.0-20150410010913-600d898af40a // indirect
	github.com/facebookgo/grace v0.0.0-20180706040059-75cf19382434
	github.com/facebookgo/httpdown v0.0.0-20180706035922-5979d39b15c2
	github.com/facebookgo/stats v0.0.0-20151006221625-1b76add642e4 // indirect
	github.com/fatih/color v1.7.0
	github.com/fogleman/gg v1.3.0
	github.com/fsnotify/fsnotify v1.4.9
	github.com/fvbock/endless v0.0.0-20170109170031-447134032cb6
	github.com/garyburd/redigo v1.6.0
	github.com/go-kit/kit v0.9.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gobwas/glob v0.2.3
	github.com/gobwas/pool v0.2.1
	github.com/gobwas/ws v1.0.3
	github.com/gocraft/web v0.0.0-20190207150652-9707327fb69b
	github.com/gogo/protobuf v1.3.2
	github.com/goji/httpauth v0.0.0-20160601135302-2da839ab0f4d
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.4.2
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/uuid v1.1.2 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20190910122728-9d188e94fb99
	github.com/gorilla/context v1.1.1
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/mux v1.7.3
	github.com/gorilla/websocket v1.4.0
	github.com/grooveshark/golib v0.1.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/grpc-ecosystem/grpc-gateway v1.8.5
	github.com/hashicorp/golang-lru v0.5.3
	github.com/hashicorp/raft v1.1.2
	github.com/hpcloud/tail v1.0.0
	github.com/huandu/skiplist v0.0.0-20191129113331-b90e16040d86
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/jmoiron/sqlx v1.2.0
	github.com/json-iterator/go v1.1.9
	github.com/julienschmidt/httprouter v1.2.0
	github.com/justinas/alice v0.0.0-20171023064455-03f45bd4b7da
	github.com/justinas/nosurf v0.0.0-20190416172904-05988550ea18
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/kinwyb/go v0.0.0-20201029032031-48239eb7299c
	github.com/klauspost/reedsolomon v1.9.9 // indirect
	github.com/kr/pretty v0.2.0
	github.com/lestrrat-go/strftime v1.0.1
	github.com/lib/pq v1.2.0
	github.com/mattn/go-colorable v0.0.9
	github.com/mattn/go-runewidth v0.0.7
	github.com/mediocregopher/okq-go v0.0.0-20160211201133-048e319dd5ee
	github.com/mediocregopher/radix.v2 v0.0.0-20181115013041-b67df6e626f9
	github.com/micro/cli/v2 v2.1.2
	github.com/mitchellh/go-homedir v1.1.0
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d
	github.com/olekukonko/tablewriter v0.0.4
	github.com/opentracing/opentracing-go v1.1.0
	github.com/ory/ladon v1.2.0
	github.com/oxtoacart/bpool v0.0.0-20190530202638-03653db5a59c
	github.com/pascaldekloe/redis v1.3.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pborman/uuid v1.2.0
	github.com/pkg/browser v0.0.0-20210621091255-c198bc921a84
	github.com/pkg/errors v0.9.1
	github.com/pkg/profile v1.4.0
	github.com/pquerna/otp v1.2.0
	github.com/prometheus/client_golang v1.1.0
	github.com/rakyll/pb v0.0.0-20160123035540-8d46b8b097ef // indirect
	github.com/robertkrimen/otto v0.0.0-20180617131154-15f95af6e78d
	github.com/rs/xid v1.2.1
	github.com/rs/zerolog v1.18.0
	github.com/satori/go.uuid v1.2.0
	github.com/shirou/gopsutil v2.18.12+incompatible
	github.com/sirupsen/logrus v1.6.0
	github.com/sjwhitworth/golearn v0.0.0-20201127221938-294d65fca392
	github.com/smartystreets/goconvey v1.6.4
	github.com/spaolacci/murmur3 v1.1.0
	github.com/spf13/cast v1.3.1
	github.com/spf13/pflag v1.0.2
	github.com/stretchr/testify v1.6.1
	github.com/stvp/tempredis v0.0.0-20181119212430-b82af8480203
	github.com/templexxx/cpufeat v0.0.0-20180724012125-cef66df7f161 // indirect
	github.com/templexxx/xor v0.0.0-20191217153810-f85b25db303b // indirect
	github.com/throttled/throttled v2.2.4+incompatible
	github.com/tidwall/gjson v1.5.0
	github.com/tidwall/wal v0.1.1
	github.com/tikv/minitrace-go v0.0.0-20210119063709-5194f6ab6fd7
	github.com/tinylib/msgp v1.1.5
	github.com/tjfoc/gmsm v1.3.2 // indirect
	github.com/urfave/cli v1.22.1
	github.com/vmihailenco/msgpack v4.0.0+incompatible
	github.com/willf/bitset v1.1.10
	github.com/xtaci/kcp-go v5.4.20+incompatible
	github.com/yangwenmai/ratelimit v0.0.0-20180104140304-44221c2292e1
	github.com/youtube/vitess v2.1.1+incompatible // indirect
	github.com/zenazn/goji v0.9.0
	go.etcd.io/etcd v3.3.22+incompatible
	go.uber.org/zap v1.14.1
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	golang.org/x/exp v0.0.0-20200331195152-e8c3332aa8e5
	golang.org/x/net v0.0.0-20201021035429-f5854403a974
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a
	golang.org/x/sys v0.0.0-20210319071255-635bc2c9138d
	golang.org/x/text v0.3.5
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	golang.org/x/tools v0.0.0-20210106214847-113979e3529a
	gonum.org/v1/gonum v0.9.1 // indirect
	google.golang.org/appengine v1.6.6
	google.golang.org/grpc v1.29.1
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/throttled/throttled.v1 v1.0.0
	gopkg.in/yaml.v2 v2.3.0
	sigs.k8s.io/yaml v1.2.0 // indirect
	sourcegraph.com/sourcegraph/appdash v0.0.0-20190731080439-ebfcffb1b5c0
)
