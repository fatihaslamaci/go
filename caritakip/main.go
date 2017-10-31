package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"strconv"

	"github.com/fatihaslamaci/go/caritakip/datalayer"
	"github.com/fatihaslamaci/go/caritakip/entity"
)

type Context struct {
	Title    string
	Message  string
	UserName string

	KayitId  string
	KayitId2 string

	Data interface{}
}

func render2(w http.ResponseWriter, r *http.Request, tmpl string, context Context) {
	files := []string{
		"./templates/base.html", "./templates/" + tmpl + ".html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	context.UserName = getUserName(r)
	err = ts.ExecuteTemplate(w, "base", context)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

}

func internalPageHandler(writer http.ResponseWriter, request *http.Request) {
	context := Context{Title: "Cari Kart Tanıtımı", Data: datalayer.ReadItemId(db, 0)}
	render2(writer, request, "auth/dashboard", context)
}

func indexPageHandler(w http.ResponseWriter, r *http.Request) {

	render2(w, r, "index", Context{})

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
	db = datalayer.InitDB(dbpath)
	defer db.Close()
	datalayer.CreateTable(db)

	var carihareketler []entity.CariHareket
	carihareketler = make([]entity.CariHareket, 200)
	for i := 0; i < len(carihareketler); i++ {
		// Display integer.
		carihareketler[i] = entity.CariHareket{
			Id:          i,
			CariHesapId: i,
			Ba:          true,
			Tutar:       (100.02),
		}

	}

	datalayer.InsertCariHareket(db, carihareketler)

	var cariKartlar []entity.CariKart
	cariKartlar = make([]entity.CariKart, 200)
	for i := 0; i < len(cariKartlar); i++ {
		// Display integer.
		cariKartlar[i] = entity.CariKart{
			Id:      i,
			Kod:     RandStringBytesMaskImprSrc(10),
			Unvan:   RandStringBytesMaskImprSrc(50),
			Email:   RandStringBytesMaskImprSrc(10) + "@mail.com",
			Telefon: "0 536 555 12 12",
			Adres:   RandStringBytesMaskImprSrc(20) + " mahallesi, bilmem ne sk. no:1 d:" + strconv.Itoa(len(cariKartlar)),
		}
	}

	datalayer.StoreItem(db, cariKartlar)

	http.HandleFunc("/", makeLoginHandler(indexPageHandler))
	http.HandleFunc("/login.html", loginHandler)
	http.HandleFunc("/login", loginHandlerPost)
	http.HandleFunc("/logout.html", logoutHandler)
	http.HandleFunc("/auth/dashboard.html", makeLoginHandler(internalPageHandler))
	http.HandleFunc("/auth/carikart.html", carikartHandler)

	http.HandleFunc("/auth/carikartkaydet", carikartHandlerPost)

	http.HandleFunc("/auth/carikartlar.html", carikartlarHandler)

	addStaticDirAll()
	http.ListenAndServe(":8000", nil)

	log.Println("server başladı")

}
