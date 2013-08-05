package main

import (
  "net/http"
  "log"
  "github.com/gorilla/mux"
  "encoding/json"
  "io/ioutil"
)

type AuthJSON struct {
    Username    string `json:"username"`
    Password 	string `json:"password"`
}

type Response struct {
    Success		bool `json:"success"`
    Token 		string `json:"token"`
	Message 	string `json:"message"`
}

type Article struct {
	Id		int 	`json:"id"`
	Title	string 	`json:"title"`
	Author	string	`json:"author"`
	Body	string	`json:"body"`
}

type Photo struct {
	Id		int 	`json:"id"`
	Src		string 	`json:"src"`
}

type PhotosJSON struct {
	Photos []Photo `json:"photos"`
}

type ArticlesJSON struct {
    Articles []Article `json:"articles"`
}

var articles []Article
var photos []Photo

func initrecords() {

var article1 Article
article1.Id = 1
article1.Title = "How to Write a Javascript Framework"
article1.Author = "Tomhuda Katzdale"
article1.Body = "Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

var article2 Article
article2.Id = 2
article2.Title = "Chronicles of an Embereño"
article2.Author = "Alerik Bryneer"
article2.Body = "Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

var article3 Article
article3.Id = 3
article3.Title = "The Eyes of Thomas"
article3.Author = "Yahuda Katz"
article3.Body = "Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

articles = append(articles, article1)
articles = append(articles, article2)
articles = append(articles, article3)

var photo1 Photo
photo1.Id = 1
photo1.Src = "images/potd.png"

var photo2 Photo
photo2.Id = 2
photo2.Src = "images/yohuda.jpg"

var photo3 Photo
photo3.Id = 3
photo3.Src = "images/easter.jpg"

photos = append(photos, photo1)
photos = append(photos, photo2)
photos = append(photos, photo3)

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
	if username == "ember" && password == "casts" {
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
	
	b, _ := ioutil.ReadAll(req.Body)
	log.Println(b)
	
	/*var body Response
	err := json.NewDecoder(req.Body).Decode(&body)
    if err != nil {
        panic(err)
    }
	
	
	userToken := body.Token */
	
	if( currentToken == "" /*|| userToken != currentToken*/ ) {
		// log.Println(errors.New("math: square root of negative number"))
	  	w.WriteHeader(401)
		w.Write([]byte(`{ "error": "Invalid token. You provided: " + userToken }`))
		return false
	}
	return true
}

func ArticlesHandler(w http.ResponseWriter, r *http.Request) {
  	if(ValidTokenProvided(w,r)) {
      w.Header().Set("Content-Type", "application/json")
      j, err := json.Marshal(articles)
      if err != nil {
          panic(err)
      }
      w.Write(j)
	}
}
func PhotosHandler(w http.ResponseWriter, r *http.Request) {
    if(ValidTokenProvided(w,r)) {
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
