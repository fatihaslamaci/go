package main

import "database/sql"

func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil { panic(err) }
	if db == nil { panic("db nil") }
	return db
}

func CreateTable(db *sql.DB) {
	// create table if not exists
	sql_table := `
	CREATE TABLE IF NOT EXISTS carihesaplar(
		id INTEGER	primary key	autoincrement,
		kod VARCHAR(50),
		unvan VARCHAR(150),
		telefon VARCHAR(15),
		adres TEXT,
		email VARCHAR(50),
		created DATE
	);

	CREATE TABLE IF NOT EXISTS kullanicilar(
		id INTEGER	primary key	autoincrement,
		ad VARCHAR(50),
		email VARCHAR(150),
		sifre VARCHAR(150)
	);
	`
	_, err := db.Exec(sql_table)


	if err != nil { panic(err) }
}




func StoreItem(db *sql.DB, items []CariKart) {
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

func UpdateCariKart(db *sql.DB, item CariKart) {
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



func ReadItem(db *sql.DB, SonId int) ([]CariKart,int){
	rows, err := db.Query(`
Select id,kod, unvan, telefon, adres, email from carihesaplar
where id>?
order by id LIMIT 25
`,SonId)

	if err != nil { panic(err) }
	defer rows.Close()

	var result []CariKart
	for rows.Next() {
		item := CariKart{}
		err2 := rows.Scan(&item.Id, &item.Kod, &item.Unvan,&item.Telefon,&item.Adres,&item.Email)
		if err2 != nil { panic(err2) }
		result = append(result, item)
	}

	var son int

	if len(result)>0{
		son=result[len(result)-1].Id
	}

	return result,son
}

func ReadItemId(db *sql.DB, id int) CariKart {
	item := CariKart{}
	row := db.QueryRow("Select id,kod, unvan, telefon, adres, email from carihesaplar where id=?",id)
	err :=row.Scan(&item.Id, &item.Kod, &item.Unvan,&item.Telefon,&item.Adres,&item.Email)
	if err != nil { panic(err) }
	return item
}
