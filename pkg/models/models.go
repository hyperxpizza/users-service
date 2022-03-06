package models

import "github.com/golang/protobuf/ptypes/timestamp"

type LoginData struct {
	Id           int64                `json:"id"`
	Username     string               `json:"username"`
	Email        string               `json:"email"`
	PasswordHash string               `json:"passwordHash"`
	PasswordSlat string               `json:"passwordSlat"`
	Created      *timestamp.Timestamp `json:"created"`
	Updated      *timestamp.Timestamp `json:"updated"`
	UserID       int64                `json:"userid"`
}

type UserData struct {
	Id int64 `json:"id"`
}

type User struct {
	login LoginData ``
}
