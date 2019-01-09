package secondary

import "strconv"

type Users struct {
	Users []User
}

type User struct {
	ID    int    "json:id"
	Name  string "json:name"
	Email string "json:email"
	First string "json:first"
	Last  string "json:last"
	Image string "json:image"
}

func (u *User) ToString() string {
	return strconv.Itoa(u.ID) + " " + u.Name + " " + u.Email + " " + u.First + " " + u.Last + " " + u.Image
}
