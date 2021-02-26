package db

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

var db *sql.DB

func InitDb(dbName string) {
	_, fileErr := os.Stat(dbName)
	log.Debug("Checking if database file exists " + dbName)
	if os.IsNotExist(fileErr) {
		log.Fatal("Database file does not exist " + dbName)
		os.Exit(1)
	}

	var err error
	db, err = sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal("Couldn't open database")
		os.Exit(1)
	}
}

func FindHashSuffixes(hashPrefix string) ([]string, error) {
	hashSuffixes := make([]string, 0)

	rows, err := db.Query("SELECT suffix FROM hash WHERE prefix = ?", hashPrefix)
	if err != nil {
		return hashSuffixes, err
	}

	var hashSuffix string
	for rows.Next() {
		err = rows.Scan(&hashSuffix)
		if err != nil {
			return hashSuffixes, err
		}
		hashSuffixes = append(hashSuffixes, hashSuffix)
	}
	rows.Close()

	return hashSuffixes, nil
}
