package main

import (
	"encoding/json"
	"io/ioutil"
)

type user struct {
	ID       int
	Username string
	Password string
	Email    string
}

var users = []user{}

func updateUsers() {
	bytes, _ := ioutil.ReadFile("users.json")
	_ = json.Unmarshal(bytes, &users)
}
