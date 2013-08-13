package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type AuthJSON struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
	Message string `json:"message"`
}

type Article struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Body   string `json:"body"`
}

type Photo struct {
	Id  int    `json:"id"`
	Src string `json:"src"`
}

var articles []Article
var photos []Photo

func initrecords() {

	article1 := Article {
		Id: 1,
		Title: "How to Write a Javascript Framework",
		Author: "Tomhuda Katzdale",
		Body: "Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
	}
	article2 := Article {
		Id: 2,
		Title: "Chronicles of an Embereño",
		Author: "Alerik Bryneer",
		Body: "Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
	}
	article3 := Article {
		Id: 3,
		Title: "The Eyes of Thomas",
		Author: "Yahuda Katz",
		Body: "Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
	}
	
	articles = append(articles, article1, article2, article3)
	
	photo1 := Photo {
		Id: 1,
		Src: "images/potd.png",
	}
	photo2 := Photo {
		Id: 2,
		Src: "images/yohuda.jpg",
	}
	photo3 := Photo {
		Id: 3,
		Src: "images/easter.jpg",
	}

	photos = append(photos, photo1, photo2, photo3)
}

var currentToken string

func AuthenticationHandler(w http.ResponseWriter, req *http.Request) {

	// Parse the incoming user/pass from the request body
	var authJSON AuthJSON
	err := json.NewDecoder(req.Body).Decode(&authJSON)
	if err != nil {
		panic(err)
	}
	username := authJSON.Username
	password := authJSON.Password

	// If username and password are correct
	if username == "7cadc15d609c4ae9b4be6265b8e1cace16e6fa78a81ab0c7db82e687a7c867a5" && password == "0653216d8920713c0db2c2571a0d6ed0c2b6b94194f37e53a7bb025391dfc8b0" {
		currentToken = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

		var res Response
		res.Token = currentToken
		res.Success = true

		// Serialize the token to JSON
		j, err := json.Marshal(res)
		if err != nil {
			panic(err)
		}

		// Write the response
		w.Header().Set("Content-Type", "application/json")
		w.Write(j)
	} else {
		var res Response
		res.Message = "Invalid username/password"
		res.Success = false

		// Serialize the token to JSON
		j, err := json.Marshal(res)
		if err != nil {
			panic(err)
		}
		// Write the response
		w.Header().Set("Content-Type", "application/json")
		w.Write(j)
	}

}

func ValidTokenProvided(w http.ResponseWriter, req *http.Request) bool {

	userToken := req.FormValue("token")

	if currentToken == "" || userToken != currentToken {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{ "error": "Invalid token. You provided: " + userToken }`))
		return false
	}
	return true
}

func ArticlesHandler(w http.ResponseWriter, r *http.Request) {
	if ValidTokenProvided(w, r) {
		w.Header().Set("Content-Type", "application/json")
		j, err := json.Marshal(articles)
		if err != nil {
			panic(err)
		}
		w.Write(j)
	}
}
func PhotosHandler(w http.ResponseWriter, r *http.Request) {
	if ValidTokenProvided(w, r) {
		w.Header().Set("Content-Type", "application/json")
		j, err := json.Marshal(photos)
		if err != nil {
			panic(err)
		}
		w.Write(j)
	}
}

func main() {
	log.Println("Starting Server")

	initrecords()

	r := mux.NewRouter()
	r.HandleFunc("/api/articles.json", ArticlesHandler).Methods("GET")
	r.HandleFunc("/api/photos.json", PhotosHandler).Methods("GET")
	r.HandleFunc("/api/auth.json", AuthenticationHandler).Methods("POST")
	http.Handle("/api/", r)

	http.Handle("/", http.FileServer(http.Dir("./frontEnd/")))

	log.Println("Listening on 8080")
	http.ListenAndServe(":8080", nil)
}
