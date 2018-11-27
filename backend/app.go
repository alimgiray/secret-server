package main

import (
	"log"
	"net/http"
	"strings"
	"strconv"
	"time"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	db "./database"
	. "./secret"
)

type response interface {
	encode(w http.ResponseWriter, s Secret)
}

type jsonResponse struct {}
type xmlResponse struct {}

func (r jsonResponse) encode (w http.ResponseWriter, s Secret) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func (r xmlResponse) encode (w http.ResponseWriter, s Secret) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/xml")
 	xml.NewEncoder(w).Encode(s)
}

var mongo = db.Database{}
var responseType response

func CreateSecret(w http.ResponseWriter, r *http.Request) {
	
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

	responseType.encode(w, secret)
}

func RetrieveSecret(w http.ResponseWriter, r *http.Request) {

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

	responseType.encode(w, secret)
}

func ChangeResponseType(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
    respType := vars["type"]

    if respType == "xml" {
    	responseType = xmlResponse{}
    } else {
		responseType = jsonResponse{}
    }

    w.WriteHeader(http.StatusOK)
}

func empty(s string) bool {
    return len(strings.TrimSpace(s)) == 0
}

func init() {
	readConfig()
	setupDatabase()
}

func readConfig() {

	f, err := ioutil.ReadFile("config")
    if err != nil {
        log.Println(err)
    }

    configStr := string(f)

    if configStr == "xml" {
    	responseType = xmlResponse{}
    } else {
		responseType = jsonResponse{}
    }

}

func setupDatabase() {
	mongo.Server = ""
	mongo.Database = "codersrank"
	mongo.Connect()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/v1/secret", CreateSecret).Methods("POST", "OPTIONS")
	r.HandleFunc("/v1/secret/{hash}", RetrieveSecret).Methods("GET", "OPTIONS")
	r.HandleFunc("/v1/responseType/{type}", ChangeResponseType).Methods("GET", "OPTIONS")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}