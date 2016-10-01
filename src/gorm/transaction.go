func CreateAnimals(db *gorm.DB) err{
  tx := db.Begin()

  if err := tx.Create(&Animal{Name: "Giraffe"}).Error; err != nil {
    tx.Rollback()
    return err
  }

  if err := tx.Create(&Animal{Name: "Lion"}).Error; err != nil {
     tx.Rollback()
     return err
  }

  tx.Commit()
  return nil
}
