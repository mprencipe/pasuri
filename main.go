package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"pasuri/db"

	log "github.com/sirupsen/logrus"
)

func getHashSuffixes(w http.ResponseWriter, req *http.Request) {
	hashPrefix := req.URL.Query().Get("prefix")
	if len(hashPrefix) != 5 {
		http.Error(w, "Hash prefix needs to be exactly 5 characters in length", http.StatusBadRequest)
		log.Error(errors.New("Received invalid hash prefix " + hashPrefix))
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
