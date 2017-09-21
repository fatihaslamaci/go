package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

var templatesNavbar = template.Must(template.ParseFiles(
	"./templates/basenavbar.html",
	"./templates/head.html",
	"./templates/basefooter.html"))

var templatesBlank = template.Must(template.ParseFiles(
	"./templates/baseblank.html",
	"./templates/head.html",
	"./templates/basefooter.html"))

type Context struct {
	Title string
}

type CariKart struct {
	Message string
	Unvan   string
	Kod     string
}

var CariKart1 CariKart

func render(w http.ResponseWriter, tmpl string, js_tmpl string, context interface{}, template *template.Template) {
	tmpl_list := []string{fmt.Sprintf("views/%s.html", tmpl),
		fmt.Sprintf("templates/%s.html", js_tmpl),
	}

	t, err := template.Clone()
	if err != nil {
		log.Print("template clone error: ", err)
	}

	t, err = t.ParseFiles(tmpl_list...)
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.Execute(w, context)
	if err != nil {
		log.Print("template executing error: ", err)

	}
}

func renderBlank(w http.ResponseWriter, tmpl string, js_tmpl string, context interface{}) {
	render(w, tmpl, js_tmpl, context, templatesBlank)
}

func renderNavbar(w http.ResponseWriter, tmpl string, js_tmpl string, context interface{}) {
	render(w, tmpl, js_tmpl, context, templatesNavbar)
}

func blankHandler(w http.ResponseWriter, req *http.Request) {
	context := Context{Title: "Blak Page"}

	renderNavbar(w, "blank", "blank_js", context)
}

func carikartHandler(response http.ResponseWriter, request *http.Request) {

	renderNavbar(response, "carikart", "blank_js", CariKart1)
}

func carikartHandlerPost(response http.ResponseWriter, request *http.Request) {

	CariKart2 := CariKart1
	CariKart2.Unvan = request.FormValue("unvan")
	CariKart2.Kod = request.FormValue("kod")

	if len(CariKart2.Kod) > 0 {
		CariKart1 = CariKart2
	} else {
		CariKart2.Message = "Lütfen Zorunlu alanları giriniz"
	}
	renderNavbar(response, "carikart", "blank_js", CariKart2)

}

func indexPageHandler(response http.ResponseWriter, request *http.Request) {
	context := Context{Title: "Login Page"}
	renderBlank(response, "login", "blank_js", context)
}

func internalPageHandler(response http.ResponseWriter, request *http.Request) {
	userName := getUserName(request)
	context := Context{Title: "Home Page"}
	if userName != "" {
		renderNavbar(response, "main", "blank_js", context)
	} else {
		http.Redirect(response, request, "/", 302)
	}

}

// logout handler

func logoutHandler(response http.ResponseWriter, request *http.Request) {
	clearSession(response)
	http.Redirect(response, request, "/", 302)
}

// server main method

var router = mux.NewRouter()

func addStaticDir(s string) {
	http.Handle("/"+s+"/", http.StripPrefix("/"+s, http.FileServer(http.Dir("./statics/"+s))))
}

func addStaticDirAll() {
	addStaticDir("bower_components")
	addStaticDir("dist")
	addStaticDir("js")
	addStaticDir("less")
}

func main() {

	CariKart1.Kod = "ilk kod"
	CariKart1.Unvan = "ilk ünvan"

	router.HandleFunc("/", indexPageHandler)
	router.HandleFunc("/loginpage", loginPageHandler)
	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/logout", logoutHandler)
	router.HandleFunc("/main", internalPageHandler)

	router.HandleFunc("/blank", blankHandler)

	router.HandleFunc("/carikart", carikartHandler)
	router.HandleFunc("/carikartkaydet", carikartHandlerPost).Methods("POST")

	http.Handle("/", router)

	addStaticDirAll()

	http.ListenAndServe(":8000", nil)
}
