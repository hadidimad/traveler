package main

import (
	"crypto/sha512"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"strings"

	"encoding/json"

	_ "github.com/mattn/go-sqlite3"
)

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	cookie, err := r.Cookie("User_Cookie")
	var onuser user
	if err == nil {
		userID, _ := strconv.Atoi(cookie.Value)
		onuser = getUserByID(userID)
		m["username"] = onuser.Username
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
	db, _ := sql.Open("sqlite3", "./database/database")
	rows, _ := db.Query("SELECT * FROM userinfo WHERE username='" + r.FormValue("username") + "'")
	for rows.Next() {
		rows.Scan(&onuser.ID, &onuser.Username, &onuser.Password, &onuser.Email, &onuser.likedTravels)
		passHash := sha512.New()
		io.WriteString(passHash, r.FormValue("password"))
		if onuser.Username == r.FormValue("username") && onuser.Password == string(passHash.Sum(nil)) {
			finded = true
			break
		}
	}
	rows.Close()
	if finded {
		cookie := &http.Cookie{
			Name:   "User_Cookie",
			Value:  strconv.Itoa(onuser.ID),
			MaxAge: 0,
		}
		http.SetCookie(w, cookie)
		m["username"] = onuser.Username
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
	if err == "invalidUsername" {
		m["invalidUsername"] = true
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
	if strings.ContainsAny(r.FormValue("username"), "/=*'><") {
		formIsValid = false
		http.Redirect(w, r, "/signup?err=invalidUsername", http.StatusFound)
	}
	if !(r.FormValue("password") == r.FormValue("password-repeat")) {
		formIsValid = false
		http.Redirect(w, r, "/signup?err=passnotmatch", http.StatusFound)
	}
	db, _ := sql.Open("sqlite3", "./database/database")
	if formIsValid {
		rows, _ := db.Query("SELECT * FROM userinfo WHERE username='" + r.FormValue("username") + "';")
		for rows.Next() {
			formIsValid = false
			http.Redirect(w, r, "/signup?err=takenusername", http.StatusFound)
		}
	}
	if formIsValid {
		rows, _ := db.Query("SELECT * FROM userinfo WHERE email='" + r.FormValue("email") + "';")
		for rows.Next() {
			formIsValid = false
			http.Redirect(w, r, "/signup?err=takenemail", http.StatusFound)
		}
	}
	if formIsValid {
		var temp = []int{}
		tempstr, _ := json.Marshal(temp)
		passHash := sha512.New()
		io.WriteString(passHash, r.FormValue("password"))

		stmt, _ := db.Prepare("INSERT INTO userinfo(username,password,email,likedTravels) values(?,?,?,?)")
		stmt.Exec(r.FormValue("username"), string(passHash.Sum(nil)), r.FormValue("email"), string(tempstr))
		r.ParseMultipartForm(32 << 20)
		file, _, err := r.FormFile("Image")
		if err == nil {
			defer file.Close()
			f, err := os.OpenFile("./userImages/"+r.FormValue("username"), os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer f.Close()
			io.Copy(f, file)
		}

		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func travelsGetHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	var onuser user
	cookie, err := r.Cookie("User_Cookie")
	if err == nil {
		userID, _ := strconv.Atoi(cookie.Value)
		onuser = getUserByID(userID)
		m["username"] = onuser.Username
	}
	mytravels := []travel{}
	db, _ := sql.Open("sqlite3", "./database/database")
	rows, _ := db.Query("SELECT * FROM travels")
	var temptravel travel
	for rows.Next() {

		rows.Scan(&temptravel.ID, &temptravel.Name, &temptravel.Start, &temptravel.End, &temptravel.Company, &temptravel.How, &temptravel.Description, &temptravel.ShareBy, &temptravel.Date, &temptravel.Likes)
		mytravels = append(mytravels, temptravel)
	}
	m["travels"] = mytravels
	Render(w, "travels", m)
}

func userGetHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	cookie, err := r.Cookie("User_Cookie")
	var onuser user
	if err == nil {
		userID, _ := strconv.Atoi(cookie.Value)
		onuser = getUserByID(userID)
		m["username"] = onuser.Username
		m["onuser"] = onuser
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	mytravels := []travel{}
	db, _ := sql.Open("sqlite3", "./database/database")
	rows, _ := db.Query("SELECT * FROM travels WHERE shareBy='" + onuser.Username + "';")
	var temptravel travel
	for rows.Next() {

		rows.Scan(&temptravel.ID, &temptravel.Name, &temptravel.Start, &temptravel.End, &temptravel.Company, &temptravel.How, &temptravel.Description, &temptravel.ShareBy, &temptravel.Date, &temptravel.Likes)
		mytravels = append(mytravels, temptravel)
	}
	m["usertravels"] = mytravels
	Render(w, "user", m)
}

func userInfoGetHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	var onuser user
	cookie, err := r.Cookie("User_Cookie")
	if err == nil {
		userID, _ := strconv.Atoi(cookie.Value)
		onuser = getUserByID(userID)
		m["username"] = onuser.Username
	}
	userIDstr := r.URL.Query().Get("user")
	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		fmt.Println("invlaid user id")
		http.Redirect(w, r, "/", http.StatusFound)
	}
	userInfo := getUserByID(userID)
	m["userInfo"] = userInfo
	mytravels := []travel{}
	db, _ := sql.Open("sqlite3", "./database/database")
	rows, _ := db.Query("SELECT * FROM travels WHERE shareBy='" + userInfo.Username + "';")
	var temptravel travel
	for rows.Next() {

		rows.Scan(&temptravel.ID, &temptravel.Name, &temptravel.Start, &temptravel.End, &temptravel.Company, &temptravel.How, &temptravel.Description, &temptravel.ShareBy, &temptravel.Date, &temptravel.Likes)
		mytravels = append(mytravels, temptravel)
	}
	m["usertravels"] = mytravels
	Render(w, "userinfo", m)
}

func userEditGetHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	var onuser user
	cookie, err := r.Cookie("User_Cookie")
	if err == nil {
		userID, _ := strconv.Atoi(cookie.Value)
		onuser = getUserByID(userID)
		m["username"] = onuser.Username
		m["onuser"] = onuser
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
	var onuser user
	if err == nil {
		userID, _ := strconv.Atoi(cookie.Value)
		onuser = getUserByID(userID)
		m["username"] = onuser.Username
		m["onuser"] = onuser
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
	var usernameData, emailData string
	db, _ := sql.Open("sqlite3", "./database/database")
	rows, _ := db.Query("SELECT username FROM userinfo WHERE username='" + r.FormValue("username") + "';")
	for rows.Next() {
		rows.Scan(&usernameData)
		if usernameData != onuser.Username {
			http.Redirect(w, r, "/useredit?err=takenusername", http.StatusFound)
			formIsValid = false
		}
	}
	rows, _ = db.Query("SELECT email FROM userinfo WHERE email='" + r.FormValue("email") + "';")
	for rows.Next() {
		rows.Scan(&emailData)
		if emailData != onuser.Email {
			http.Redirect(w, r, "/useredit?err=takenemail", http.StatusFound)
			formIsValid = false
		}
	}
	if formIsValid {
		db, _ := sql.Open("sqlite3", "./database/database")
		stmt, _ := db.Prepare("UPDATE userinfo SET username=? , password=? , email=? WHERE uid=?;")
		stmt.Exec(r.FormValue("username"), r.FormValue("password"), r.FormValue("email"), onuser.ID)
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("Image")
		if err == nil {
			if handler.Filename != "" {
				defer file.Close()
				f, err := os.OpenFile("./userImages/"+r.FormValue("username"), os.O_WRONLY|os.O_CREATE, 0666)
				if err != nil {
					fmt.Println(err)
					return
				}
				defer f.Close()
				io.Copy(f, file)
			}
		}
		http.Redirect(w, r, "/user", http.StatusFound)
	}
	Render(w, "userEdit", m)
}

func userDeleteGetHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	cookie, err := r.Cookie("User_Cookie")
	var userID int
	var onuser user
	if err == nil {
		userID, _ = strconv.Atoi(cookie.Value)
		onuser = getUserByID(userID)
		m["username"] = onuser.Username
		m["onuser"] = onuser
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	os.Remove("userImages/" + onuser.Username)
	os.Remove("./statics/users/" + onuser.Username)
	db, _ := sql.Open("sqlite3", "./database/database")
	stmt, _ := db.Prepare("delete	from	userinfo	where	uid=?")
	stmt.Exec(onuser.ID)
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
	var onuser user
	var isLogin bool
	isLogin = false
	if err == nil {
		isLogin = true
		userID, _ = strconv.Atoi(cookie.Value)
		onuser = getUserByID(userID)
		m["username"] = onuser.Username
	}
	travelIDstr := r.URL.Query().Get("travel")
	var temptravel travel
	db, _ := sql.Open("sqlite3", "./database/database")
	rows, _ := db.Query("SELECT * FROM travels WHERE id=" + travelIDstr + ";")
	for rows.Next() {
		rows.Scan(&temptravel.ID, &temptravel.Name, &temptravel.Start, &temptravel.End, &temptravel.Company, &temptravel.How, &temptravel.Description, &temptravel.ShareBy, &temptravel.Date, &temptravel.Likes)
	}
	travelLiked := false
	var travelsID []int
	json.Unmarshal([]byte(onuser.likedTravels), &travelsID)
	for i := 0; i < len(travelsID); i++ {
		if travelsID[i] == temptravel.ID {
			travelLiked = true
			break
		}
	}
	m["liked"] = travelLiked
	m["isLogin"] = isLogin
	m["travel"] = temptravel
	Render(w, "travelinfo", m)
}

func travelLikeGetHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	cookie, err := r.Cookie("User_Cookie")
	var userID int
	var onuser user
	if err == nil {
		userID, _ = strconv.Atoi(cookie.Value)
		onuser = getUserByID(userID)
		m["username"] = onuser.Username
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	travelIDstr := r.URL.Query().Get("travel")
	travelID, _ := strconv.Atoi(travelIDstr)
	var likedID []int
	json.Unmarshal([]byte(onuser.likedTravels), &likedID)
	likedID = append(likedID, travelID)
	jsonValue, _ := json.Marshal(likedID)
	db, _ := sql.Open("sqlite3", "./database/database")
	stmt, _ := db.Prepare("UPDATE userinfo SET likedTravels=? WHERE uid=" + strconv.Itoa(onuser.ID) + ";")
	stmt.Exec(jsonValue)
	rows, _ := db.Query("SELECT likes FROM travels WHERE id=" + travelIDstr + ";")
	var likes int
	for rows.Next() {
		rows.Scan(&likes)
	}
	likes++
	stmt, _ = db.Prepare("UPDATE travels SET likes=? WHERE id=" + travelIDstr + ";")
	stmt.Exec(likes)
	http.Redirect(w, r, "/travelinfo?travel="+travelIDstr, http.StatusFound)
}

func travelUnLikeGetHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	cookie, err := r.Cookie("User_Cookie")
	var userID int
	var onuser user
	if err == nil {
		userID, _ = strconv.Atoi(cookie.Value)
		onuser = getUserByID(userID)
		m["username"] = onuser.Username
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	travelIDstr := r.URL.Query().Get("travel")
	travelID, _ := strconv.Atoi(travelIDstr)
	var likedID []int
	json.Unmarshal([]byte(onuser.likedTravels), &likedID)
	i := 0
	for i = 0; i < len(likedID); i++ {
		if likedID[i] == travelID {
			break
		}
	}
	likedID = append(likedID[:i], likedID[i+1:]...)
	jsonValue, _ := json.Marshal(likedID)
	db, _ := sql.Open("sqlite3", "./database/database")
	stmt, _ := db.Prepare("UPDATE userinfo SET likedTravels=? WHERE uid=" + strconv.Itoa(onuser.ID) + ";")
	stmt.Exec(jsonValue)
	rows, _ := db.Query("SELECT likes FROM travels WHERE id=" + travelIDstr + ";")
	var likes int
	for rows.Next() {
		rows.Scan(&likes)
	}
	likes--
	stmt, _ = db.Prepare("UPDATE travels SET likes=? WHERE id=" + travelIDstr + ";")
	stmt.Exec(likes)
	http.Redirect(w, r, "/travelinfo?travel="+travelIDstr, http.StatusFound)
}

func newTravelGetHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	cookie, err := r.Cookie("User_Cookie")
	var userID int
	var onuser user
	if err == nil {
		userID, _ = strconv.Atoi(cookie.Value)
		onuser = getUserByID(userID)
		m["username"] = onuser.Username
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	Render(w, "newtravel", m)
}

func newTravelPostHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	cookie, err := r.Cookie("User_Cookie")
	var userID int
	var onuser user
	if err == nil {
		userID, _ = strconv.Atoi(cookie.Value)
		onuser = getUserByID(userID)
		m["username"] = onuser.Username
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
		dateyear, _ := strconv.Atoi(r.FormValue("Date-year"))
		dateday, _ := strconv.Atoi(r.FormValue("Date-day"))
		datemonth, _ := strconv.Atoi(r.FormValue("Date-month"))
		datehour, _ := strconv.Atoi(r.FormValue("Time-hour"))
		dateminute, _ := strconv.Atoi(r.FormValue("Time-minute"))
		loc := time.Now().Location()
		time := time.Date(dateyear, time.Month(datemonth), dateday, datehour, dateminute, 6, 0, loc)

		tempTravel.Date = time.Unix()
		tempTravel.How = r.FormValue("how")
		tempTravel.Company = r.FormValue("company")
		tempTravel.Description = r.FormValue("description")
		tempTravel.ShareBy = getUserByID(onuser.ID).Username
		db, _ := sql.Open("sqlite3", "./database/database")
		stmt, _ := db.Prepare("INSERT INTO travels(name,start,end,company,how,description,shareby,date,likes) values(?,?,?,?,?,?,?,?,?)")
		res, _ := stmt.Exec(tempTravel.Name, tempTravel.Start, tempTravel.End, tempTravel.Company, tempTravel.How, tempTravel.Description, tempTravel.ShareBy, tempTravel.Date, 0)
		travelID, _ := res.LastInsertId()
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("image")
		if err == nil {
			if handler.Filename != "" {
				defer file.Close()
				os.MkdirAll("./travels/"+strconv.Itoa(int(travelID)), 0777)
				f, err := os.Create("./travels/" + strconv.Itoa(int(travelID)) + "/image.")
				if err != nil {
					fmt.Println(err)
				}
				defer f.Close()
				io.Copy(f, file)
			}
		} else {
			os.MkdirAll("./travels/"+strconv.Itoa(int(travelID)), 0777)
			f, err := os.Create("./travels/" + strconv.Itoa(int(travelID)) + "/image.")
			if err != nil {
				fmt.Println(err)
			}
			f.Close()
		}
		http.Redirect(w, r, "/user", 302)
	}
}
func deleteTravelGetHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	cookie, err := r.Cookie("User_Cookie")
	var userID int
	var onuser user
	if err == nil {
		userID, _ = strconv.Atoi(cookie.Value)
		onuser = getUserByID(userID)
		m["username"] = onuser.Username
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	travelIDstr := r.URL.Query().Get("travel")
	travelID, _ := strconv.Atoi(travelIDstr)
	var temptravel travel
	db, _ := sql.Open("sqlite3", "./database/database")
	rows, _ := db.Query("SELECT * FROM travels WHERE id=" + travelIDstr + ";")
	for rows.Next() {
		rows.Scan(&temptravel.ID, &temptravel.Name, &temptravel.Start, &temptravel.End, &temptravel.Company, &temptravel.How, &temptravel.Description, &temptravel.ShareBy, &temptravel.Date, &temptravel.Likes)
	}
	stmt, _ := db.Prepare("DELETE FROM travels WHERE id=?")
	stmt.Exec(travelID)
	stmt.Close()
	db.Close()
	os.RemoveAll("./travels/" + strconv.Itoa(temptravel.ID))
	os.Remove("./statics/travels/" + strconv.Itoa(temptravel.ID))
	http.Redirect(w, r, "/user", 302)

}

func getUserByID(id int) user {
	db, _ := sql.Open("sqlite3", "./database/database")
	rows, _ := db.Query("SELECT * FROM userinfo WHERE uid=" + strconv.Itoa(id) + ";")
	var onuser user
	for rows.Next() {
		rows.Scan(&onuser.ID, &onuser.Username, &onuser.Password, &onuser.Email, &onuser.likedTravels)
	}
	return onuser
}
