package main

import(
	"fmt"
	"net/http"
	"encoding/json"
	"math/rand" 
	"strings"
	"net/url"
)


//json respones : we use struct , then encode it using encoding/json 
//***for JSON ,, struct key should always start from capital letter or else , it will be 
//invisible to the json


type getURL struct{
	LongUrl string 
	ShortUrl string
}



var urldata = make(map[string]string)


//fucn to generate shortcode 
func generateShortCode()string{
	var shortcode strings.Builder
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for i:=0;i<=5;i++{
           shortcode.WriteByte(charset[rand.Intn(len(charset))])
	}
	return shortcode.String()
}

//create and return shorturl : input : long output : short  + readjson & return json

func createShortURL(w http.ResponseWriter , r *http.Request){
	var Url getURL 
	err:=json.NewDecoder(r.Body).Decode(&Url)
	if err!=nil{
		fmt.Fprintf(w,"error in decoding url")
		return 
	}
	var code string
	code = generateShortCode()
	//check is code is already being used or not 
    _,exists :=urldata[code]
	if exists{
		code = generateShortCode()
	}
	//update hasmap 
	urldata[code]=Url.LongUrl 

	//return the json with shorturl 
	shortURL:=fmt.Sprintf("http://localhost:8080/%v",code)
    
	json.NewEncoder(w).Encode(shortURL) 
	
   
}

//func to redirect : input:short output :long   + read and return json 

func redirect(w http.ResponseWriter , r *http.Request){
	//get the long url 
   var Url getURL 
   err:=json.NewDecoder(r.Body).Decode(&Url) 
   if err!=nil{
		fmt.Fprintf(w,"error in decoding url")
		return 
	}
	//get the code from the url 
	u,err:=url.Parse(Url.ShortUrl)
	if err != nil {
		fmt.Fprintf(w,"error in getting code ")
		return
	}
	code:=strings.TrimPrefix(u.Path,"/")
    
	//check if code exist in map 
	 _,exists :=urldata[code]
	if !exists{
		fmt.Fprintf(w,"invalid shorturl ")
		return
	}
    
	//send map 
  
   	json.NewEncoder(w).Encode(urldata[code]) 

}



func main(){
	http.HandleFunc("/createShortURL", createShortURL)
	http.HandleFunc("/redirectToLongURL", redirect)
	fmt.Println("server started at localhost 8080 ")
    err := http.ListenAndServe(":8080",nil) //this returns just a error and not usual res,err
	if(err!=nil){
      fmt.Println("error starting the server at 8080 ",err)
	}

}



