package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type date struct {
	Day   int
	Month time.Month
	Year  int
}
type Time struct {
	Hour   int
	Minute int
}

type travel struct {
	ID          int
	Name        string
	Start       string
	End         string
	Time        Time
	Date        date
	Company     string
	How         string
	Description string
	Path        string
	ShareBy     string
}

var travels = []travel{}

type ByID []travel

func (this ByID) Len() int {
	return len(this)
}
func (this ByID) Less(i, j int) bool {
	return this[i].ID < this[j].ID
}
func (this ByID) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func updateTravels() {
	travels = nil
	err := filepath.Walk("./travels/", func(path string, info os.FileInfo, err error) error {
		if !strings.Contains(path, ".") {
			bytes, _ := ioutil.ReadFile(path + "/data.json")
			var tempTravel travel
			json.Unmarshal(bytes, &tempTravel)
			tempTravel.Path = path
			travels = append(travels, tempTravel)
		}
		return err
	})
	sort.Sort(ByID(travels))
	if err != nil {
		fmt.Println(err)
	}
}
