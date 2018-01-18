package main

import (
    "fmt"
    "log"
    "net/http"
	"html/template"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
)

//Struct for Data to sent to the webpage
//An array of a struct containing a single element is for listing purposes in the webpage
type Index struct{
	DataSet []DataList
}

//Data Array
type DataList struct{
	Data string
}

//html template 
var indexTemplate = template.Must(template.ParseFiles("index.html"))

//Handler for the index page when the application first opens
func indexHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//fmt.Println(r.Form)
	if err := indexTemplate.Execute(w, getData()); err != nil {
		log.Println(err)
	}
}

//Handler for when the user inputs data
func dataInputHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//If no input from the user when clicking submit
	if (r.Form["input"]==nil){ 
		if err := indexTemplate.Execute(w, getData()); err != nil {
			log.Println(err)
		}
		return
	}
	addData(r.Form["input"][0])
	if err := indexTemplate.Execute(w, getData()); err != nil {
		log.Println(err)
	}
}

//retrieve data from mysql server
func getData() (Index){
	data :=Index{}
	//connect to local mysql database, replace username,password,and databasename
	db, err := sql.Open("mysql", "username:password@/databasename")
    if err != nil {
        panic(err.Error())
    }
	
	//retrieve data entries from table
	result, err := db.Query("SELECT * FROM dataentries")
	 if err != nil {
        panic(err.Error())
    }
	for result.Next(){
		var datavalue string
		err = result.Scan(&datavalue)
		if err != nil {
			panic(err.Error()) 
		}
		data.DataSet=append(data.DataSet,DataList{datavalue})
	}
	db.Close()
	return data
}

//add a piece of data to the mysql server when user clicks submit
func addData(data string){
	//connect to local sql database, replace username,password,and database name
	db, err := sql.Open("mysql", "username:password@/databasename")
    if err != nil {
        panic(err.Error())
    }
	
	//insert into database
	_, erri := db.Query("insert into dataentries values ('"+data+"')")
	if erri != nil {
        panic(err.Error())
    }
	db.Close()
}

func main() {	
	fmt.Println("Launching GoApp at http://localhost:8081/")
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/dataInput",dataInputHandler)
    err := http.ListenAndServe(":8081", nil) // setting listening port
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

