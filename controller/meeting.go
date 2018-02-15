package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ya42/go_webapp/model"
	"github.com/ya42/go_webapp/service"
	"github.com/ya42/go_webapp/common/passhash"
	"github.com/ya42/go_webapp/common/session"

	"github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
)

// MeetingpadReadGET displays the meetings in the meetingpad
func MeetingList(w http.ResponseWriter, r *http.Request) {
	// Get session
	fmt.Println("controller")
	sess := session.Instance(r)

	//userID := fmt.Sprintf("%s", sess.Values["id"])
  userID := ""
	fmt.Println("get meeting")
	meetings, err := model.MeetingsByUserID(userID)
	if err != nil {
		log.Println(err)
		meetings = []model.Meeting{}
	}

	// Display the view
	v := NewView(r)
	v.Name = "meeting/index"
	v.Vars["first_name"] = sess.Values["first_name"]
	v.Vars["meetings"] = meetings
	v.Render(w)
}

func CreateNewMeeting(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Display the view
	v := NewView(r)
	v.Name = "meeting/newmeeting"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Render(w)
}

func SaveMeeting(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Validate with required fields
	if validate, missingField := Validate(r, []string{"location","title","starttime"}); !validate {
		sess.AddFlash(Flash{"Field missing: " + missingField, FlashError})
		sess.Save(r, w)
		CreateNewMeeting(w, r)
		return
	}
	fmt.Println("content")
		// Get form values
	content := make([]string, 3)
	content[0] = r.FormValue("title")
	content[1] = r.FormValue("location")
	content[2] = r.FormValue("starttime")
	fmt.Println(content)
	var meetingID string
	if r.FormValue("meetingID") != ""{
		meetingID = r.FormValue("meetingID")
	}
	userID := fmt.Sprintf("%s", sess.Values["id"])
	fmt.Println(meetingID)

	// Get database result
	var err error
	if meetingID == ""{
		fmt.Println("new meeting")
	  err = service.CreateMeeting(content, userID)
	}else{
		fmt.Println("update meeting")
		err = service.UpdateMeeting(content, meetingID)
	}
	// Will only error if there is a problem with the query
	if err != nil {
		fmt.Println("error")
		log.Println(err)
		sess.AddFlash(Flash{"An error occurred on the server. Please try again later.", FlashError})
		sess.Save(r, w)
	} else {
		sess.AddFlash(Flash{"Meeting added!", FlashSuccess})
		sess.Save(r, w)
		fmt.Println("redirect")
		http.Redirect(w, r, "/meeting/index", http.StatusFound)
		return
	}
	// Display the same page
	fmt.Println("return list")
	MeetingList(w, r)
}

// UpdateMeeting displays the meeting update page
func UpdateMeeting(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Get the meeting id
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	meetingID := params.ByName("id")

	//userID := fmt.Sprintf("%s", sess.Values["id"])

	// Get the meeting
	meeting, err := model.MeetingByID(meetingID)
	if err != nil { // If the meeting doesn't exist
		log.Println(err)
		sess.AddFlash(Flash{"An error occurred on the server. Please try again later.", FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/meeting/index", http.StatusFound)
		return
	}

	// Display the view
	v := NewView(r)
	v.Name = "meeting/update"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Vars["title"] = meeting.Title
	v.Vars["location"] = meeting.Location
	v.Vars["starttime"] = meeting.Starttime
	v.Render(w)
}

// DeleteMeeting handles the meeting deletion
func DeleteMeeting(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	//userID := fmt.Sprintf("%s", sess.Values["id"])
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	meetingID := params.ByName("id")

	// Get database result
	err := service.DeleteMeeting(meetingID)
	// Will only error if there is a problem with the query
	if err != nil {
		log.Println(err)
		sess.AddFlash(Flash{"An error occurred on the server. Please try again later.", FlashError})
		sess.Save(r, w)
	} else {
		sess.AddFlash(Flash{"Meeting deleted!", FlashSuccess})
		sess.Save(r, w)
	}
	http.Redirect(w, r, "/meeting/index", http.StatusFound)
	return
}
