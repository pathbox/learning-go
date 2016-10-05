db.DB()

db.DB().Ping()
// 连接池
db.DB().SetMaxIdleConns(100)

//Set multiple fields as primary key to enable composite primary key

type Product struct {
    ID           string `gorm:"primary_key"`
    LanguageCode string `gorm:"primary_key"`
}

// Enable Logger, show detailed log
db.LogMode(true)

// Diable Logger, don't show any log
db.LogMode(false)

// Debug a single operation, show detailed log for this operation
db.Debug().Where("name = ?", "jinzhu").First(&User{})


//Refer GORM's default logger for how to customize it https://github.com/jinzhu/gorm/blob/master/logger.go

db.SetLogger(gorm.Logger{revel.TRACE})
db.SetLogger(log.New(os.Stdout, "\r\n", 0))
