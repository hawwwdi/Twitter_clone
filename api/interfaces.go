package api

type DB interface {
	RegisterUser(username, password string) error
	LogIn(username, password string) (string, error)
	LogOut(session string) error
	Follow(follower, followed string) error
	Post(body, owner string) error
	ShowUserPosts(id string, start, count int) (map[string][]string, error)
	GetSessionUserID(session string) (string, error)
}
