package service

import (
	"github.com/ya42/go_webapp/common/boltAdapter"
	"github.com/ya42/go_webapp/model"
	"github.com/ya42/go_webapp/model/message"
	"fmt"
)

type MeetingService struct{
	  Db boltAdapter.BoltConnection
}

func NewMeetingService(connStr string) *MeetingService{
	ser := new(MeetingService)
	ser.Db = boltAdapter.Connect("placeholder")
	return ser
}

func (ser *MeetingService) Dispose(){
	ser.Db.Database.Close()
	ser = nil
}

func (ms *MeetingService) MeetingsByUserID(userID string) ([]model.Meeting, string){
  result := []model.Meeting{}
	var errMsg string
  rawres,err := ms.Db.Seek("meeting", userID)
	fmt.Print(rawres)
	if err != nil {
		errMsg = message.DB_TRANSACTION
	}else if len(rawres) ==  0{
		errMsg = message.DB_NOTFOUND
	}
  for _,e := range rawres{
		fmt.Println(e)
		res, _ := ms.MeetingByID("user@test.commeeting 1")
		fmt.Println(res)
		result = append(result,res)
	}
	return result, errMsg
}

func (ms *MeetingService) MeetingByID(meetingID string) (model.Meeting, string){
  result := model.Meeting{}
	var err error
	var errMsg string
  err = ms.Db.View("meeting", meetingID, &result)
	if err != nil {
		errMsg = message.DB_TRANSACTION
	}else if result.Title == ""{
		errMsg = message.DB_NOTFOUND
	}
	return result, errMsg
}

func (ms *MeetingService) CreateMeeting(meeting model.Meeting) string {
	var errMsg string
	result := model.Meeting{}
	var err error
	err = ms.Db.Update("meeting", meeting.ID, &result)
	if err != nil{
		errMsg = message.DB_TRANSACTION
	}
	return errMsg
}

func (ms *MeetingService) UpdateMeeting(meeting model.Meeting) string {
	var errMsg string
  var err error
	err = ms.Db.Update("meeting", meeting.Title+meeting.CreatedBy, &meeting)
	if err != nil{
    errMsg = message.DB_TRANSACTION
	}
	return errMsg
}

func (ms *MeetingService) DeleteMeeting(meetingID string) string {
	var errMsg string
	err := ms.Db.Delete("meeting", meetingID)
  if err != nil{
			 errMsg = message.DB_TRANSACTION
	}
	return errMsg
}
