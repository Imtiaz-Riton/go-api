package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

type Publication struct {
	Id 			int		`json:"id"`
	Title 		string	`json:"title"`
	Description string	`json:"description"`
}

type Author struct {
	Id 		int		`json:"id"`
	Name 	string 	`json:"name"`
}

var(
	author = Author{Id: 1, Name: "Imtiaz"}
	publication = Publication{Id: 1, Title: "go programming", Description: "go is easy"}
)
var db *sql.DB
var err error
func main()  {
	fmt.Println("hello")
	host := "localhost"
	dbPort := "5432"
	user := "postgres"
	dbName := "author_publication"
	password := "abc123"
	dbURI := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		host, user, password, dbName, dbPort)

	db, err = sql.Open("postgres", dbURI)
	if err != nil{
		panic(err)
	}else{
		fmt.Println("successfully connected database")
	}
	defer db.Close()

	_, err = db.Query(`INSERT INTO author(id, name)VALUES($1, $2);`, author.Id, author.Name)
	if err != nil{
		panic(err)
	}
	_, err = db.Query(`INSERT INTO publication(id, title, description)VALUES($1, $2, $3);`, publication.Id, publication.Title, publication.Description)
	if err != nil{
		panic(err)
	}
	handleRequests()
}

func handleRequests()  {
	const port string = ":8080"
	router := mux.NewRouter()
	router.HandleFunc("/", homePage)
	router.HandleFunc("/author", findAuthor)
	router.HandleFunc("/publication", findPublication)
	log.Fatal(http.ListenAndServe(port, router))
}

func homePage(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprintf(w, "welcome to home page!")
}

func findAuthor(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("content-type","application/json")
	rows, err := db.Query("SELECT * FROM author")
	if err != nil{
		panic(err)
	}
	var author1 = Author{}
	var authorArr = []Author{}
	for rows.Next(){
		rows.Scan(&author1.Id, &author1.Name)
		authorArr = append(authorArr, author1)
	}
	json.NewEncoder(w).Encode(&authorArr)
}

func findPublication(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("content-type","application/json")
	rows, err := db.Query("SELECT * FROM publication")
	if err != nil{
		panic(err)
	}
	var publication1 = Publication{}
	for rows.Next(){
		rows.Scan(&publication1.Id, &publication1.Title, &publication1.Description)
	}
	json.NewEncoder(w).Encode(&publication1)
}