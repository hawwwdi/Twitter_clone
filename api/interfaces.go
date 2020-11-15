package api

type DB interface {
	RegisterUser(username, password string) (string, error)
	LogIn(username, password string) (string, error)
	LogOut(session string) error
	Follow(follower, followed string) error
	Post(body, owner string) (string, error)
	ShowTimeLinePosts(count int) (map[string][]string, error)
	ShowUserPosts(id string, start, count int) (map[string][]string, error)
	GetUser(id string) (map[string]string, error)
	GetSessionUserID(session string) (string, error)
}
