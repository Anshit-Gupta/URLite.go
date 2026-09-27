package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	
	godotenv.Load()
    databaseURL := os.Getenv("DATABASE_URL")
	ctx := context.Background()
	var err error
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		fmt.Println("error creating database pool", err)
		return
	}
	fmt.Println("pool working ")
	defer pool.Close()
    exitCode := m.Run()  // actually runs all your Test... functions
    // optional cleanup here
    os.Exit(exitCode)
}


func TestGenerateShortCode(t *testing.T){
   code:=generateShortCode()
   expected_len:=6
   if len(code)!=expected_len{
	  t.Errorf("expected length %d, got %d", expected_len, len(code))
   }

   const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
   for _,char := range code {
	   if !strings.ContainsRune(charset, char){
        t.Errorf("unexpected char %q found in code", char)
	   }
   } 
}

func TestRedirect ( t *testing.T){
	//we  create a custom code-url pair and instert it into the db 
	testCode :="test123"  
	testURL := "https://anshitgupta.in"

	_,err := pool.Exec(context.Background(),"INSERT INTO urls (code ,long_url) VALUES($1 , $2)", testCode,testURL)
	if err!=nil{
		t.Fatalf("failed to insert test data : %v",err)
	}

	//CLEANUP : remove the row once this test is done : fail or pass
	defer pool.Exec(context.Background(),"DELETE FROM urls WHERE code = $1", testCode)

	//build fake incoming request - as if somone is visiting the shorturl 

	req :=httptest.NewRequest("GET","/"+testCode,nil)
	recorder :=httptest.NewRecorder()
   
	//call the handler directly 
    redirect(recorder,req)

	 //check 1 : did it respond with a redirect status code ?
	 if recorder.Code !=http.StatusFound{
		t.Errorf("expected status %d, got %d ", http.StatusFound , recorder.Code )

	 }
    //check 2 : did it redirect to the correct long URL?
	 if location:=recorder.Header().Get("Location"); location!=testURL{
		t.Errorf("expected redirect to %s , got  %s ",testURL,location)
	 }



}