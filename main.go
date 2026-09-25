package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

//json respones : we use struct , then encode it using encoding/json
//***for JSON ,, struct key should always start from capital letter or else , it will be
//invisible to the json

type ShortenRequest struct {
	LongURL string
}

//unused
type ShortenResposne struct {
	ShortURL string
}

var urldata = make(map[string]string)
var mu sync.Mutex

//fucn to generate shortcode
func generateShortCode() string {
	var shortcode strings.Builder
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for i := 0; i <= 5; i++ {
		shortcode.WriteByte(charset[rand.Intn(len(charset))])
	}
	return shortcode.String()
}

//create and return shorturl : input : long output : short  + readjson & return json

func createShortURL(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var Url ShortenRequest
	err := json.NewDecoder(r.Body).Decode(&Url)
	if err != nil {
		http.Error(w, "Error decoding the Longurl", 400)
		return
	}
	var code string
	code = generateShortCode()
	//check is code is already being used or not
	//keep generating new code if matching with already present
	_, exists := urldata[code]
	for exists {
		code = generateShortCode()

		_, exists = urldata[code] //re-check with the new code
		//here its = and not := , since we are updating the same exist and not creating a new one
	}
	//update hasmap
	urldata[code] = Url.LongURL

	//return the json with shorturl
	shorturl := fmt.Sprintf("http://localhost:8080/%v", code)

	json.NewEncoder(w).Encode(ShortenResposne{ShortURL: shorturl})

}

//func to redirect

func redirect(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	//get the code from the short url clicked
	u, err := url.Parse(r.URL.Path)
	if err != nil {
		http.Error(w, "Error getting the code ", http.StatusBadRequest)
		return
	}
	code := strings.TrimPrefix(u.Path, "/")

	//check if code exist in map
	_, exists := urldata[code]
	if !exists {
		http.Error(w, "Error : code not found", http.StatusNotFound)
		return
	}

	//send map

	http.Redirect(w, r, urldata[code], http.StatusFound)

}

func main() {
	http.HandleFunc("/createShortURL", createShortURL)
	http.HandleFunc("/", redirect)
	fmt.Println("server started at localhost 8080 ")
	err := http.ListenAndServe(":8080", nil) //this returns just a error and not usual res,err
	if err != nil {
		fmt.Println("error starting the server at 8080 ", err)
	}

}
