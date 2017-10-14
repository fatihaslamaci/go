package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var templatesNavbar = template.Must(template.ParseFiles(
	"./templates/basenavbar.html",
	"./templates/head.html",
	"./templates/topnav.html",
	"./templates/main_container.html",
	"./templates/sidebarmenu.html",
	"./templates/basefooter.html",
	"./templates/carikart.tmpl.html",
"./templates/carilist.tmpl.html"))




var templatesLogin = template.Must(template.ParseFiles(
	"./templates/logintamplate.html"))

type Context struct {
	Title   string
	Message string
	Data    interface{}

}



type CariKart struct {
	Unvan string
	Kod   string
	Telefon string
	Adres string
	Email string
}

var cariKartlar [1000]CariKart

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

func render2(w http.ResponseWriter, tmpl string, context interface{}) {
	err := templatesNavbar.ExecuteTemplate(w, tmpl, context)
	if err != nil {
		//log.Print("template executing error: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderLogin(w http.ResponseWriter, context interface{}) {
	err := templatesLogin.ExecuteTemplate(w, "logintamplate.html", context)
	if err != nil {
		//log.Print("template executing error: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderNavbar(w http.ResponseWriter, tmpl string, js_tmpl string, context interface{}) {
	render(w, tmpl, js_tmpl, context, templatesNavbar)
}

func blankHandler(w http.ResponseWriter, req *http.Request) {
	context := Context{Title: "Blank Page"}

	renderNavbar(w, "blank", "blank_js", context)
}

func carikartHandler(response http.ResponseWriter, request *http.Request) {

	context := Context{Title: "Cari Kart Tanıtımı", Data: cariKartlar[0]}

	render2(response, "carikart.tmpl.html", context)
}

func carilistHandler(response http.ResponseWriter, request *http.Request) {

	context := Context{Title: "Cari Kart Listesi", Data: cariKartlar}
	render2(response, "carilist.tmpl.html", context)

}

func formadvancedHandler(response http.ResponseWriter, request *http.Request) {

	context := Context{Title: "Cari Kart Listesi", Data: cariKartlar}

	renderNavbar(response, "form_advanced", "blank_js", context)
}

func carikartHandlerPost(response http.ResponseWriter, request *http.Request) {

	CariKart2 := cariKartlar[0]
	CariKart2.Unvan = request.FormValue("unvan")
	CariKart2.Kod = request.FormValue("kod")
	context := Context{Title: "Blak Page"}

	if len(CariKart2.Kod) > 0 {
		cariKartlar[0] = CariKart2
	} else {
		context.Message = "Lütfen Zorunlu alanları giriniz"
	}

	context.Data = CariKart2

	renderNavbar(response, "carikart", "blank_js", context)

}

func indexPageHandler(response http.ResponseWriter, request *http.Request) {
	context := Context{Title: "Login Page"}
	renderLogin(response, context)
}

func internalPageHandler(response http.ResponseWriter, request *http.Request) {
	context := Context{Title: "Home Page"}
	renderNavbar(response, "main", "blank_js", context)
}

func makeLoginHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userName := getUserName(r)
		fmt.Println("-" + userName + "-")
		if userName == "" {
			http.Redirect(w, r, "/", 302)
			return
		}
		fn(w, r)
	}
}

// server main method

var router = mux.NewRouter()

func addStaticDir(s string) {
	http.Handle("/"+s+"/", http.StripPrefix("/"+s, http.FileServer(http.Dir("./statics/"+s))))
}

func addStaticDirAll() {
	addStaticDir("build")
	addStaticDir("images")
	addStaticDir("production")
	addStaticDir("src")
	addStaticDir("vendors")

}

func main() {

	for i := 0; i < len(cariKartlar); i++ {
		// Display integer.
		cariKartlar[i] = CariKart{
			Kod: RandStringBytesMaskImprSrc(10),
			Unvan: RandStringBytesMaskImprSrc(50) ,
			Email:RandStringBytesMaskImprSrc(15)+ "@mail.com",
			Telefon:"0 536 555 12 12",
			Adres: RandStringBytesMaskImprSrc(20)+ " mahallesi, bilmem ne sk. no:1 d:"+strconv.Itoa(len(cariKartlar)),
		}
	}







	http.HandleFunc("/", indexPageHandler)
	http.HandleFunc("/login.html", loginHandler)
	http.HandleFunc("/login", loginHandlerPost)
	http.HandleFunc("/logout.html", logoutHandler)
	http.HandleFunc("/main", makeLoginHandler(internalPageHandler))
	http.HandleFunc("/blank", blankHandler)
	http.HandleFunc("/carikart.html", carikartHandler)
	http.HandleFunc("/carikartkaydet", carikartHandlerPost)
	http.HandleFunc("/carilist.html", carilistHandler)
	http.HandleFunc("/form_advanced.html", formadvancedHandler)



	addStaticDirAll();
	http.ListenAndServe(":8000", nil)

	/*
	router.HandleFunc("/", indexPageHandler)
	router.HandleFunc("/login.html", loginHandler)
	router.HandleFunc("/login", loginHandlerPost).Methods("POST")
	router.HandleFunc("/logout.html", logoutHandler)
	router.HandleFunc("/main", makeLoginHandler(internalPageHandler))

	router.HandleFunc("/blank", blankHandler)

	router.HandleFunc("/carikart.html", carikartHandler)
	router.HandleFunc("/carikartkaydet", carikartHandlerPost).Methods("POST")

	router.HandleFunc("/carilist.html", carilistHandler)

	router.HandleFunc("/form_advanced.html", formadvancedHandler)

	http.Handle("/", router)

	addStaticDirAll()

	http.ListenAndServe(":8000", nil)
	*/
}
