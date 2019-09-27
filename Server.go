package main
import (
	"strconv"
	
	"encoding/json"
	"fmt"
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	jwt "github.com/dgrijalva/jwt-go"
	
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var mySigningKey=[]byte("supersecretkey")
var key1="admin"
var key2="admin123"

type User struct {
	//Gorm.model (ID,CreateTime,UpdateTime,DeleteTime)Tutan bir struct
	ID uint
	Name  string
	Surname string
	Age int
}

func homePage(w http.ResponseWriter,r *http.Request) {
	fmt.Fprintf(w, "Wellcome To Home Page")
}

func allUsers(w http.ResponseWriter, r *http.Request) {
	//gorm.open
	//:= tanımla ve ata
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	//defer fonksiyon bitmeden önce çalışır
	defer db.Close()

	var users []User
	db.Find(&users)
	fmt.Println("{}", users)

	//json formatına kodlanması
	json.NewEncoder(w).Encode(users)
}

func oneUser(w http.ResponseWriter, r *http.Request)  {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	id := vars["id"]

	var user User
	db.Where("id = ?",id).Find(&user)
	fmt.Println("",user)
	json.NewEncoder(w).Encode(user)
}

func newUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New User Endpoint Hit")

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]
	surname := vars["surname"]
	age := vars["age"]

	fmt.Println(name)
	fmt.Println(surname)
	fmt.Println(age)


	i1, err := strconv.Atoi(age)

	db.Create(&User{Name: name, Surname: surname, Age: i1})
	fmt.Fprintf(w, "New User Successfully Created")
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	id := vars["id"]

	var user User
	db.Where("id = ?", id).Find(&user)
	db.Delete(&user)

	fmt.Fprintf(w, "Successfully Deleted User")
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	id := vars["id"]
	surname := vars["surname"]

	var user User
	db.Where("id = ?", id).Find(&user)

	user.Surname = surname

	db.Save(&user)
	fmt.Fprintf(w, "Successfully Updated User")
}


func isAuthorized(endpoint func(http.ResponseWriter,*http.Request)) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter,  r *http.Request){

		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0],func(token *jwt.Token)(interface{},error){
				if _, ok:= token.Method.(*jwt.SigningMethodHMAC); !ok{
					return nil,fmt.Errorf("That was an error")
				}
				claims:=token.Claims.(jwt.MapClaims)
				if !(claims["userName"]==key1 && claims["userPassword"]==key2){
					return nil,fmt.Errorf("Wrong Authentication")
				}
				return mySigningKey,  nil
				
			})
			if err!=nil{
				fmt.Fprintf(w,err.Error())
			}
			if token.Valid{
				endpoint(w,r)
			}
		}else{
			fmt.Fprintf(w, "Not Authorized")
		}

	})
}


func handleRequest()  {
	
	myRouter := mux.NewRouter().StrictSlash(true)
	
	myRouter.Handle("/{userName}/{userPassword}/",isAuthorized(homePage)).Methods("GET")
	myRouter.Handle("/{userName}/{userPassword}/users", isAuthorized(allUsers)).Methods("GET")
	myRouter.Handle("/{userName}/{userPassword}/user/{id}", isAuthorized(oneUser)).Methods("GET")
	myRouter.Handle("/{userName}/{userPassword}/user/{id}", isAuthorized(deleteUser)).Methods("DELETE")
	myRouter.Handle("/{userName}/{userPassword}/user/{id}/{surname}", isAuthorized(updateUser)).Methods("PATCH")
	myRouter.Handle("/{userName}/{userPassword}/user/{name}/{surname}/{age}", isAuthorized(newUser)).Methods("POST")
	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func initialMigration() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema(Eğer yoksa ilgili tabloları modele bakarak oluşturucak güncelleme varsa yapıcak)
	db.AutoMigrate(&User{})
}	

func main()  {
	fmt.Println("My Server")
	
	initialMigration() 
	handleRequest()
}