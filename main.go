package main

import(
	"fmt"
	"net/http"
	"encoding/json" 
)


//json respones : we use struct , then encode it using encoding/json 
//***for JSON ,, struct key should always start from capital letter or else , it will be 
//invisible to the json
type  greetings struct{
   Message string
}

type Longurl struct{
	URL string 
}

//exatc fucntion handle to be made 
func printmsg(w http.ResponseWriter , r *http.Request){
    // write somehtingn to w 
	fmt.Fprint(w,"hello this is my first go backend code :)")
	//can also write this using w.Write([]byte(""))
}

func readjson(w http.ResponseWriter , r *http.Request){
	var req Longurl 
	err:= json.NewDecoder(r.Body).Decode(&req)//similar to sacnl(&s) , we give address of the varibale where we want it to be stored 
	if err!=nil{
		fmt.Fprintf(w,"error decoding : %v",err)
		return 
	}
	fmt.Fprintf(w,"got:%s" , req.URL) 

}

func prams(w http.ResponseWriter , r * http.Request){
    name:=r.URL.Query().Get("name")
	msg := greetings{Message: "hello "+name}  
	json.NewEncoder(w).Encode(msg)
}


func main(){
     
	http.HandleFunc("/",printmsg) //does not return aything 
	http.HandleFunc("/greet",prams)
	http.HandleFunc("/shorten",readjson)
	fmt.Println("server started at localhost 8080 ")
    err := http.ListenAndServe(":8080",nil) //this returns just a error and not usual res,err
	if(err!=nil){
      fmt.Println("error starting the server at 8080 ",err)
	}

}



