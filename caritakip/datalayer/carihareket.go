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

func ToplamRapor(db *sql.DB) (borc float64, alacak float64){
	var tborc float64
	var talacak float64
	row := db.QueryRow(`
Select
  ifnull((select sum(ifnull(tutar,0)) from carihareket where ba=1),0) as Borc
  ,ifnull((select sum(ifnull(tutar,0)) from carihareket where ba=0),0) as Alacak
`)
	err := row.Scan(&tborc, &talacak)
	if err != nil { panic(err) }

	return tborc,talacak
}




func Ekstre(db *sql.DB, Id int) ([]entity.CariHareket,int){
	rows, err := db.Query(`
Select id,carihesapid,ba,tutar
from carihareket
where CariHesapId=?

`,Id)

	if err != nil { panic(err) }
	defer rows.Close()

	var result []entity.CariHareket
	for rows.Next() {
		item := entity.CariHareket{}
		err2 := rows.Scan(&item.Id, &item.CariHesapId, &item.Ba,&item.Tutar)
		if err2 != nil { panic(err2) }
		result = append(result, item)
	}

	var son int

	if len(result)>0{
		son=result[len(result)-1].Id
	}

	return result,son
}
