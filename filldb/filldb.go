package main

import (
	"bufio"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

const chunkSize = 7000

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
		processFile(fileName, db)
	}

	log.Info("..done!")
}

func processFile(fileNameAndType string, db *sql.DB) {
	parts := strings.Split(fileNameAndType, ":")
	fileName := parts[0]
	fileType := parts[1]

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal("Unable to read file ", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	log.Info("Starting to write file ", fileName)

	rowCount := 0
	chunks := []string{}
	for scanner.Scan() {
		rowCount++
		var hash string
		if fileType == "hash" {
			hash = scanner.Text()
		} else {
			hash = MakeHash(scanner.Text())
		}
		chunks = append(chunks, hash)
		if len(chunks) == chunkSize {
			log.Info("Read ", rowCount, " lines")
			saveHashes(chunks, db)
			chunks = nil
		}
	}
	if len(chunks) > 0 {
		saveHashes(chunks, db)
	}
	log.Info("..done writing ", fileName)
}

func saveHashes(hashes []string, db *sql.DB) {
	valueStrings := []string{}
	valueArgs := []interface{}{}

	for _, hash := range hashes {
		valueStrings = append(valueStrings, "(?, ?, ?, ?)")
		hashPrefix, err := strconv.ParseInt(hash[:5], 16, 64)
		exitOnError(err)
		hashPart1, err := strconv.ParseInt(hash[5:17], 16, 64)
		exitOnError(err)
		hashPart2, err := strconv.ParseInt(hash[17:29], 16, 64)
		exitOnError(err)
		hashPart3, err := strconv.ParseInt(hash[29:], 16, 64)
		exitOnError(err)
		valueArgs = append(valueArgs, hashPrefix)
		valueArgs = append(valueArgs, hashPart1)
		valueArgs = append(valueArgs, hashPart2)
		valueArgs = append(valueArgs, hashPart3)
	}

	smt := `INSERT OR IGNORE INTO hash(prefix, part1, part2, part3) VALUES %s`
	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))
	_, err := db.Exec(smt, valueArgs...)
	exitOnError(err)
}

func MakeHash(str string) string {
	hasher := sha1.New()
	io.WriteString(hasher, str)
	return hex.EncodeToString(hasher.Sum(nil))
}

func exitOnError(err error) {
	if err != nil {
		log.Fatal("Unable to insert hash", err)
	}
}
