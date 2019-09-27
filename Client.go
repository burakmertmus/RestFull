package main

import (
	"io/ioutil"
	"fmt"
	"net/http"
	"time"
	"log"
	
	"github.com/gorilla/mux"
	jwt "github.com/dgrijalva/jwt-go"
)

//var mySigningKey = os.Get("My_JWT_TOKEN")
var mySingningKey = []byte ("supersecretkey")

func homePage(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	userName := vars["userName"]
	userPassword := vars["userPassword"]

	validToken, err := GenereteJWT(userName,userPassword)
	if err!=nil {
		fmt.Fprintf(w, err.Error())
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET","http://localhost:8081/"+userName+"/"+userPassword+"/",nil)
	req.Header.Set("Token", validToken)
	res,err :=client.Do(req)
	if err!=nil{
		fmt.Fprintf(w,"Error: %s",err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err!=nil {
		fmt.Fprintf(w,err.Error())
	}
	
	fmt.Println(validToken)
	fmt.Fprintf(w,string(body))
}

func allUsers(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	userName := vars["userName"]
	userPassword := vars["userPassword"]

	validToken, err := GenereteJWT(userName,userPassword)
	if err!=nil {
		fmt.Fprintf(w, err.Error())
	}
	
	client := &http.Client{}
	req, _ := http.NewRequest("GET","http://localhost:8081/"+userName+"/"+userPassword+"/users",nil)
	req.Header.Set("Token", validToken)
	res,err :=client.Do(req)
	if err!=nil{
		fmt.Fprintf(w,"Error: %s",err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err!=nil {
		fmt.Fprintf(w,err.Error())
	}
	
	fmt.Println(validToken)
	fmt.Fprintf(w,string(body))
}
func newUser(w http.ResponseWriter, r *http.Request)  {

	
	vars := mux.Vars(r)
		userName := vars["userName"]
	userPassword := vars["userPassword"]
	name := vars["name"]
	surname := vars["surname"]
	age := vars["age"]


	validToken, err := GenereteJWT(userName,userPassword)
	if err!=nil {
		fmt.Fprintf(w, err.Error())
	}
	client := &http.Client{}
	req, _ := http.NewRequest("POST","http://localhost:8081/"+userName+"/"+userPassword+"/user/"+name+"/"+surname+"/"+age+"",nil)
	req.Header.Set("Token", validToken)
	res,err :=client.Do(req)
	if err!=nil{
		fmt.Fprintf(w,"Error: %s",err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err!=nil {
		fmt.Fprintf(w,err.Error())
	}
	
	fmt.Println(validToken)
	fmt.Fprintf(w,string(body))
}

func oneUser(w http.ResponseWriter, r *http.Request)  {

	
	vars := mux.Vars(r)
	id := vars["id"]
	userName := vars["userName"]
	userPassword := vars["userPassword"]
	

	validToken, err := GenereteJWT(userName,userPassword)
	if err!=nil {
		fmt.Fprintf(w, err.Error())
	}
	client := &http.Client{}
	req, _ := http.NewRequest("GET","http://localhost:8081/"+userName+"/"+userPassword+"/user/"+id+"",nil)
	req.Header.Set("Token", validToken)
	res,err :=client.Do(req)
	if err!=nil{
		fmt.Fprintf(w,"Error: %s",err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err!=nil {
		fmt.Fprintf(w,err.Error())
	}
	
	fmt.Println(validToken)
	fmt.Fprintf(w,string(body))
}

func deleteUser(w http.ResponseWriter, r *http.Request)  {

	
	vars := mux.Vars(r)
	id := vars["id"]
	userName := vars["userName"]
	userPassword := vars["userPassword"]

	validToken, err := GenereteJWT(userName,userPassword)
	if err!=nil {
		fmt.Fprintf(w, err.Error())
	}
	client := &http.Client{}
	req, _ := http.NewRequest("DELETE","http://localhost:8081/"+userName+"/"+userPassword+"/user/"+id+"",nil)
	req.Header.Set("Token", validToken)
	res,err :=client.Do(req)
	if err!=nil{
		fmt.Fprintf(w,"Error: %s",err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err!=nil {
		fmt.Fprintf(w,err.Error())
	}
	
	fmt.Println(validToken)
	fmt.Fprintf(w,string(body))
}

func updateUser(w http.ResponseWriter, r *http.Request)  {

	
	vars := mux.Vars(r)
	id := vars["id"]
	surname := vars["surname"]
	userName := vars["userName"]
	userPassword := vars["userPassword"]

	validToken, err := GenereteJWT(userName,userPassword)
	if err!=nil {
		fmt.Fprintf(w, err.Error())
	}
	client := &http.Client{}
	req, _ := http.NewRequest("PATCH","http://localhost:8081/"+userName+"/"+userPassword+"/user/"+id+"/"+surname+"",nil)
	req.Header.Set("Token", validToken)
	res,err :=client.Do(req)
	if err!=nil{
		fmt.Fprintf(w,"Error: %s",err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err!=nil {
		fmt.Fprintf(w,err.Error())
	}
	
	fmt.Println(validToken)
	fmt.Fprintf(w,string(body))
}



func GenereteJWT(userName string,userPassword string) (string,error)  {

	token:= jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"authorized":true,
		"userName":userName,
		"userPassword":userPassword,
		"exp":time.Now().Add(time.Minute * 30).Unix(),
	})

	tokenString, err := token.SignedString(mySingningKey)

	if err !=nil{
		fmt.Errorf("Something went wrong: s%",err.Error())
		return "",err
	}

	return tokenString,nil
}

func handleRequest()  {


	
	myRouter := mux.NewRouter().StrictSlash(true)
	
	myRouter.HandleFunc("/{userName}/{userPassword}/",homePage).Methods("GET")
	myRouter.HandleFunc("/{userName}/{userPassword}/users", allUsers).Methods("GET")
	myRouter.HandleFunc("/{userName}/{userPassword}/user/{id}", oneUser).Methods("GET")
	//myRouter.HandleFunc("/user/{name}/{userPassword}/user/{name}/{surname}/{age}", newUser).Methods("POST") SUNUMDAKİ HALİ
	myRouter.HandleFunc("/{userName}/{userPassword}/user/{name}/{surname}/{age}", newUser).Methods("POST")
	myRouter.HandleFunc("/{userName}/{userPassword}/user/{id}", deleteUser).Methods("DELETE")
	myRouter.HandleFunc("/{userName}/{userPassword}/user/{id}/{surname}", updateUser).Methods("PATCH")
	
	log.Fatal(http.ListenAndServe(":9001",myRouter))
}

func main()  {
	fmt.Println("My Client")	
	handleRequest()
}