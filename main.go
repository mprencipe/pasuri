package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"pasuri/db"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func getHashSuffixes(w http.ResponseWriter, req *http.Request) {
	hashPrefix := req.URL.Query().Get("prefix")
	if len(hashPrefix) != 5 {
		http.Error(w, "Hash prefix needs to be exactly 5 characters in length", http.StatusBadRequest)
		log.Error(errors.New("Received invalid hash prefix " + hashPrefix))
		return
	}
	_, err := strconv.ParseUint(hashPrefix, 16, 64)
	if err != nil {
		http.Error(w, "Hash prefix needs to be in hexadecimal format: 0-9, a-f", http.StatusBadRequest)
		log.Error(errors.New("Received non-hex string " + hashPrefix))
		return
	}
	hashSuffixes, err := db.FindHashSuffixes(hashPrefix)
	if err != nil {
		http.Error(w, "Error reading hash count", http.StatusInternalServerError)
		return
	}
	hashSuffixesJson, err := json.Marshal(hashSuffixes)
	if err != nil {
		http.Error(w, "Error reading hash count", http.StatusInternalServerError)
		log.Error(err)
		return
	}

	if cors, ok := os.LookupEnv("CORS"); ok {
		w.Header().Set("Access-Control-Allow-Origin", cors)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", string(hashSuffixesJson))
}

func main() {
	log.Info("Starting up..")

	db.InitDb("pass.db")

	log.Info("..running!")

	http.HandleFunc("/hash", getHashSuffixes)
	http.ListenAndServe(":8080", nil)
}
