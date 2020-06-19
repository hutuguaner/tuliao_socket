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

type sendPositionData struct {
	Lng   float64 `json:"lng"`
	Lat   float64 `json:"lat"`
	Email string `json:"email"`
}

type sendBroadcastData struct {
	Broadcast string `json:"broadcast"`
	Email     string `json:"email"`
}

type sendChatMsgData struct {
	Msg  string `json:"msg"`
	From string `json:"from"`
	To   string `json:"to"`
}
