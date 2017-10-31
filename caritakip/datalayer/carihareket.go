package datalayer

import (
	"database/sql"
	"github.com/fatihaslamaci/go/caritakip/entity"
)

func InsertCariHareket(db *sql.DB, items []entity.CariHareket) {

	stmt, err := db.Prepare("INSERT INTO carihareket(carihesapid,ba,tutar) VALUES (?,?,?)")
	if err != nil { panic(err) }
	defer stmt.Close()

	for _, item := range items {
		_, err2 := stmt.Exec(item.CariHesapId, item.Ba,item.Tutar)
		if err2 != nil { panic(err2) }
	}
}
