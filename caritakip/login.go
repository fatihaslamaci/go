package main

import (
	"github.com/gorilla/securecookie"
	"net/http"
	"fmt"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))


func getUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["Username"]
		}
	}
	return userName
}

func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"Username": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

// login handler

func loginHandler(response http.ResponseWriter, request *http.Request) {
	clearSession(response);
	context := Context{Title: "Login Page"}
	renderLogin(response, context)
}


func loginHandlerPost(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("Username")
	pass := request.FormValue("Password")

	fmt.Println("name")
	fmt.Println(pass)

	redirectTarget := "/"
	if (name == "Fatih") && (pass == "1") {
		// .. check credentials ..
		setSession(name, response)
		redirectTarget = "/main"
		fmt.Println(name)

	}
	http.Redirect(response, request, redirectTarget, 302)
}

// logout handler
func logoutHandler(response http.ResponseWriter, request *http.Request) {
	clearSession(response)
	http.Redirect(response, request, "/", 302)
}

