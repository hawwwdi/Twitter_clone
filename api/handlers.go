package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func signUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	username, password := r.FormValue("username"), r.FormValue("password")
	if username == "" || password == "" {
		http.Error(w, "invalid username or password", http.StatusBadRequest)
		return
	}
	err := hub.db.RegisterUser(username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	/* 	res, err := json.Marshal(newUser)
	   	if err != nil {
	   		panic(err)
	   	} */
	http.Redirect(w, r, "/logIn", http.StatusSeeOther)
}

func logIn(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	username, password := r.FormValue("username"), r.FormValue("password")
	if username == "" || password == "" {
		http.Error(w, "invalid username or password", http.StatusBadRequest)
		return
	}
	uuid, err := hub.db.LogIn(username, password)
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

func logOut(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session, _ := getSession(r)
	err := hub.db.LogOut(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func follow(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//todo follow user by username
	toFollow := ps.ByName("user")
	if toFollow == "" {
		http.Error(w, "user id not found", http.StatusBadRequest)
		return
	}
	session, _ := getSession(r)
	followerID, _ := hub.db.GetSessionUserID(session)
	err := hub.db.Follow(followerID, toFollow)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func post(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body := r.FormValue("body")
	if body == "" {
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}
	session, _ := getSession(r)
	owner, _ := hub.db.GetSessionUserID(session)
	err := hub.db.Post(body, owner)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func showUserPosts(w http.ResponseWriter, r *http.Request, sp httprouter.Params) {
	sstart, scount := r.FormValue("start"), r.FormValue("count")
	start, err := strconv.Atoi(sstart)
	count, err1 := strconv.Atoi(scount)
	if err != nil || err1 != nil {
		http.Error(w, "require start and count field", http.StatusBadRequest)
		return
	}
	session, _ := getSession(r)
	id, _ := hub.db.GetSessionUserID(session)
	posts, err := hub.db.ShowUserPosts(id, start, count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	res, _ := json.Marshal(posts)
	w.Header()["content-type"] = []string{"application/json"}
	w.Write(res)
}

func getSession(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func authenticate(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		session, err := getSession(r)
		if session == "" || err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		id, err := hub.db.GetSessionUserID(session)
		if id == "" || err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		handler(w, r, ps)
	}
}
