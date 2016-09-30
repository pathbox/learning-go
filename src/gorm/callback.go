func (u *User) BeforeUpdate() (err error) {
    if u.readonly() {
        err = errors.New("read only user")
    }
    return
}

// Rollback the insertion if user's id greater than 1000
func (u *User) AfterCreate() (err error) {
    if (u.Id > 1000) {
        err = errors.New("user id is already greater than 1000")
    }
    return
}
Save/Delete operations in gorm are running in transactions,
so changes made in that transaction are not visible
unless it is commited. If you want to use those changes in your
callbacks, you need to run your SQL in the same transaction.
So you need to pass current transaction to callbacks like this:

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
    tx.Model(u).Update("role", "admin")
    return
}
func (u *User) AfterCreate(scope *gorm.Scope) (err error) {
  scope.DB().Model(u).Update("role", "admin")
    return
}
