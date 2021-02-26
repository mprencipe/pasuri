package db

import (
	"database/sql"
	"os"
	"strconv"

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
	hashPrefixInt, err := strconv.ParseInt(hashPrefix, 16, 64)
	if err != nil {
		return hashSuffixes, err
	}

	rows, err := db.Query("SELECT prefix, part1, part2, part3 FROM hash WHERE prefix = ?", hashPrefixInt)
	if err != nil {
		return hashSuffixes, err
	}

	var dbHashPrefix int
	var dbHashPart1 int64
	var dbHashPart2 int64
	var dbHashPart3 int64
	for rows.Next() {
		err = rows.Scan(&dbHashPrefix, &dbHashPart1, &dbHashPart2, &dbHashPart3)
		if err != nil {
			return hashSuffixes, err
		}
		hashSuffix := strconv.FormatInt(dbHashPart1, 16) + strconv.FormatInt(dbHashPart2, 16) + strconv.FormatInt(dbHashPart3, 16)
		hashSuffixes = append(hashSuffixes, hashSuffix)
	}
	rows.Close()

	return hashSuffixes, nil
}
