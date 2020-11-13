package user

type User struct {
	ID         string `json:"ID"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Auth       string `json:"auth"`
	Followings []string
	Followers  []string
}

func NewUser(username, password string) *User {
	return &User{
		Username:   username,
		Password:   password,
		Followings: make([]string, 0),
		Followers:  make([]string, 0),
	}
}

func (u *User) SetId(id string) {
	u.ID = id
}

func (u *User) Info() (string, string, string) {
	return u.ID, u.Username, u.Password
}
