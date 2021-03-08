// go build app.go && app.exe
// TODO:
// Add logging to .log file
// Maybe make the home page more prettier? Adding CSS colors might help

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
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
	r.HandleFunc("/url", getShortLinks).Methods("GET")
	r.HandleFunc("/url/{id}", redirectURL).Methods("GET")
	r.HandleFunc("/url/{id}", deleteShortLink).Methods("DELETE")

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

func isLetter(s string) bool {
	// Ref: https://stackoverflow.com/questions/38554353/how-to-check-if-a-string-only-contains-alphabetic-characters-in-go
	for _, r := range s {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
			return false
		}
	}
	return true
}

func isValidURL(testurl string) bool {
	// Ref: https://golangcode.com/how-to-check-if-a-string-is-a-url/
	// Might want to test by sending a request, but this might be a potential security isssue
	_, err := url.ParseRequestURI(testurl)
	if err != nil {
		return false
	}

	u, err := url.Parse(testurl)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

func randString(n int) string {
	key := true
	var b []byte
	// Prevent random key collision (highly unlikely but nice to have)
	for key == true {
		b = make([]byte, n)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		_, key = linkMap[string(b)]
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
}

// Delete shortened link
func deleteShortLink(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	log.Print("Request received to delete link with ID ", id)
	_, key := linkMap[id]
	if key {
		delete(linkMap, id)
		log.Print("ID ", id, " is successfully deleted. New list: ", linkMap)

		saveToFile()

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w)
	} else {
		log.Print("ID ", id, " not found. [For DELETE]")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w)
	}
}

// Generate new shortened link (randomize the shortened URL)
func generateShortLink(w http.ResponseWriter, r *http.Request) {
	var link Link
	decoder := json.NewDecoder(r.Body)
	// Uncomment line below to strictly allow known fields only
	// decoder.DisallowUnknownFields()
	err := decoder.Decode(&link)

	if err != nil {
		log.Print("Error trying to decode POST request with link ", link, " Error: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if link is valid
	if !isValidURL(link.URL) {
		http.Error(w, "The specified URL is invalid.", http.StatusBadRequest)
		return
	}

	log.Print("Request came through to shorten link ", link.URL)

	// Generate random string for the ID (overwrites ID from POST request if specified)
	if link.ID == "" {
		id := randString(6)
		link.ID = id
	} else {
		// Check if custom ID consists of letters only
		if !isLetter(link.ID) {
			http.Error(w, "The requested custom ID must consist of letters only.", http.StatusBadRequest)
			return
		}

		_, key := linkMap[link.ID]
		// Check if custom ID is in use
		if key {
			http.Error(w, "The requested custom ID is in use.", http.StatusConflict)
			return
		}
	}

	linkMap[link.ID] = link

	saveToFile()

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(link)
}

func getShortLinks(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	_, err := json.Marshal(linkMap)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(linkMap)
}

// Redirect to shortened link
func redirectURL(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
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
