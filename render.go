package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var templates *template.Template

func InitRender() {
	templates = template.New("")
	FuncMap := template.FuncMap{
		"getTravelsImage": func(id int) string {
			image, err := ioutil.ReadFile(travels[id].Path + "/image.")
			if err != nil {
				fmt.Println(err)
			}
			ioutil.WriteFile("./statics/travels/"+travels[id].Name, image, 0644)
			return "/statics/travels/" + travels[id].Name
		},
		"getUserImage": func(id int) string {
			image, err := ioutil.ReadFile("./userImages/" + users[id].Username)
			if err != nil {
				fmt.Println(err)
				image, _ := ioutil.ReadFile("./userImages/noimage")
				ioutil.WriteFile("./statics/users/"+users[id].Username, image, 0644)
				return "/statics/users/" + users[id].Username
			}
			ioutil.WriteFile("./statics/users/"+users[id].Username, image, 0644)
			return "/statics/users/" + users[id].Username
		},
		"getTravelInfoLink": func(id int) string {
			return "/travelinfo?travel=" + strconv.Itoa(id)
		},
		"getTravelDeleteLink": func(id int) string {
			return "/deletetravel?travel=" + strconv.Itoa(id)
		},
		"getTravelEditLink": func(id int) string {
			return "/edittravel?travel=" + strconv.Itoa(id)
		},
	}
	templates.Funcs(FuncMap)
	err := filepath.Walk("./view", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".html") {
			_, err = templates.ParseFiles(path)
		}
		return err
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func Render(w io.Writer, name string, m interface{}) {
	err := templates.ExecuteTemplate(w, name, m)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
