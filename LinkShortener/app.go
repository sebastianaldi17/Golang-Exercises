// go build app.go && app.exe
// TODO:
// Add form at home page to shorten a website
// URL validation (string validation, or try to send a request to the specified webpage)
// Maybe make the home page more prettier? CSS might help

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Link Json request payload is as follows,
// {
//  "id": "SHORTlink",
//  "url": "https://www.google.com/"
// }
type Link struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type content struct {
	Links map[string]Link
}

var linkMap map[string]Link

func main() {
	rand.Seed(time.Now().UnixNano())

	linkMap = make(map[string]Link)
	r := mux.NewRouter()

	// Set handlers
	r.HandleFunc("/", getHomePage).Methods("GET")
	r.HandleFunc("/url", generateShortLink).Methods("POST")
	r.HandleFunc("/url/{id}", redirectURL).Methods("GET")

	if fileExists("db.json") {
		// Load json file
		jsonFile, err := os.Open("db.json")
		if err != nil {
			log.Fatal(err)
		}
		byteData, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteData, &linkMap)
		log.Print("Existing db.json file found: ", linkMap)

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

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
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
	data, err := json.MarshalIndent(linkMap, "", " ")
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
// Show all shortened links
func getHomePage(w http.ResponseWriter, r *http.Request) {
	log.Print("Request came through to Home Page")
	p := &content{Links: linkMap}
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, p)
	// w.Header().Add("Content-Type", "application/json")
	// _, err := json.Marshal(linkMap)
	// if err != nil {
	// 	log.Print(err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(linkMap)
}

// Generate new shortened link (randomize the shortened URL)
func generateShortLink(w http.ResponseWriter, r *http.Request) {
	var link Link
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&link)

	if err != nil {
		log.Print("Error trying to decode POST request with link ", link, " Error: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Check extra fields with decoder.More() (returns boolean) if needed

	// Generate random string for the ID (overwrites ID from POST request if specified)
	id := randString(6)
	link.ID = id
	linkMap[id] = link

	saveToFile()

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(link)
}

// Redirect to shortened link
func redirectURL(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) // Gets the :id from the URL
	id := params["id"]
	log.Print("Request received to get URL with shortened link ", id)
	link, key := linkMap[id]
	w.Header().Add("Content-Type", "application/json")

	if key {
		// Link found, redirect to URL
		http.Redirect(w, r, link.URL, http.StatusSeeOther)
	} else {
		log.Print("Shortened link with ID ", id, " not found.")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w)
	}

}