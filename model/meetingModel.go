package model

import (
	"time"
	//"gopkg.in/mgo.v2/bson"
)

// *****************************************************************************
// Meeting
// *****************************************************************************

// Meeting table contains the information for each meeting
type Meeting struct {
	ID  string `db:"id" bson:"id"`
	Title   string  `db:"title" bson:"title"`
	Location string `db:"location" bson:"location"`
	Starttime string `db:"starttime" bson:"starttime"`
	CreatedBy  string `bson:"createdby"`
	CreatedOn time.Time `db:"createdon" bson:"createdon"`
	Deleted   uint8 `db:"deleted" bson:"deleted"`
}
