package service

type User struct {
	Id             string
	Email          string
	HashedPassword string
}

type Note struct {
	UserId string
	Title  string
	Body   string
}
