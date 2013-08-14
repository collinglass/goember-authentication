# Go Ember Token-based Authentication

### Note: ###

Small addition to Embercasts tutorial, I used crypto.js and added a basic SHA256 encryption hash to the user/pass before sending it to server. [See here](http://code.google.com/p/crypto-js/#SHA-2) and [here](http://code.google.com/p/crypto-js/#The_Hasher_Output) for algorithm used, (changed to hex string). And [here](http://www.movable-type.co.uk/scripts/sha256.html) is a simple resource to generate any hash string if you want to change the authentication user/pass.

## Documentation ##

Simple Go authentication server to use with [embercasts authentication screencasts](http://www.embercasts.com/).


### Import ###

```
import (
	"net/http"
```
The library for web apps/protocols.
```
"log"
```
The library to log to the console
```
"github.com/gorilla/mux"
```
We need this for routing.
```
    "encoding/json"
)
```
We need this for handling JSON (encoding/decoding).

### Articles ###

```
type Article struct {
	Id		int 	`json:"id"`
	Title	string 	`json:"title"`
	Author	string	`json:"author"`
	Body	string	`json:"body"`
}

...

var articles []Article

func initrecords(){

var article1 Article
article1.Id = 1
article1.Title = "How to Write a Javascript Framework"
article1.Author = "Tomhuda Katzdale"
article1.Body = "..."

var article2 Article
article2.Id = 2
article2.Title = "Chronicles of an Embereño"
article2.Author = "Alerik Bryneer"
article2.Body = "..."

var article3 Article
article3.Id = 3
article3.Title = "The Eyes of Thomas"
article3.Author = "Yahuda Katz"
article3.Body = "..."

articles = append(articles, article1)
articles = append(articles, article2)
articles = append(articles, article3)

...

}

func main(){
	...
	initrecords()
	...
}
```


Create a go struct for the article object. On initialization: Instantiate an array of Article objects. Then individually instantiate the 3 article objects and append them to articles array. This gives us a clean structure to handle articles.

```
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

func main(){
	...
    r := mux.NewRouter()
    r.HandleFunc("/api/articles.json", ArticlesHandler).Methods("GET")
    ...
    http.Handle("/api/", r)
	...
}
```

Handle the GET request for displaying the Articles with JSON. If valid token is provided, encode articles array to JSON and send to client.

### Photos ###

Same as Articles.

### Authentication ###

```
type AuthJSON struct {
    Username    string `json:"username"`
    Password 	string `json:"password"`
}
```

AuthJSON struct creates go object for incoming username and password.

```
type Response struct {
    Success		bool `json:"success"`
    Token 		string `json:"token"`
	Message 	string `json:"message"`
}
```

Response struct creates go object for outgoing response. If valid username/password combination is submitted token/success response is sent; if invalid username/password combination is submitted, message/success response is sent.

```
var currentToken string
```

currentToken holds authentication token and persists for login session. 

```
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
		
	    // Serialize the response to JSON
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
		
	    // Serialize the response to JSON
	    j, err := json.Marshal(res)
	    if err != nil {
	        panic(err)
	    }
	    // Write the response
	    w.Header().Set("Content-Type", "application/json")
	    w.Write(j)
	}
	
}
```

Decodes incoming username/password JSON, puts it in AuthJSON object.

If username/password combination is valid, set currentToken and success:true serialize it to JSON and pass back to client.
If username/password combination is invalid, set error message and success:false serialize it to JSON and pass back to client.

```
func ValidTokenProvided(w http.ResponseWriter, req *http.Request) bool {
	userToken := req.FormValue("token")
	
	if( currentToken == "" || userToken != currentToken ) {
	  	w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{ "error": "Invalid token. You provided: " + userToken }`))
		return false
	}
	return true
}
```

If currentToken is empty or userToken does not match currentToken then raise 401 error and return false, else return true.

### Full server.go Code ###
<a id="fullcode"></a> 

``` server.go

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
		
	    // Serialize the response to JSON
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
		
	    // Serialize the response to JSON
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
	
	if( currentToken == "" || userToken != currentToken ) {
	  	w.WriteHeader(http.StatusUnauthorized)
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

```

Now follow the [embercasts](http://www.embercasts.com/) to understand how the client interacts with this code.