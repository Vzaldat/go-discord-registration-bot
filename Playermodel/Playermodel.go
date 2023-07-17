package Playermodel

type Player struct {
	Name       string `json:"Name"`
	InGameName string `json:"InGameName"`
	InGameTag  string `json:"GameTag"`
	Rank       string
}

var Registrations map[string]Player = make(map[string]Player)
