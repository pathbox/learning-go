if err := db.Where("name = ?", "jinzhu").First(&user).Error; err != nil {
  // error handling
}

db.First(&user).Limit(10).Find(&user).GetErrors()

db.Where("name = ?", "hello world").First(&user).RecordNotFound()

if db.Model(&user).Related(&credit_card).ReordNotFound(){

}
