package main

import "log"

type Twot struct {
	Sequence int
	User     User
	Content  string
}

type User struct {
	Name string
}

type Twotter struct {
	Twots []Twot
	Users []User
}

func LoadTwots() *Twotter {
	u1 := User{"Alice"}
	u2 := User{"Bob"}

	twotter := Twotter{
		Twots: []Twot{
			{1, u1, "Hello"},
			{2, u2, "Hi"},
			{3, u1, "You're a bit of a twot"},
			{4, u2, "What!?"},
			{5, u1, "What?"},
			{6, u2, "..."},
		},
		Users: []User{u1, u2},
	}

	return &twotter
}
