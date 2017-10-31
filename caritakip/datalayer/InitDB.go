package datalayer

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

	CREATE TABLE IF NOT EXISTS carihareket(
		id INTEGER	primary key	autoincrement,
		carihesapid INTEGER,
		ba BIT,
		tutar DOUBLE,
		FOREIGN KEY(carihesapid) REFERENCES carihesaplar(id)
	);
`
	_, err := db.Exec(sql_table)
	if err != nil { panic(err) }


}

