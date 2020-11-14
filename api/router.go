package api

import (
	"errors"

	"github.com/julienschmidt/httprouter"
)

var hub *Hub

type Hub struct {
	Router *httprouter.Router
	db     DB
}

func InitHub(db DB) {
	hub = &Hub{
		Router: newRouter(),
		db:     db,
	}
}

func GetHub() (*Hub, error) {
	if hub == nil {
		return nil, errors.New("init hub first")
	}
	return hub, nil
}

func newRouter() *httprouter.Router {
	mux := httprouter.New()
	mux.POST("/signUp", signUp)
	mux.POST("/logIn", logIn)
	mux.GET("/logOut", authenticate(logOut))
	mux.GET("/follow/:user", authenticate(follow))
	mux.POST("/post", authenticate(post))
	mux.GET("/home", authenticate(showUserPosts))
	mux.GET("/timeline", authenticate(showTimeLinePosts))
	return mux
}
