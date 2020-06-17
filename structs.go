package main

import (
	_ "github.com/go-sql-driver/mysql"
)

type userOfflineData struct {
	Type  string `json:"type"`
	Email string `json:"email"`
}

type positionData struct {
	Type  string `json:"type"`
	Email string `json:"email"`
	Lng   int64  `json:"lng"`
	Lat   int64  `json:"lat"`
}
