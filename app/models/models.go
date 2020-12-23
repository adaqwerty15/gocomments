package models

import (
	"time"
)

type Comment struct {
	Id       int `json:"id"`
	PageId   int `json:"page_id"`
	PageUrl   string `json:"page_url"`
	UserId   int `json:"user_id"`
	UserFirstName   string `json:"user_first_name"`
	UserLastName   string `json:"user_last_name"`
	Timestamp  time.Time `json:"timestamp"`
	Text   string `json:"text"`
	Status   string `json:"status"`
	Important bool `json:"important"`
}

type User struct {
	Id   int `json:"id"`
	FirstName   string `json:"first_name"`
	LastName   string `json:"last_name"`
}

type Page struct {
	Id   int `json:"id"`
	Url   string `json:"url"`
}

type CommentInfo struct {
	Id       int `json:"id"`
	Page     Page `json:"page"`
	User   User `json:"user"`
	Timestamp  time.Time `json:"timestamp"`
	Text   string `json:"text"`
	Status   string `json:"status"`
	Important bool `json:"important"`
}