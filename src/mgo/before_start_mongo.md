#### db 初始化

测试用本地 db: (./test 单元测试必须)
mongouri: "app_test:app123456@127.0.0.1:27017/app_test"

```
mongo  # db管理员权限
use app_test
db.createUser(
 {
   user: "app_test",
   pwd: "app123456",
   roles:
     [
       { role: "readWrite", db: "app_test" }
     ]
 }
)

db.config.insert({"AppName": "app"})
db.config.insert({"Env": "test"})

```

开发用本地 db: (本地开发必须)
mongouri: "app_dev:app123456@127.0.0.1:27017/app_dev"

```
mongo  # db管理员权限
use app_dev
db.createUser(
 {
   user: "app_dev",
   pwd: "app123456",
   roles:
     [
       { role: "readWrite", db: "app_dev" }
     ]
 }
)
db.config.insert({"AppName": "app"})
db.config.insert({"Env": "development"})
```