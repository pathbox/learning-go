package account


type Account interface{
	Query(id string) int
	Update(id string, value int)
}

type AccountImpl struct {
	Id string
	Name string
	Value int
}

func (a *AccountImpl) Query(_ string) int {
    fmt.Println("AccountImpl.Query")
    return 100
}

func (a *AccountImpl) Update(_ string, _ int) {
    fmt.Println("AccountImpl.Update")
}

var New = func(id, name string, value int) Account {
    return &AccountImpl{id, name, value}
}