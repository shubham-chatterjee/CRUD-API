package db

import (
	"encoding/json"
	"errors"
	_ "fmt"
	"io/ioutil"
	_ "log"
	"os"
)

type Person struct {
	ID        string `json:"id"`
	Password string `json:"password"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
}

func Retrieve(id string) (Person, error) {
	file, err := os.Open("database/db.json")
	if err != nil {
		return Person{}, err
	}
	decoder := json.NewDecoder(file)
	decoder.Token()
	var person = new(Person)
	for decoder.More() {
		decoder.Decode(person)
		if person.ID == id {
			return *person, nil
		}
	}
	return Person{}, errors.New("Person not found")
}

func All() ([]Person, error) {
	var population []Person
	data, err := ioutil.ReadFile("database/db.json")
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &population)
	return population, nil
}

func Update(id string, person Person) error {
	population, err := All() 
	if err != nil {
		return err 
	}  
	for index, p := range population {
		if p.ID == id { 
			person.ID = p.ID 
			person.Password = p.Password
			population[index] = person 
			break 
		}
	}
	file, err := os.OpenFile("database/db.json", os.O_CREATE, os.ModePerm)
	json.NewEncoder(file).Encode(population)
	return err
}

func Add(person Person) error {
	data, err := ioutil.ReadFile("database/db.json")
	if err != nil {
		return err
	}
	var population []Person
	json.Unmarshal(data, &population)
	population = append(population, person)
	file, err := os.OpenFile("database/db.json", os.O_CREATE, os.ModePerm)
	json.NewEncoder(file).Encode(population)
	return err
}

func Delete(id string) error {
	person, _ := Retrieve(id)
	population, _ := All()
	var indices []int
	for index, value := range population {
		if value.ID == person.ID {
			indices = append(indices, index)
		}
	}
	if len(indices) == 0 {
		return errors.New("Person not found")
	}
	for _, index := range indices {
		population = append(population[:index], population[index+1:]...)
	}
	file, err := os.Create("database/db.json")
	json.NewEncoder(file).Encode(population)
	return err
}
