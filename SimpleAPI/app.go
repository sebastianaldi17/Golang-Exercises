// https://medium.com/@supun.muthutantrige/lets-go-create-an-awesome-rest-api-in-go-part-iii-d9c43e797abc
// To do list:
// Save and load to local json

// go build app.go && app.exe

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

// Account Json request payload is as follows,
// {
//  "id": "1",
//  "first_name": "John",
//  "last_name":  "Doe",
//  "user_name":  "JDoe123"
// }
type Account struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"user_name"`
}

var accountMap map[string]Account

func main() {
	accountMap = make(map[string]Account)
	r := mux.NewRouter()
	r.HandleFunc("/", homePageHandler).Methods("GET")
	r.HandleFunc("/accounts", accountPageHaandler).Methods("GET")
	r.HandleFunc("/accounts", createAccountHandler).Methods("POST")
	r.HandleFunc("/accounts/{id}", getAccountHandler).Methods("GET")
	r.HandleFunc("/accounts/{id}", deleteAccountHandler).Methods("DELETE")

	if fileExists("db.json") {
		// Load json file
		jsonFile, err := os.Open("db.json")
		if err != nil {
			log.Fatal(err)
		}
		byteData, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteData, &accountMap)
		log.Print("Existing db.json file found: ", accountMap)

	} else {
		log.Print("No existing db.json file found.")
		// Create json file
		file, err := os.Create("db.json")
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
	}

	log.Print("Starting API at port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// IO Operations
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func saveToFile() {
	data, err := json.MarshalIndent(accountMap, "", " ")
	if err != nil {
		log.Print("An error occured while converting map to json.")
		log.Fatal(err)
	}
	err = ioutil.WriteFile("db.json", data, 0644)
	if err != nil {
		log.Print("An error occured while writing to file.")
		log.Fatal(err)
	}
}

// Handlers
func homePageHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Request came through to Home Page")
	w.Header().Add("Content-Type", "application/json")
	_, err := json.Marshal(accountMap)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accountMap)
}

func accountPageHaandler(w http.ResponseWriter, r *http.Request) {
	// Tells user to use POST instead of GET
	// This can be replaced with a form that POSTS to itself (Potential feature)
	http.Error(w, "Please use POST instead of GET when trying to create a new account.", http.StatusMethodNotAllowed)
}

func createAccountHandler(w http.ResponseWriter, r *http.Request) {
	var account Account
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&account) // Decode JSON data and dump to var account
	log.Print("Request received to create an Account with data ", account)

	// JSON checking
	if err != nil {
		log.Print("Request received to create an Account with data ", account, " rejected for bad format")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if account.ID == 0 || account.FirstName == "" || account.LastName == "" || account.UserName == "" {
		log.Print("Request received to create an Account with data ", account, " rejected for missing/default value field(s)")
		http.Error(w, "Missing or default value field", http.StatusBadRequest)
		return
	}
	if decoder.More() {
		log.Print("Request received to create an Account with data ", account, " rejected for additional field(s)")
		http.Error(w, "Extra data on JSON object", http.StatusBadRequest)
		return
	}

	// Also serves as update, unless you want to make a separate function for it...
	// Check existing ID by using
	// if _, key := accountMap[id]; key {
	//     // Do something
	// }
	id := strconv.Itoa(account.ID)
	accountMap[id] = account
	log.Print("Added ", account, " to list of accounts. List: ", accountMap)

	saveToFile()

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

func getAccountHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) // Gets the :id from the URL
	id := params["id"]
	log.Print("Request received to get an Account with ID ", id)
	account, key := accountMap[id]
	w.Header().Add("Content-Type", "application/json")
	if key {
		log.Print("ID ", id, " is found (OK).")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(account)
	} else {
		log.Print("ID ", id, " not found [For GET].")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w)
	}
}

func deleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	log.Print("Request received to delete an Account with ID ", id)
	_, key := accountMap[id] // Account information not needed so just _ the first var
	if key {
		delete(accountMap, id)
		log.Print("ID ", id, " is successfully deleted. New list: ", accountMap)

		saveToFile()

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w)
	} else {
		log.Print("ID ", id, " not found. [For DELETE]")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w)
	}
}
