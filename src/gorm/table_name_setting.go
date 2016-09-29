type User struct {} // default table name is `users`

// set User's table name to be `profiles
func (User) TableName() string {
  return "profiles"
}

func (u User) TableName() string {
    if u.Role == "admin" {
        return "admin_users"
    } else {
        return "users"
    }
}

// Disable table name's pluralization globally
db.SingularTable(true) // if set this to true,
//`User`'s default table name will be `user`, table name setted with `TableName` won't be affected
