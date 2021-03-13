package db

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

const testDbName = "test.db"

func TestFindHashSuffixes(t *testing.T) {
	db := createDatabase()
	_, err := db.Exec("INSERT OR IGNORE INTO hash(prefix, part1, part2, part3) VALUES(?,?,?,?)", 375002, 250735282941349, 26977265306808, 7420971621823)
	if err != nil {
		panic(err)
	}
	db.Close()
	InitDb(testDbName)

	hashes, _ := FindHashSuffixes("5b8da")

	if len(hashes) != 1 {
		t.Errorf("Hash was incorrect, got: %d, wanted: %s.", len(hashes), "correctHash")
	}
}

func createDatabase() *sql.DB {
	_, err := os.Create(testDbName)
	if err != nil {
		panic(err)
	}
	db, err := sql.Open("sqlite3", testDbName)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("create table hash (prefix integer not null, part1 integer not null, part2 integer not null, part3 integer not null, primary key(prefix, part1, part2, part3)) without rowid")
	if err != nil {
		panic(err)
	}
	return db
}
