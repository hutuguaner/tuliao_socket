package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

func socketIOServerStart() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected : ", s.ID())
		return nil
	})
	server.OnEvent("/", sendPosition, func(s socketio.Conn, msg string) {

		bytes := []byte(msg)
		data := sendPositionData{}
		err := json.Unmarshal(bytes, &data)

		if err == nil {
			fmt.Println("sendposition:", msg)
			s.SetContext(data.Email)
			s.Join(all)
			s.Join(data.Email)
			server.BroadcastToRoom("", all, receivePosition, msg)
		}

	})
	server.OnEvent("/", sendBroadcast, func(s socketio.Conn, msg string) {
		fmt.Println("sendBroadcast:", msg)
		//s.Emit("sendposition", "have "+msg)
		bytes := []byte(msg)
		data := sendBroadcastData{}
		err := json.Unmarshal(bytes, &data)
		if err == nil {
			server.BroadcastToRoom("", all, receiveBroadcast, msg)

		}

	})
	server.OnEvent("/", sendMsg, func(s socketio.Conn, msg string) {
		fmt.Println("sendposition:", msg)
		//s.Emit("sendposition", "have "+msg)
		bytes := []byte(msg)
		data := sendChatMsgData{}
		err := json.Unmarshal(bytes, &data)
		if err == nil {
			server.BroadcastToRoom("", data.To, receiveMsg, msg)
		}

	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", s.ID())
		email := s.Context().(string)
		s.Close()
		server.BroadcastToRoom("", all, offline, email)
	})
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", s.ID())
		email := s.Context().(string)
		s.Close()
		server.BroadcastToRoom("", all, offline, email)
	})
	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

//event
var sendPosition = "send_position"
var sendBroadcast = "send_broadcast"
var sendMsg = "send_chat_msg"
var offline = "offline"
var receivePosition = "receive_position"
var receiveBroadcast = "receive_broadcast"
var receiveMsg = "receive_chat_msg"

//room
var all = "all"
