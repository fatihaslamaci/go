package main

import (
	"net/http"
	"strconv"
	"github.com/fatihaslamaci/go/caritakip/datalayer"
	"strings"
)

func carikartlarHandler(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	id :=request.FormValue("id")

	if (id==""){
		id="0"
	}

	idlist := strings.Split(id, "-")

	Tip:=request.FormValue("T")

	idint:=0

	if Tip=="1"{
		idlist = idlist[:len(idlist)-1]
	}

	if len(idlist)>0 {
		idint, _ = strconv.Atoi(idlist[len(idlist)-1])
	}

	fData, son := datalayer.ReadItem(db, idint)

	idlist = append(idlist, strconv.Itoa(son))

	s:=""

	for i, eleman := range idlist {
		if i>0 {
			s +="-"+eleman
		}else{
			s=eleman
		}
	}

	s2:=""
	for i, eleman := range idlist[:len(idlist)-1] {
		if i>0 {
			s2 +="-"+eleman
		}else{
			s2=eleman
		}
	}

	context := Context{Data: fData, KayitId: s, KayitId2: s2}
	render2(response, request, "auth/carikartlar", context)
}


func carikartHandler(response http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id:=r.FormValue("id")
	var i int
	i,_ = strconv.Atoi(id)
	context := Context{Title: "Cari Kart Tanıtımı", Data: datalayer.ReadItemId(db,i)}
	render2(response,r, "auth/carikart", context)
}


func carikartHandlerPost(response http.ResponseWriter, request *http.Request) {

	request.ParseForm()
	id,_:= strconv.Atoi(request.FormValue("id"))

	CariKart2 := datalayer.ReadItemId(db,id)
	CariKart2.Unvan = request.FormValue("unvan")
	CariKart2.Kod = request.FormValue("kod")
	CariKart2.Adres = request.FormValue("adres")
	CariKart2.Telefon = request.FormValue("telefon")
	CariKart2.Email = request.FormValue("email")

	context := Context{}


	if len(CariKart2.Kod) > 0 {
		datalayer.UpdateCariKart(db,CariKart2)
		context.Message = "Kayıt yapıldı"
	} else {
		context.Message = "Lütfen Zorunlu alanları giriniz"
	}

	context.Data= CariKart2
	render2(response,request, "auth/carikart", context)

}

