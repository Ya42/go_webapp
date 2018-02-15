package service

import (
	//"github.com/ya42/go_webapp/common/boltAdapter"
	"github.com/ya42/go_webapp/model"
)

var (
	user model.User
	//db boltAdapter.BoltConnection
)

// UserByEmail gets user information from email
func UserByEmail(email string) (User, error) {
	fmt.Println(user)
	var err error
	result := User{}
		err = database.View("user", email, &result)
		if err != nil {
			err = ErrNoResult
		}
    return result, standardizeError(err)
}

// UserCreate creates user
func UserCreate(firstName, lastName, email, password string) error {
	var err error
	now := time.Now()
		user := &User{
			ObjectID:  bson.NewObjectId(),
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
			Password:  password,
			StatusID:  1,
			CreatedAt: now,
			UpdatedAt: now,
			Deleted:   0,
		}
		err = database.Update("user", user.Email, &user)
	return standardizeError(err)
}
