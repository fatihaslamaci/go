package main_test

import (
	"database/sql"
	"testing"

	"github.com/fatihaslamaci/go/caritakip/datalayer"
	_ "github.com/mattn/go-sqlite3"
)



func TestSelectCariHareket(t *testing.T) {
	var db *sql.DB
	const dbpath = "testdb.sqlite"
	db = datalayer.InitDB(dbpath)
	defer db.Close()
	datalayer.CreateTable(db)
	datalayer.ReadItem(db, 100000)

}
