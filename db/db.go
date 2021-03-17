package db

import (
	"database/sql"
	"fmt"
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
	}

	var err error
	db, err = sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal("Couldn't open database")
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
		hashSuffix := leftPadWithZeroes(strconv.FormatInt(dbHashPart1, 16), 12) + leftPadWithZeroes(strconv.FormatInt(dbHashPart2, 16), 12) + leftPadWithZeroes(strconv.FormatInt(dbHashPart3, 16), 11)
		hashSuffixes = append(hashSuffixes, hashSuffix)
	}
	rows.Close()

	return hashSuffixes, nil
}

func leftPadWithZeroes(hashPart string, length int) string {
	return fmt.Sprintf("%0*s", length, hashPart)
}
