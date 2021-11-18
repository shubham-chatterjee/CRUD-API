package main

import (
	"encoding/json"
	"fmt"
	"log"
	db "module/database"
	"net/http"
	"strconv"

	snowflake "github.com/godruoyi/go-snowflake"

	routes "github.com/gorilla/mux"
)

func create(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var person db.Person
	json.NewDecoder(request.Body).Decode(&person)
	person.Password = strconv.FormatUint(snowflake.ID(), 10)
	db.Add(person)
	json.NewEncoder(response).Encode(person)
}

func retrieve(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	values := routes.Vars(request)
	person, err := db.Retrieve(values["id"])
	if err != nil {
		log.Fatalln(err.Error())
	}
	json.NewEncoder(response).Encode(person)

}

func all(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	population, err := db.All()
	if err != nil {
		log.Fatalln(err.Error())
	}
	json.NewEncoder(response).Encode(population)
}

func update(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	params := routes.Vars(request)
	var person db.Person
	json.NewDecoder(request.Body).Decode(&person)
	db.Update(params["id"], person)
	person, _ = db.Retrieve(params["id"])
	json.NewEncoder(response).Encode(person)
}

func remove(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	params := routes.Vars(request)
	db.Delete(params["id"])
	population, err := db.All()
	if err != nil {
		log.Fatalln(err.Error())
	}
	json.NewEncoder(response).Encode(population)
}

func main() {
	router := routes.NewRouter()
	router.HandleFunc("/app/person", create).Methods("POST")
	router.HandleFunc("/app/people", all).Methods("GET")
	router.HandleFunc("/app/person/{id}", retrieve).Methods("GET")
	router.HandleFunc("/app/person/{id}", remove).Methods("DELETE")
	router.HandleFunc("/app/person/{id}", update).Methods("PUT")
	fmt.Print("Listening at port :8080")
	log.Fatalln(http.ListenAndServe(":8080", router).Error())
}
