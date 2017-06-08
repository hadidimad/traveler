package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var templates *template.Template

func InitRender() {
	templates = template.New("")
	FuncMap := template.FuncMap{
		"getyear": func(unixT int64) string {
			time := time.Unix(unixT, 0)
			return strconv.Itoa(time.Year())
		},
		"getmonth": func(unixT int64) string {
			time := time.Unix(unixT, 0)
			return strconv.Itoa(int(time.Month()))
		},
		"getday": func(unixT int64) string {
			time := time.Unix(unixT, 0)
			return strconv.Itoa(time.Day())
		},
		"getmin": func(unixT int64) string {
			time := time.Unix(unixT, 0)
			return strconv.Itoa(time.Minute())
		},
		"gethour": func(unixT int64) string {
			time := time.Unix(unixT, 0)
			return strconv.Itoa(time.Hour())
		},
		"getTravelsImage": func(id int) string {
			image, err := ioutil.ReadFile("travels/" + strconv.Itoa(id) + "/image.")
			if err != nil {
				fmt.Println(err)
			}
			ioutil.WriteFile("./statics/travels/"+strconv.Itoa(id), image, 0644)
			return "/statics/travels/" + strconv.Itoa(id)
		},
		"getUserImage": func(id int) string {
			onuser := getUserByID(id)
			image, err := ioutil.ReadFile("./userImages/" + onuser.Username)
			if err != nil {
				fmt.Println(err)
				image, _ := ioutil.ReadFile("./userImages/noimage")
				ioutil.WriteFile("./statics/users/"+onuser.Username, image, 0644)
				return "/statics/users/" + onuser.Username
			}
			ioutil.WriteFile("./statics/users/"+onuser.Username, image, 0644)
			return "/statics/users/" + onuser.Username
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
		"getTravelLikeLink": func(id int) string {
			return "/travellike?travel=" + strconv.Itoa(id)
		},
		"getTravelUnLikeLink": func(id int) string {
			return "/travelunlike?travel=" + strconv.Itoa(id)
		},
		"getUserInfoLink": func(username string) string {
			db, _ := sql.Open("sqlite3", "./database/database")
			rows, _ := db.Query("SELECT uid FROM userinfo WHERE username='" + username + "';")
			var id int
			for rows.Next() {
				rows.Scan(&id)
			}
			return "/userinfo?user=" + strconv.Itoa(id)
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
