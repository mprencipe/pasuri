package main

import (
	"bufio"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"io"
	"os"

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
			log.Debug("Read", i, "lines")
		}

		hasher := sha1.New()
		io.WriteString(hasher, scanner.Text())
		hash := hex.EncodeToString(hasher.Sum(nil))
		hashPrefix := hash[:5]
		hashSuffix := hash[5:]
		_, err = db.Exec("INSERT OR IGNORE INTO hash(prefix, suffix) VALUES(?, ?)", hashPrefix, hashSuffix)
		if err != nil {
			log.Fatal("Unable to insert hash", err)
			os.Exit(1)
		}

		i++
	}
	log.Info(fileName, "..done!")
}
