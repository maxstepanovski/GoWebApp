package models

import "strconv"

type Users struct {
	Users []User
}

type User struct {
	ID    int    "json:id"
	Name  string "json:name"
	Email string "json:email"
	First string "json:email"
	Last  string "json:email"
}

func (u *User) ToString() string {
	return strconv.Itoa(u.ID) + " " + u.Name + " " + u.Email + " " + u.First + " " + u.Last
}
