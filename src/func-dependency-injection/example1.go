package injection

type User struct {
	Name string
	ID   string
}

type getUser func(id string) User

type selectUserByID func(id string) string

func newGetUser(selectUser selectUserByID) getUser {
	return func(id string) User {
		name := selectUser(id)
		return User{
			ID:   id,
			Name: name,
		}
	}
}

func newSelectUserByID() selectUserByID {
	return func(id string) string {
		return "Cary"
	}
}
