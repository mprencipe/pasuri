package main

import (
	"bufio"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"io"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Starting up..")

	if len(os.Args) < 2 {
		log.Fatal("Supply at least one password file")
	}

	db, err := sql.Open("sqlite3", "./pass.db")
	if err != nil {
		log.Fatal("Couldn't open database", err)
	}

	for _, fileName := range os.Args[1:] {
		saveHashes(fileName, db)
	}

	log.Info("..done!")
}

func saveHashes(fileNameAndType string, db *sql.DB) {
	parts := strings.Split(fileNameAndType, ":")
	fileName := parts[0]
	fileType := parts[1]

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal("Unable to read file ", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	i := 0
	log.Info("Starting to write file ", fileName)
	for scanner.Scan() {
		if i%100 == 0 {
			log.Info("Read ", i, " lines")
		}

		var hash string
		if fileType == "hash" {
			hash = scanner.Text()
		} else {
			hash = MakeHash(scanner.Text())
		}

		writeHash(hash, db)

		i++
	}
	log.Info("..done writing ", fileName)
}

func MakeHash(str string) string {
	hasher := sha1.New()
	io.WriteString(hasher, str)
	return hex.EncodeToString(hasher.Sum(nil))
}

func writeHash(hash string, db *sql.DB) {
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
}

func exitOnError(err error) {
	if err != nil {
		log.Fatal("Unable to insert hash", err)
	}
}
