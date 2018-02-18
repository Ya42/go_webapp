package model

import (
	"time"
)

type User struct {
	FirstName string        `db:"first_name" bson:"first_name"`
	LastName  string        `db:"last_name" bson:"last_name"`
	Email     string        `db:"email" bson:"email"`
	Password  string        `db:"password" bson:"password"`
	CreatedOn time.Time     `db:"created_on" bson:"created_on"`
	UpdatedOn time.Time     `db:"updated_on" bson:"updated_on"`
	Deleted   uint8         `db:"deleted" bson:"deleted"`
}
