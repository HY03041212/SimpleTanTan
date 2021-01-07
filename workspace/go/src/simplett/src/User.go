package main

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Liked    []int  `json:"liked"`
	Matched  []int  `json:"matched"`
	Disliked []int  `json:"disliked"`
}

type Users []User

type BasicUser struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type BasicUsers []BasicUser

type Relation struct {
	Id    int    `json:"id"`
	State string `json:"state"`
	Type  string `json:"type"`
}

type Relations []Relation
