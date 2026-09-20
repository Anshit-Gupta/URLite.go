package main

import(
	"fmt"
	"net/http"
)

//exatc fucntion handle to be made 
func printmsg(w http.ResponseWriter , r *http.Request){
    // write somehtingn to w 
}



func main(){
     
	http.HandleFunc("/",printmsg)

	

}


