package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type Context struct {
	Title   string
	Message string
	Data    interface{}
}

type CariKart struct {
	Unvan   string
	Kod     string
	Telefon string
	Adres   string
	Email   string
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
	files := []string{
		"./templates/base.html", "./templates/" + tmpl + ".html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.ExecuteTemplate(w, "base", context)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func carikartHandler(response http.ResponseWriter, request *http.Request) {

	context := Context{Title: "Cari Kart Tanıtımı", Data: cariKartlar[0]}

	render2(response, "carikart.tmpl.html", context)
}

func carilistHandler(response http.ResponseWriter, request *http.Request) {

	context := Context{Title: "Cari Kart Listesi", Data: cariKartlar}
	render2(response, "carilist.tmpl.html", context)

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

	//renderNavbar(response, "carikart", "blank_js", context)

}

func indexPageHandler(w http.ResponseWriter, r *http.Request) {

	render2(w, "index", nil)

}

func makeLoginHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if (len(r.URL.Path) > 6) && (r.URL.Path[0:6] == "/auth/") {
			userName := getUserName(r)
			fmt.Println("-" + userName + "-")
			if userName == "" {
				http.Redirect(w, r, "/login.html", 302)
				return
			}
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
			Kod:     RandStringBytesMaskImprSrc(10),
			Unvan:   RandStringBytesMaskImprSrc(50),
			Email:   RandStringBytesMaskImprSrc(15) + "@mail.com",
			Telefon: "0 536 555 12 12",
			Adres:   RandStringBytesMaskImprSrc(20) + " mahallesi, bilmem ne sk. no:1 d:" + strconv.Itoa(len(cariKartlar)),
		}
	}

	http.HandleFunc("/", makeLoginHandler(indexPageHandler))
	http.HandleFunc("/login.html", loginHandler)
	http.HandleFunc("/login", loginHandlerPost)
	http.HandleFunc("/logout.html", logoutHandler)
	http.HandleFunc("/auth/main.html", makeLoginHandler(internalPageHandler))
	http.HandleFunc("/carikart.html", carikartHandler)
	http.HandleFunc("/carikartkaydet", carikartHandlerPost)
	http.HandleFunc("/carilist.html", carilistHandler)

	addStaticDirAll()
	http.ListenAndServe(":8000", nil)

}
func internalPageHandler(writer http.ResponseWriter, request *http.Request) {
	context := Context{Title: "Cari Kart Tanıtımı", Data: cariKartlar[0]}
	render2(writer, "auth/main", context)
}
