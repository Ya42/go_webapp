package model

import (
	"time"
	//"gopkg.in/go.v2/bson"
)

type Meeting struct {
	ID      string  `db:"id" bson:"id"`
	Title   string  `db:"title" bson:"title"`
	Location string `db:"location" bson:"location"`
	Starttime string `db:"starttime" bson:"starttime"`
	CreatedBy  string `db:"createdby" bson:"createdby"`
	CreatedOn time.Time `db:"createdon" bson:"createdon"`
	Deleted   uint8 `db:"deleted" bson:"deleted"`
}
