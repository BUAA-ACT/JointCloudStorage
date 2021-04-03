package model

type User struct {
	Id string
}

func Authenticate(sid string) (*User, error) {
	return &User{Id: "tester"}, nil
}
