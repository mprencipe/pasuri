package main

import (
	"bufio"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"io"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Starting up..")

	if len(os.Args) < 2 {
		log.Fatal("Supply at least one password file")
		os.Exit(1)
	}

	db, err := sql.Open("sqlite3", "./pass.db")
	if err != nil {
		log.Fatal("Couldn't open database", err)
		os.Exit(1)
	}

	for _, fileName := range os.Args[1:] {
		saveHashes(fileName, db)
	}

	log.Info("..done!")
}

func saveHashes(fileName string, db *sql.DB) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal("Unable to read file", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	i := 0
	log.Info("Starting to hash", fileName)
	for scanner.Scan() {
		if i%100 == 0 {
			log.Info("Read ", i, " lines")
		}

		hasher := sha1.New()
		io.WriteString(hasher, scanner.Text())
		hash := hex.EncodeToString(hasher.Sum(nil))
		hashPrefix, err := strconv.ParseInt(hash[:5], 16, 64)
		exitOnError(err)
		hashPart1, err := strconv.ParseInt(hash[5:17], 16, 64)
		exitOnError(err)
		hashPart2, err := strconv.ParseInt(hash[17:29], 16, 64)
		exitOnError(err)
		hashPart3, err := strconv.ParseInt(hash[29:], 16, 64)
		exitOnError(err)
		_, err = db.Exec("INSERT OR IGNORE INTO hash(prefix, part1, part2, part3) VALUES(?,?,?,?)", hashPrefix, hashPart1, hashPart2, hashPart3)
		exitOnError(err)

		i++
	}
	log.Info(fileName, "..done!")
}

func exitOnError(err error) {
	if err != nil {
		log.Fatal("Unable to insert hash", err)
		os.Exit(1)
	}
}
