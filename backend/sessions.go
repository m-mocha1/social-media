package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gofrs/uuid/v5"
)

var (
	sessionsMap  = map[string]string{}
	userSession  = map[string]string{}
	emailSession = map[string]string{}
)

// make a session id ---------------------------------------------
func makeSessionId() string {
	id, err := uuid.NewV4()
	if err != nil {
		fmt.Print("failed to generate UUID:", err)
		os.Exit(0)
	}
	return id.String()
}

// creat a new session id ---------------------------------------
func setSession(w http.ResponseWriter, userName string) {
	old, already := userSession[userName]
	old2, already2 := emailSession[userName]
	if already {
		delete(sessionsMap, old)
	}
	if already2 {
		delete(emailSession, old2)
	}
	sId := makeSessionId()
	sessionsMap[sId] = userName
	userSession[userName] = sId
	http.SetCookie(w, &http.Cookie{
		Name:     "sID",
		Value:    sId,
		Secure:   true,
		HttpOnly: true,
	})
}

func setEmailSession(w http.ResponseWriter, email string) {
	old, already := userSession[email]
	if already {
		delete(sessionsMap, old)
	}
	sId := makeSessionId()
	sessionsMap[sId] = email
	emailSession[email] = sId
	http.SetCookie(w, &http.Cookie{
		Name:     "sID",
		Value:    sId,
		Secure:   true,
		HttpOnly: true,
	})
}

// gets a session id from the request---------------------------
func getSession(r *http.Request) (string, bool) {
	cookie, err := r.Cookie("sID")
	if err != nil {
		return "", false
	}
	username, exists := sessionsMap[cookie.Value]
	return username, exists
}

func getEmailSession(r *http.Request) (string, bool) {
	cookie, err := r.Cookie("sID")
	if err != nil {
		return "", false
	}
	username, exists := emailSession[cookie.Value]
	return username, exists
}

