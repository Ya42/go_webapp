package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ya42/go_webapp/model"
	"github.com/ya42/go_webapp/model/message"
	"github.com/ya42/go_webapp/service"
	"github.com/ya42/go_webapp/common/session"

	"github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
)

var meetingService *service.MeetingService

// MeetingpadReadGET displays the meetings in the meetingpad
func MeetingList(w http.ResponseWriter, r *http.Request) {
	// Get session
	fmt.Println("controller")
	sess := session.Instance(r)

	//userID := fmt.Sprintf("%s", sess.Values["id"])
  /*userID := ""
	fmt.Println("get meeting")
	meetings, err := model.MeetingsByUserID(userID)
	if err != nil {
		log.Println(err)
		meetings = []model.Meeting{}
	}*/

	// Display the view
	v := NewView(r)
	v.Name = "meeting/index"
	v.Vars["first_name"] = sess.Values["first_name"]
	v.Vars["meetings"] = ""
	v.Render(w)
}

func CreateNewMeeting(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Display the view
	v := NewView(r)
	v.Name = "meeting/new"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Render(w)
}

func SaveMeeting(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)
  curuser := sess.Values["email"].(string)
	// Validate with required fields
	if validate, missingField := Validate(r, []string{"location","title","starttime"}); !validate {
		sess.AddFlash(Flash{"Field missing: " + missingField, FlashError})
		sess.Save(r, w)
		CreateNewMeeting(w, r)
		return
	}
  meeting := model.Meeting{
		Title:r.FormValue("title"),
		Location:r.FormValue("location"),
		Starttime:r.FormValue("starttime"),
	  ID:r.FormValue("title")+curuser}
	var err string
	meetingService = service.NewMeetingService("")
	if meeting.ID == ""{
	  err = meetingService.CreateMeeting(meeting)
	}else{
		err = meetingService.UpdateMeeting(meeting)
	}
	// Will only error if there is a problem with the query
	if err == message.DB_TRANSACTION{
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

	// Display the view
	v := NewView(r)
	v.Name = "meeting/update"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Render(w)
}

// DeleteMeeting handles the meeting deletion
func DeleteMeeting(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	//userID := fmt.Sprintf("%s", sess.Values["id"])
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	meetingTitle := params.ByName("title")

	// Get database result
	err := meetingService.DeleteMeeting(meetingTitle)
	// Will only error if there is a problem with the query
	if err != message.DB_TRANSACTION {
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
