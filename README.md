# Go Ember Token-based Authentication

## Documentation ##

Simple Go authentication server to use with [embercasts authentication screencasts](http://www.embercasts.com/).


### Import ###

{% codeblock %}
import (
	"net/http"
{% endcodeblock %}
The library for web apps/protocols.
{% codeblock %}
"log"
{% endcodeblock %}
The library to log to the console
{% codeblock %}
"github.com/gorilla/mux"
{% endcodeblock %}
We need this for routing.
{% codeblock %}
	"encoding/json"
)
{% endcodeblock %}
We need this for handling JSON (encoding/decoding).

### Articles ###

{% codeblock %}
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
{% endcodeblock %}

Create a go struct for the article object. On initialization: Instantiate an array of Article objects. Then individually instantiate the 3 article objects and append them to articles array. This gives us a clean structure to handle articles.

{% codeblock %}
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
{% endcodeblock %}

Handle the GET request for displaying the Articles with JSON. If valid token is provided, encode articles array to JSON and send to client.

### Photos ###

Same as Articles.

### Authentication ###

{% codeblock %}
type AuthJSON struct {
    Username    string `json:"username"`
    Password 	string `json:"password"`
}
{% endcodeblock %}

AuthJSON struct creates go object for incoming username and password.

{% codeblock %}
type Response struct {
    Success		bool `json:"success"`
    Token 		string `json:"token"`
	Message 	string `json:"message"`
}
{% endcodeblock %}

Response struct creates go object for outgoing response. If valid username/password combination is submitted token/success response is sent; if invalid username/password combination is submitted, message/success response is sent.

{% codeblock %}
var currentToken string
{% endcodeblock %}

currentToken holds authentication token and persists for login session. 

{% codeblock %}
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
{% endcodeblock %}

Decodes incoming username/password JSON, puts it in AuthJSON object.

If username/password combination is valid, set currentToken and success:true serialize it to JSON and pass back to client.
If username/password combination is invalid, set error message and success:false serialize it to JSON and pass back to client.

{% codeblock %}
func ValidTokenProvided(w http.ResponseWriter, req *http.Request) bool {
	userToken := req.FormValue("token")
	
	if( currentToken == "" || userToken != currentToken ) {
	  	w.WriteHeader(401)
		w.Write([]byte(`{ "error": "Invalid token. You provided: " + userToken }`))
		return false
	}
	return true
}
{% endcodeblock %}

If currentToken is empty or userToken does not match currentToken then raise 401 error and return false, else return true.

Now follow the [embercasts](http://www.embercasts.com/) to understand how the client interacts with this code.