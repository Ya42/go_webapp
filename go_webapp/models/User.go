package model

import(
  "time"
)

type Account struct{
  Email string
  Username string
  Password string
  Bio string
  Location string
  CreatedOn time.Time
  LastLogin time.Time
}
