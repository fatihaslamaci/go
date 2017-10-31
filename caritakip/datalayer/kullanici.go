package datalayer

import (
	"database/sql"
	"github.com/fatihaslamaci/go/caritakip/entity"
)

func ReadKullanici(db *sql.DB, id int) entity.Kullanici {
	item := entity.Kullanici{}
	row := db.QueryRow("Select id,ad,email,sifre from kullanicilar where id=?",id)
	err :=row.Scan(&item.Id, &item.Ad, &item.Email,&item.Sifre)
	if err != nil { panic(err) }
	return item
}

func Login(db *sql.DB, email string, sifre string) (entity.Kullanici) {
	item := entity.Kullanici{}
	row := db.QueryRow("Select id,ad,email,sifre from kullanicilar where email=? and sifre=?",email,sifre)
	err :=row.Scan(&item.Id, &item.Ad, &item.Email,&item.Sifre)

	if err != nil { println("Hatalı kullanıcı girişi")  }

	return item
}



