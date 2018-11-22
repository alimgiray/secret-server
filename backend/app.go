package main

import (
	"log"
	"net/http"
	"strings"
	"strconv"
	"time"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	db "./database"
	. "./secret"
)

var mongo = db.Database{}

func CreateSecret(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Access-Control-Allow-Origin", "*")

	r.ParseForm()
	var secretText = r.Form["secret"][0]
	var expireAfterViews = r.Form["expireAfterViews"][0]
	var expireAfter = r.Form["expireAfter"][0]

	if empty(secretText) || empty(expireAfterViews) || empty(expireAfter) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	remainingViews, err := strconv.Atoi(expireAfterViews)
	if err != nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	remainingMinutes, err := strconv.Atoi(expireAfter)
	if err != nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	randomHash, err := uuid.NewRandom()
	if err != nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	secret := Secret {
		Hash: randomHash.String(), 
		SecretText: secretText, 
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Local().Add(time.Minute * time.Duration(remainingMinutes)), 
		RemainingViews: remainingViews,
	}

	mongo.Insert(secret)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(secret)
}

func RetrieveSecret(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
    hash := vars["hash"]

    if empty(hash) {
    	w.WriteHeader(http.StatusNotFound)
		return
    }

    secret, err := mongo.Find(hash)

    if err != nil {
    	w.WriteHeader(http.StatusNotFound)
		return
    }

    if secret.RemainingViews == 0 || secret.ExpiresAt.Sub(time.Now()) < 0 {
    	w.WriteHeader(http.StatusNotFound)
		return
    }

    secret.RemainingViews = secret.RemainingViews - 1
    mongo.Update(secret)

    w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(secret)
}

func empty(s string) bool {
    return len(strings.TrimSpace(s)) == 0
}

func init() {

	mongo.Server = ""
	mongo.Database = "codersrank"
	mongo.Connect()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/v1/secret", CreateSecret).Methods("POST", "OPTIONS")
	r.HandleFunc("/v1/secret/{hash}", RetrieveSecret).Methods("GET", "OPTIONS")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}