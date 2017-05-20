package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	cookie, err := r.Cookie("User_Cookie")
	if err == nil {
		userID, _ := strconv.Atoi(cookie.Value)
		m["username"] = users[userID].Username
	}

	Render(w, "index", m)
}

func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	if r.URL.Query().Get("err") == "invalid" {
		m["invalid"] = true
	}
	Render(w, "login", m)
}

func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	var onuser user
	var finded bool
	finded = false
	for i := 0; i < len(users); i++ {
		if r.FormValue("username") == users[i].Username && r.FormValue("password") == users[i].Password {
			onuser = users[i]
			finded = true
			break
		}
	}
	if finded {
		cookie := &http.Cookie{
			Name:   "User_Cookie",
			Value:  strconv.Itoa(onuser.ID),
			MaxAge: 0,
		}
		http.SetCookie(w, cookie)
		m["username"] = users[onuser.ID].Username
		http.Redirect(w, r, "/user", http.StatusFound)
	} else {
		http.Redirect(w, r, "/login?err=invalid", http.StatusFound)
	}
}

func logoutGetHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("User_Cookie")
	if err == nil {
		cookie.MaxAge = -1
		cookie.Value = "logout"
		http.SetCookie(w, cookie)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func signupGetHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	err := r.URL.Query().Get("err")
	if err == "passnotmatch" {
		m["passnotmatch"] = true
	}
	if err == "takenusername" {
		m["takenusername"] = true
	}
	if err == "takenemail" {
		m["takenemail"] = true
	}
	if err == "emptyfield" {
		m["emptyfield"] = true
	}

	Render(w, "signup", m)
}

func signupPostHandler(w http.ResponseWriter, r *http.Request) {
	var formIsValid bool
	formIsValid = true
	if r.FormValue("password") == "" || r.FormValue("username") == "" || r.FormValue("email") == "" {
		formIsValid = false
		http.Redirect(w, r, "/signup?err=emptyfield", http.StatusFound)
	}
	if !(r.FormValue("password") == r.FormValue("password-repeat")) {
		formIsValid = false
		http.Redirect(w, r, "/signup?err=passnotmatch", http.StatusFound)
	}
	for i := 0; i < len(users); i++ {
		if users[i].Username == r.FormValue("username") {
			formIsValid = false
			http.Redirect(w, r, "/signup?err=takenusername", http.StatusFound)
			break
		}
		if users[i].Email == r.FormValue("email") {
			formIsValid = false
			http.Redirect(w, r, "/signup?err=takenemail", http.StatusFound)
		}
	}
	if formIsValid {
		var onuser user
		onuser.Username = r.FormValue("username")
		onuser.Password = r.FormValue("password")
		onuser.Email = r.FormValue("email")
		onuser.ID = len(users)
		r.ParseMultipartForm(32 << 20)
		file, _, err := r.FormFile("Image")
		if err == nil {
			defer file.Close()
			f, err := os.OpenFile("./userImages/"+onuser.Username, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer f.Close()
			io.Copy(f, file)
		}

		users = append(users, onuser)
		bytes, _ := json.Marshal(users)
		_ = ioutil.WriteFile("users.json", bytes, 0644)
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func travelsGetHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	cookie, err := r.Cookie("User_Cookie")
	if err == nil {
		userID, _ := strconv.Atoi(cookie.Value)
		m["username"] = users[userID].Username
	}
	m["travels"] = travels
	Render(w, "travels", m)
}

func userGetHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	cookie, err := r.Cookie("User_Cookie")
	var userID int
	if err == nil {
		userID, _ = strconv.Atoi(cookie.Value)
		m["username"] = users[userID].Username
		m["onuser"] = users[userID]
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	var userTravels = []*travel{}
	for i := 0; i < len(travels); i++ {
		if travels[i].ShareBy == users[userID].Username {
			userTravels = append(userTravels, &travels[i])
		}
	}
	m["usertravels"] = userTravels
	Render(w, "user", m)
}

func userEditGetHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	cookie, err := r.Cookie("User_Cookie")
	if err == nil {
		userID, _ := strconv.Atoi(cookie.Value)
		m["username"] = users[userID].Username
		m["onuser"] = users[userID]
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	error := r.URL.Query().Get("err")
	if error == "passnotmatch" {
		m["passnotmatch"] = true
	}
	if error == "takenusername" {
		m["takenusername"] = true
	}
	if error == "takenemail" {
		m["takenemail"] = true
	}
	if error == "emptyfield" {
		m["emptyfield"] = true
	}
	Render(w, "userEdit", m)
}

func userEditPostHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	cookie, err := r.Cookie("User_Cookie")
	var userID int
	if err == nil {
		userID, _ = strconv.Atoi(cookie.Value)
		m["username"] = users[userID].Username
		m["onuser"] = users[userID]
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	var formIsValid bool
	formIsValid = true
	if r.FormValue("password") == "" || r.FormValue("username") == "" || r.FormValue("email") == "" {
		formIsValid = false
		http.Redirect(w, r, "/useredit?err=emptyfield", http.StatusFound)
	}
	if !(r.FormValue("password") == r.FormValue("password-repeat")) {
		formIsValid = false
		http.Redirect(w, r, "/useredit?err=passnotmatch", http.StatusFound)
	}
	for i := 0; i < len(users); i++ {
		if i != userID {
			if users[i].Username == r.FormValue("username") {
				formIsValid = false
				http.Redirect(w, r, "/useredit?err=takenusername", http.StatusFound)
				break
			}
			if users[i].Email == r.FormValue("email") {
				formIsValid = false
				http.Redirect(w, r, "/useredit?err=takenemail", http.StatusFound)
			}
		}
	}
	if formIsValid {
		users[userID].Username = r.FormValue("username")
		users[userID].Password = r.FormValue("password")
		users[userID].Email = r.FormValue("email")

		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("Image")
		if err == nil {
			if handler.Filename != "" {
				defer file.Close()
				f, err := os.OpenFile("./userImages/"+users[userID].Username, os.O_WRONLY|os.O_CREATE, 0666)
				if err != nil {
					fmt.Println(err)
					return
				}
				defer f.Close()
				io.Copy(f, file)
			}
		}

		bytes, _ := json.Marshal(users)
		_ = ioutil.WriteFile("users.json", bytes, 0644)

		http.Redirect(w, r, "/user", http.StatusFound)
	}
	Render(w, "userEdit", m)
}

func userDeleteGetHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	cookie, err := r.Cookie("User_Cookie")
	var userID int
	if err == nil {
		userID, _ = strconv.Atoi(cookie.Value)
		m["username"] = users[userID].Username
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	os.Remove("userImages/" + users[userID].Username)
	users = append(users[:userID], users[userID+1:]...)
	for i := 0; i < len(users); i++ {
		users[i].ID = i
	}
	bytes, _ := json.Marshal(users)
	ioutil.WriteFile("users.json", bytes, 0644)
	if err == nil {
		cookie.MaxAge = -1
		cookie.Value = "logout"
		http.SetCookie(w, cookie)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func travelInfoGetHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	cookie, err := r.Cookie("User_Cookie")
	var userID int
	if err == nil {
		userID, _ = strconv.Atoi(cookie.Value)
		m["username"] = users[userID].Username
	}
	travelIDstr := r.URL.Query().Get("travel")
	travelID, _ := strconv.Atoi(travelIDstr)
	if travelID >= len(travels) {
		http.Redirect(w, r, "/travels", http.StatusFound)
	} else {
		m["travel"] = travels[travelID]
		Render(w, "travelinfo", m)
	}
}

func newTravelGetHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	cookie, err := r.Cookie("User_Cookie")
	var userID int
	if err == nil {
		userID, _ = strconv.Atoi(cookie.Value)
		m["username"] = users[userID].Username
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	Render(w, "newtravel", m)
}

func newTravelPostHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	cookie, err := r.Cookie("User_Cookie")
	var userID int
	if err == nil {
		userID, _ = strconv.Atoi(cookie.Value)
		m["username"] = users[userID].Username
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	var formValid bool
	formValid = true

	if formValid {
		var tempTravel travel
		tempTravel.Name = r.FormValue("name")
		tempTravel.Start = r.FormValue("start")
		tempTravel.End = r.FormValue("end")
		tempTravel.Date.Day, _ = strconv.Atoi(r.FormValue("Date-day"))
		month, _ := strconv.Atoi(r.FormValue("Date-month"))
		tempTravel.Date.Month = time.Month(month)
		tempTravel.Date.Year, _ = strconv.Atoi(r.FormValue("Date-day"))
		tempTravel.Time.Hour, _ = strconv.Atoi(r.FormValue("Time-hour"))
		tempTravel.Time.Minute, _ = strconv.Atoi(r.FormValue("Time-minute"))
		tempTravel.How = r.FormValue("how")
		tempTravel.Company = r.FormValue("company")
		tempTravel.Description = r.FormValue("description")
		tempTravel.ID = len(travels)
		tempTravel.ShareBy = users[userID].Username
		os.Mkdir("./travels/"+tempTravel.Name, os.ModePerm)
		bytes, _ := json.Marshal(tempTravel)
		ioutil.WriteFile("./travels/"+tempTravel.Name+"/data.json", bytes, 0644)
		fmt.Println(tempTravel)
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("image")
		if err == nil {
			if handler.Filename != "" {
				defer file.Close()
				f, err := os.OpenFile("./travels/"+tempTravel.Name+"/image.", os.O_WRONLY|os.O_CREATE, 0666)
				if err != nil {
					fmt.Println(err)
					return
				}
				defer f.Close()
				io.Copy(f, file)
			}
		} else {
			f, _ := os.OpenFile("./travels/"+tempTravel.Name+"/image.", os.O_WRONLY|os.O_CREATE, 0666)
			f.Close()
		}
		updateTravels()
		http.Redirect(w, r, "/user", 302)
	}
}
func deleteTravelGetHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	cookie, err := r.Cookie("User_Cookie")
	var userID int
	if err == nil {
		userID, _ = strconv.Atoi(cookie.Value)
		m["username"] = users[userID].Username
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	travelIDstr := r.URL.Query().Get("travel")
	travelID, _ := strconv.Atoi(travelIDstr)
	os.RemoveAll(travels[travelID].Path)
	travels = append(travels[:travelID], travels[travelID+1:]...)
	for i := 0; i < len(travels); i++ {
		travels[i].ID = i
		bytes, _ := json.Marshal(travels[i])
		ioutil.WriteFile(travels[i].Path+"/data.json", bytes, 0644)
	}
	updateTravels()
	http.Redirect(w, r, "/user", 302)

}
