package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ya42/go_webapp/common/boltAdapter"
	"github.com/ya42/go_webapp/model"
)

var (
	meeting model.Meeting
	db boltAdapter.BoltConnection
)

func Initialize(){
	database := boltAdapter.Connect()
	db.Database = database
}

func MeetingsByID(meetingID string) (*[]model.Meeting, string) {
  result := db.database.View("meeting", meetingID)
	return result, standardizeError(err)
}
// MeetingByID gets meeting by ID
func MeetingsByName(meetingName string) (*[]model.Meeting, string) {
  db.Database.View("meeting", meetingID, )
	return result, standardizeError(err)
}

// MeetingsByUserID gets all meetings for a user
func MeetingsByUserID(userID string) (*[]model.Meeting, string) {
	var err error
	var result []Meeting
  err = db.Database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("meeting"))
		if b == nil {
			log.Println(bolt.ErrBucketNotFound)
			return nil, model.Error.SYSTEM_DBERROR
		}
		c := b.Cursor()
		prefix := []byte(userID)
		for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = c.Next() {
		  var single Meeting
		  err := json.Unmarshal(v, &single)
		  if err != nil {
			log.Println(err)
			continue
		  }
		  result = append(result, single)
	    }
	  })
	return &result, err
}

// MeetingCreate creates a meeting
func CreateMeeting(meeting model.Meeting) error {

}

// MeetingUpdate updates a meeting
func UpdateMeeting(meetingID string, userID string) error {
	var err error
		var meeting Meeting
		meeting, err = MeetingByID(meetingID)
		if err == nil {
			// Confirm the owner is attempting to modify the meeting
				meeting.Title = content[0]
				meeting.Location = content[1]
				meeting.Starttime = content[2]
				err = database.Update("meeting", meeting.ObjectID.Hex(), &meeting)
			} else {
				err = ErrUnauthorized
			}

	return standardizeError(err)
}

// MeetingDelete deletes a meeting
func DeleteMeeting(meetingID string, userID string) error {
	var err error
		var meeting Meeting
		meeting, err = MeetingByID(meetingID)
		if err == nil {
			// Confirm the owner is attempting to modify the meeting
				err = database.Delete("meeting", meeting.ObjectID.Hex())
			}
	return standardizeError(err)
}
