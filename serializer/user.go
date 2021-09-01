package serializer

import (
	"blog/model"
	"strconv"
)

type User struct {
	UserName    string `json:"UserName"`
	FansNum     uint   `json:"FansNum"`
	FollowerNum uint   `json:"FollowerNum"`
	Briefly     string `json:"Briefly"`
}

func BuildUserList(users []model.User) []User {
	userlist := make([]User, len(users))
	for k, user := range users {
		userlist[k] = User{
			UserName:    user.UserName,
			FansNum:     user.FansNum,
			FollowerNum: user.FollowerNum,
			Briefly:     user.Briefly,
		}
	}
	return userlist
}
func BuildUserListResponse(userlist []model.User) Response {
	return Response{
		Data: BuildUserList(userlist),
		Msg:  strconv.Itoa(len(userlist)) + " user Display Succ!",
	}
}
