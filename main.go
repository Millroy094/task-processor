package main

import ("encoding/json", "net/http", "log")

type Task struct {
	ID int `json:"id"`
	Type int `json:"type"`
	Payload int `json:"payload"`
}