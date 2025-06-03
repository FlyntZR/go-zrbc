package service

import (
	"go-zrbc/db"
	"go-zrbc/view"
)

func DBToViewUser(dUser *db.Member) *view.Member {
	vUser := view.Member{
		ID:       dUser.ID,
		User:     dUser.User,
		UserName: dUser.UserName,
	}
	return &vUser
}
