package biz

import (
	"database/sql"
	"log"
	"os/exec"
)

var (
	sqliteClient *sql.DB
)

func ConnSqlite()  {
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	sqliteClient =db
}

func GetSqliteClient()*sql.DB  {
	if sqliteClient != nil {

	}
	return sqliteClient
}


//执行cmd
func demo()  {
	//
	exec.Command("")
}