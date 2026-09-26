package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool" //pgx
	"github.com/joho/godotenv"        //env
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
)

//json respones : we use struct , then encode it using encoding/json
//***for JSON ,, struct key should always start from capital letter or else , it will be
//invisible to the json

type ShortenRequest struct {
	LongURL string
}

// unused
type ShortenResposne struct {
	ShortURL string
}

// declare pool var at package level since we need to use in handlers
var pool *pgxpool.Pool

// fucn to generate shortcode
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

	var Url ShortenRequest
	err := json.NewDecoder(r.Body).Decode(&Url)
	if err != nil {
		http.Error(w, "Error decoding the Longurl", 400)
		return
	}
	var code string
	code = generateShortCode()
	//check is code is already being used or not

	//update hasmap
	_, err = pool.Exec(r.Context(), "INSERT INTO urls (code,long_url) VALUES($1 , $2)", code, Url.LongURL)
	if err != nil {
		http.Error(w, "Error updating the db , try again ", 400)
		return
	}
	//return the json with shorturl
	shorturl := fmt.Sprintf("http://localhost:8080/%v", code)

	json.NewEncoder(w).Encode(ShortenResposne{ShortURL: shorturl})

}

//func to redirect

func redirect(w http.ResponseWriter, r *http.Request) {

	//get the code from the short url clicked
	u, err := url.Parse(r.URL.Path)
	if err != nil {
		http.Error(w, "Error getting the code ", http.StatusBadRequest)
		return
	}
	code := strings.TrimPrefix(u.Path, "/")

	//check if code exist in map
	var longURL string
	err = pool.QueryRow(r.Context(), "SELECT long_url FROM urls WHERE code = $1", code).Scan(&longURL)
	if err != nil {
		http.Error(w, "Error : code not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, longURL, http.StatusFound)

}

func main() {
	//create a connection pool
	godotenv.Load()
	databaseURL := os.Getenv("DATABASE_URL")
	ctx := context.Background()
	var err error
	pool, err = pgxpool.New(ctx, databaseURL)
	if err != nil {
		fmt.Println("error creating database pool", err)
		return
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		fmt.Println("error connecting to database:", err)
		return
	}

	fmt.Println("database connection sucessful")

	http.HandleFunc("/createShortURL", createShortURL)
	http.HandleFunc("/", redirect)
	fmt.Println("server started at localhost 8080 ")
	err = http.ListenAndServe(":8080", nil) //this returns just a error and not usual res,err
	if err != nil {
		fmt.Println("error starting the server at 8080 ", err)
	}

}
