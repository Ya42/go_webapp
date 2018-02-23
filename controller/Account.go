package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ya42/go_webapp/model"
	"github.com/ya42/go_webapp/model/message"

	"github.com/ya42/go_webapp/service"
	"github.com/ya42/go_webapp/common/passhash"
	"github.com/ya42/go_webapp/common/session"

	"github.com/gorilla/sessions"
	"github.com/josephspurrier/csrfbanana"
)

const (
	// Name of the session variable that tracks login attempts
	sessLoginAttempt = "login_attempt"
	connStr = "db connection string";
)

var userService *service.UserService

// loginAttempt increments the number of login attempts in sessions variable
func loginAttempt(sess *sessions.Session) {
	// Log the attempt
	if sess.Values[sessLoginAttempt] == nil {
		sess.Values[sessLoginAttempt] = 1
	} else {
		sess.Values[sessLoginAttempt] = sess.Values[sessLoginAttempt].(int) + 1
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	// Get session
	fmt.Println("login")
	sess := session.Instance(r)
	// Display the view
	v := NewView(r)
	v.Name = "account/login"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	// Refill any form fields
	Repopulate([]string{"email"}, r.Form, v.Vars)
	v.Render(w)
}

func Home(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)
	// Display the view
	v := NewView(r)
	v.Name = "account/home"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	// Refill any form fields
	v.Vars["firstname"] = sess.Values["firstname"]
	v.Render(w)
}

func AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)
	// Prevent brute force login attempts by not hitting MySQL and pretending like it was invalid :-)
	if sess.Values[sessLoginAttempt] != nil && sess.Values[sessLoginAttempt].(int) >= 5 {
		log.Println("Brute force login prevented")
		sess.AddFlash(Flash{"Sorry, no brute force :-)", FlashNotice})
		sess.Save(r, w)
		Login(w, r)
		return
	}
	// Validate with required fields
	if validate, missingField := Validate(r, []string{"email", "password"}); !validate {
		sess.AddFlash(Flash{"Field missing: " + missingField, FlashError})
		sess.Save(r, w)
		Login(w, r)
		return
	}
	// Form values
	email := r.FormValue("email")
	password := r.FormValue("password")
	// Get database result
	userService = service.NewUserService("")
	result := model.User{}
	var err string
	result, err = userService.UserByEmail(email)
	fmt.Println(result)
  userService.Dispose()
	fmt.Println(result.Password, password)
	// Determine if user exists
	if err == message.DB_NOTFOUND {
		fmt.Println("norecord")
		loginAttempt(sess)
		sess.AddFlash(Flash{"No user found", FlashError})
		sess.Save(r, w)
	} else if err != ""{
		// Display error message
		fmt.Println(err)
		log.Println(err)
		sess.AddFlash(Flash{err, FlashError})
		sess.Save(r, w)
	} else if result.Deleted == 1{
		fmt.Println("deleted")
			// User inactive and display inactive message
		sess.AddFlash(Flash{"Account is inactive so login is disabled.", FlashNotice})
		sess.Save(r, w)
	} else if !passhash.MatchString(result.Password, password){
		fmt.Println("wrong pwd")
		loginAttempt(sess)
		sess.AddFlash(Flash{"Password is incorrect - Attempt: " + fmt.Sprintf("%v", sess.Values[sessLoginAttempt]), FlashWarning})
		sess.Save(r, w)
	} else {
			// Login successfully
			session.Empty(sess)
			fmt.Println("good")
			sess.AddFlash(Flash{"Login successful!", FlashSuccess})
			sess.Values["id"] = email
			sess.Values["firstname"] = result.FirstName
			sess.Save(r, w)
			fmt.Println(sess.Values["id"])
			http.Redirect(w, r, "/meeting/index", http.StatusFound)
			return
	}
	Login(w,r)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)
	// If user is authenticated
	if sess.Values["id"] != nil {
		session.Empty(sess)
		sess.AddFlash(Flash{"Goodbye!", FlashNotice})
		sess.Save(r, w)
	}
	http.Redirect(w, r, "/account/login", http.StatusFound)
}

func Register(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)
	// Display the view
	v := NewView(r)
	v.Name = "account/register"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	// Refill any form fields
  v.Render(w)
}

func SaveUser(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)
	// Validate with required fields
	validate, missingField := Validate(r, []string{"email", "firstname","lastname", "password", "confirmpassword"})
	if !validate {
		sess.AddFlash(Flash{"Field missing: " + missingField, FlashError})
		sess.Save(r, w)
		Register(w, r)
		return
	}
	// Form values
	hash,_ :=  passhash.HashString(r.FormValue("password"))
	user := model.User{
		Email:r.FormValue("email"),
		FirstName:r.FormValue("firstname"),
		LastName:r.FormValue("lastname"),
	  Password:hash}
	// Get database result
	userService = service.NewUserService("")
	var err string
	err = userService.UserCreate(user)
  userService.Dispose()
	// Determine if user exists
  if err == message.DB_TRANSACTION{
		// Display error message
		log.Println(err)
		sess.AddFlash(Flash{err, FlashError})
		sess.Save(r, w)
	} else {
			// Login successfully
	  session.Empty(sess)
	  sess.AddFlash(Flash{"Login successful!", FlashSuccess})
		sess.Values["email"] = user.Email
		sess.Values["firstname"] = user.FirstName
		sess.Save(r, w)
		http.Redirect(w, r, "/account/login", http.StatusFound)
		return
	}
}
