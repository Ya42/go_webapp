package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ya42/go_webapp/common/databaseAdapter/boltAdapter"
	"github.com/ya42/go_webapp/model"
)

// MeetingByID gets meeting by ID
func MeetingByName(meetingID string) (Meeting, string) {

}

	return result, standardizeError(err)
}

// MeetingsByUserID gets all meetings for a user
func MeetingsByUserID(userID string) (*[]model.Meeting, string) {
	var err error
	var result []Meeting
  err = database.BoltDB.View(func(tx *bolt.Tx) error {
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
			return nil
		})
	}
	return &result, err
}

// MeetingCreate creates a meeting
func CreateMeeting(meeting model.Meeting) error {
	var err error
	now := time.Now()

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
		_, err = database.SQL.Exec("INSERT INTO meeting (content, user_id) VALUES (?,?)", content, userID)
	case database.TypeMongoDB:
		if database.CheckConnection() {
			// Create a copy of mongo
			session := database.Mongo.Copy()
			defer session.Close()
			c := session.DB(database.ReadConfig().MongoDB.Database).C("meeting")

			meeting := &Meeting{
				ObjectID:  bson.NewObjectId(),
				Title:   content[0],
				Location: content[1],
				Starttime: content[2],
				CreatedOn: now,
				CreatedBy: userID,
				Deleted:   0,
			}
			err = c.Insert(meeting)
		} else {
			err = ErrUnavailable
		}
	case database.TypeBolt:
      meeting := &Meeting{
					ObjectID:  bson.NewObjectId(),
					Title:   content[0],
					Location: content[1],
					Starttime: content[2],
					CreatedOn: now,
					CreatedBy: userID,
					Deleted:   0,
		}

		err = database.Update("meeting", meeting.ObjectID.Hex(), &meeting)
	default:
		err = ErrCode
	}

	return standardizeError(err)
}

// MeetingUpdate updates a meeting
func UpdateMeeting(meetingID string, userID string) error {
	var err error

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
		_, err = database.SQL.Exec("UPDATE meeting SET content=? WHERE id = ? AND user_id = ? LIMIT 1", content[0], meetingID,meetingID )
	case database.TypeMongoDB:
		if database.CheckConnection() {
			// Create a copy of mongo
			session := database.Mongo.Copy()
			defer session.Close()
			c := session.DB(database.ReadConfig().MongoDB.Database).C("meeting")
			var meeting Meeting
			meeting, err = MeetingByID(meetingID)
			if err == nil {
				// Confirm the owner is attempting to modify the meeting
					meeting.Title = content[0]
					meeting.Location = content[1]
					meeting.Starttime = content[2]
					err = c.UpdateId(bson.ObjectIdHex(meetingID), &meeting)
			}
		} else {
			err = ErrUnavailable
		}
	case database.TypeBolt:
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
	default:
		err = ErrCode
	}

	return standardizeError(err)
}

// MeetingDelete deletes a meeting
func DeleteMeeting(meetingID string, userID string) error {
	var err error

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
		//_, err = database.SQL.Exec("DELETE FROM meeting WHERE id = ? AND user_id = ?", meetingID, userID)
	case database.TypeMongoDB:
		if database.CheckConnection() {
			// Create a copy of mongo
			session := database.Mongo.Copy()
			defer session.Close()
			c := session.DB(database.ReadConfig().MongoDB.Database).C("meeting")

			var meeting Meeting
			meeting, err = MeetingByID(meetingID)
			if err == nil {
				// Confirm the owner is attempting to modify the meeting
					err = c.RemoveId(bson.ObjectIdHex(meetingID))
					fmt.Println(meeting)
				} else {
					err = ErrUnauthorized
			}
		} else {
			err = ErrUnavailable
		}
	case database.TypeBolt:
		var meeting Meeting
		meeting, err = MeetingByID(meetingID)
		if err == nil {
			// Confirm the owner is attempting to modify the meeting
				err = database.Delete("meeting", meeting.ObjectID.Hex())
			}
	default:
		err = ErrCode
	}

	return standardizeError(err)
}
