package api

import "github.com/hawwwdi/Twitter_clone/model"

type DB interface {
	RegisterUser(username, password string) (string, error)
	LogIn(username, password string) (string, error)
	LogOut(session string) error
	Follow(follower, followed string) error
	Post(post model.Post) (string, error)
	ShowTimeLinePosts(count int) (map[string]interface{}, error)
	ShowUserPosts(id string, start, count int) (map[string]interface{}, error)
	GetUser(id string) (map[string]string, error)
	GetSessionUserID(session string) (string, error)
}
