db.Set("gorm:save_associations", false).Create(&user)

type Animal struct {
    ID   int64
    Name string `gorm:"default:'galeone'"`
    Age  int64
}

var animal = Animal{Age: 99, Name: ""}

db.Create(&animal)
// INSERT INTO animals("age") values('99');
// SELECT name from animals WHERE ID=111; // the returning primary key is 111
// animal.Name => 'galeone'

// If you want to set primary field's value in BeforeCreate callback, you could use scope.SetColumn, for example:
func (user *User) BeforeCreate(scope *gorm.Scope) error {
  scope.SetColumn("ID", uuid.New())
  return nil
}

db.Set("gorm:insert_option", "ON CONFLICT").Create(&product)

// INSERT INTO products (name, code) VALUES ("name", "code") ON CONFLICT;
