package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Context struct {
	Title   string
	Message string
	UserName string

	KayitId int

	Data    interface{}
}


type CariKart struct {
	Id int
	Unvan   string
	Kod     string
	Telefon string
	Adres   string
	Email   string
}



func render2(w http.ResponseWriter,r *http.Request , tmpl string, context Context) {
	files := []string{
		"./templates/base.html", "./templates/" + tmpl + ".html",
	}


	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	context.UserName=getUserName(r)
	err = ts.ExecuteTemplate(w, "base", context)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

}



func internalPageHandler(writer http.ResponseWriter, request *http.Request) {
	context := Context{Title: "Cari Kart Tanıtımı", Data: ReadItemId(db,0)}
	render2(writer,request, "auth/dashboard", context)
}


func carikartHandler(response http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	id:=r.FormValue("id")
	var i int
	i,_ = strconv.Atoi(id)
	context := Context{Title: "Cari Kart Tanıtımı", Data: ReadItemId(db,i)}
	render2(response,r, "auth/carikart", context)
}


func carikartlarHandler(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	id,_:=strconv.Atoi(request.FormValue("id"))

	fData, son := ReadItem(db, id)
	context := Context{Data: fData, KayitId: son}
	render2(response, request, "auth/carikartlar", context)

}

func carikartHandlerPost(response http.ResponseWriter, request *http.Request) {

	id,_:= strconv.Atoi(request.FormValue("id"))

	CariKart2 := ReadItemId(db,id)
	CariKart2.Unvan = request.FormValue("unvan")
	CariKart2.Kod = request.FormValue("kod")
	CariKart2.Adres = request.FormValue("adres")
	CariKart2.Telefon = request.FormValue("telefon")
	CariKart2.Email = request.FormValue("email")

	context := Context{}


	if len(CariKart2.Kod) > 0 {
		UpdateCariKart(db,CariKart2)
		context.Message = "Kayıt yapıldı"
	} else {
		context.Message = "Lütfen Zorunlu alanları giriniz"
	}

	context.Data= CariKart2
	render2(response,request, "auth/carikart", context)

}

func indexPageHandler(w http.ResponseWriter, r *http.Request) {

	render2(w, r,"index", Context{})

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



func addStaticDir(s string) {
	http.Handle("/"+s+"/", http.StripPrefix("/"+s, http.FileServer(http.Dir("./statics/"+s))))
}

func addStaticDirAll() {
	addStaticDir("css")
	addStaticDir("images")
	addStaticDir("production")
	addStaticDir("src")
	addStaticDir("vendors")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

var db *sql.DB

func main() {

	const dbpath = "./db/foo.sqlite"
	db = InitDB(dbpath)
	defer db.Close()
	CreateTable(db)

/*
	var cariKartlar []CariKart
	cariKartlar = make([]CariKart,2000000)
	for i := 0; i < len(cariKartlar); i++ {
		// Display integer.
		cariKartlar[i] = CariKart{
			Id:i,
			Kod:     RandStringBytesMaskImprSrc(10),
			Unvan:   RandStringBytesMaskImprSrc(50),
			Email:   RandStringBytesMaskImprSrc(10) + "@mail.com",
			Telefon: "0 536 555 12 12",
			Adres:   RandStringBytesMaskImprSrc(20) + " mahallesi, bilmem ne sk. no:1 d:" + strconv.Itoa(len(cariKartlar)),
		}
	}

	StoreItem(db,cariKartlar)
*/
	http.HandleFunc("/", makeLoginHandler(indexPageHandler))
	http.HandleFunc("/login.html", loginHandler)
	http.HandleFunc("/login", loginHandlerPost)
	http.HandleFunc("/logout.html", logoutHandler)
	http.HandleFunc("/auth/dashboard.html", makeLoginHandler(internalPageHandler))
	http.HandleFunc("/auth/carikart.html", carikartHandler)

	http.HandleFunc("/auth/carikartkaydet", carikartHandlerPost)

	http.HandleFunc("/auth/carikartlar.html", makeLoginHandler(carikartlarHandler))


	addStaticDirAll()
	http.ListenAndServe(":8000", nil)

}



