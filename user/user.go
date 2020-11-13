package user

type user struct {
	id        string
	username  string
	password  string
	Following []int
	Followers []int
}

func (u *user) setId(id string) {
	u.id = id
}

func (u *user) info() (string, string, string) {
	return u.id, u.username, u.password
}
