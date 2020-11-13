package api

import (
	"encoding/json"
	"net/http"

	"github.com/hawwwdi/Twitter_clone/db"
	"github.com/hawwwdi/Twitter_clone/user"
	"github.com/julienschmidt/httprouter"
)

func NewRouter() *httprouter.Router {
	mux := httprouter.New()
	mux.POST("/signUp", signUp)
	mux.POST("/logIn", logIn)
	return mux
}

func signUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	username, password := r.FormValue("username"), r.FormValue("password")
	if username == "" || password == "" {
		http.Error(w, "invalid username or password", http.StatusBadRequest)
		return
	}
	newUser := user.NewUser(username, password)
	_, err := db.RegisterUser(newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := json.Marshal(newUser)
	if err != nil {
		panic(err)
	}
	w.Write(res)
}

func logIn(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	username, password := r.FormValue("username"), r.FormValue("password")
	if username == "" || password == "" {
		http.Error(w, "invalid username or password", http.StatusBadRequest)
		return
	}
	uuid, err := db.LogIn(username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cookie := &http.Cookie{
		Name:   "session",
		Value:  uuid,
		MaxAge: 60 * 60,
	}
	http.SetCookie(w, cookie)
	//http.Redirect(w, r, "/home", http.StatusSeeOther)
	w.Write([]byte("logged in"))
}
