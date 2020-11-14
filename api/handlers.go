package api

import (
	"net/http"

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

func getSession(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func authenticate(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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
		handler(w, r, p)
	}
}
