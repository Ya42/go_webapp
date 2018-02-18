package service

import (
	"time"
	"github.com/ya42/go_webapp/model"
	"github.com/ya42/go_webapp/common/boltAdapter"
	"github.com/ya42/go_webapp/model/message"
)

type UserService struct{
	  Db boltAdapter.BoltConnection
}

func NewUserService(connStr string) *UserService{
	ser := new(UserService)
	ser.Db = boltAdapter.Connect(connStr)
	return ser
}

func (ser *UserService) Dispose(){
	ser.Db.Database.Close()
	ser = nil
}

func (us *UserService) UserByEmail(email string) (model.User, string){
	result := model.User{}
	var err error
	var errMsg string
	err = us.Db.View("user", email, &result)
	if err != nil {
    errMsg = message.DB_TRANSACTION
	}else if result.Email == ""{
		errMsg = message.DB_NOTFOUND
	}
  return result, errMsg
}

func (us *UserService) UserCreate(userModel model.User) string{
  userModel.CreatedOn = time.Now()
	userModel.Deleted = 0
	var errMsg string
	var err error
	err = us.Db.Update("user", userModel.Email, &userModel)
	if err != nil {
		errMsg = message.DB_TRANSACTION
	}
	return errMsg
}

func (us *UserService) UserUpdate(userModel model.User) string {
	userModel.UpdatedOn = time.Now()
	var err error
	var errMsg string
	err = us.Db.Update("user", userModel.Email, &userModel)
	if err != nil {
		errMsg = message.DB_TRANSACTION
	}
	return errMsg
}
