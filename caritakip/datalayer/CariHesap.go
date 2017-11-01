package datalayer

import (
	"database/sql"
	"github.com/fatihaslamaci/go/caritakip/entity"
)


func StoreItem(db *sql.DB, items []entity.CariKart) {
	/*
	sql_additem := `
	INSERT OR REPLACE INTO carihe (
	) values(?, ?, ?, CURRENT_TIMESTAMP)
	`
*/
	stmt, err := db.Prepare("INSERT OR REPLACE INTO carihesaplar(id,kod, unvan, telefon, adres, email, created ) VALUES (?,?,?,?,?,?,CURRENT_TIMESTAMP)")
	if err != nil { panic(err) }
	defer stmt.Close()

	for _, item := range items {
		_, err2 := stmt.Exec(item.Id, item.Kod, item.Unvan,item.Telefon,item.Adres,item.Email)
		if err2 != nil { panic(err2) }
	}
}

func UpdateCariKart(db *sql.DB, item entity.CariKart) {
	/*
	sql_additem := `
	INSERT OR REPLACE INTO carihe (
	) values(?, ?, ?, CURRENT_TIMESTAMP)
	`
*/
	stmt, err := db.Prepare("Update carihesaplar set kod=?, unvan=?, telefon=?, adres=?, email=? WHERE id=?")
	if err != nil { panic(err) }
	defer stmt.Close()
	_, err2 := stmt.Exec(item.Kod, item.Unvan,item.Telefon,item.Adres,item.Email,item.Id)
	if err2 != nil { panic(err2) }

}



func ReadItem(db *sql.DB, SonId int) ([]entity.CariKart,int){
	rows, err := db.Query(`
Select id,kod, unvan, telefon, adres, email

  ,ifnull((select sum(ifnull(tutar,0)) from carihareket where carihesapid=carihesaplar.id and ba=1),0) as Borc
  ,ifnull((select sum(ifnull(tutar,0)) from carihareket where carihesapid=carihesaplar.id and ba=0),0) as Alacak

from carihesaplar
where id>?
order by id LIMIT 30
`,SonId)

	if err != nil { panic(err) }
	defer rows.Close()

	var result []entity.CariKart
	for rows.Next() {
		item := entity.CariKart{}
		err2 := rows.Scan(&item.Id, &item.Kod, &item.Unvan,&item.Telefon,&item.Adres,&item.Email,&item.Borc,&item.Alacak)
		if err2 != nil { panic(err2) }
		result = append(result, item)
	}

	var son int

	if len(result)>0{
		son=result[len(result)-1].Id
	}

	return result,son
}

func ReadItemId(db *sql.DB, id int) entity.CariKart {
	item := entity.CariKart{}
	row := db.QueryRow("Select id,kod, unvan, telefon, adres, email from carihesaplar where id=?",id)
	err :=row.Scan(&item.Id, &item.Kod, &item.Unvan,&item.Telefon,&item.Adres,&item.Email)
	if err != nil { panic(err) }

	return item
}
