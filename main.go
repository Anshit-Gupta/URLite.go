package main

import(
	"fmt"
	"net/http"
)

//exatc fucntion handle to be made 
func printmsg(w http.ResponseWriter , r *http.Request){
    // write somehtingn to w 
	fmt.Fprint(w,"hello this is my first go backend code :)")
	//can also write this using w.Write([]byte(""))
}

func prams(w http.ResponseWriter , r * http.Request){
    name:=r.URL.Query().Get("name")
	fmt.Fprintf(w,"hello %v ",name)
}


func main(){
     
	http.HandleFunc("/",printmsg) //does not return aything 
	http.HandleFunc("/greet",prams)
	fmt.Println("server started at localhost 8080 ")
    err := http.ListenAndServe(":8080",nil) //this returns just a error and not usual res,err
	if(err!=nil){
      fmt.Println("error starting the server at 8080 ",err)
	}

}



