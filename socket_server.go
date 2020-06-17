package main

import (
	"container/list"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/tidwall/gjson"
)

func socketServerStart() {
	//
	go broadCastMsg()
	//
	service := ":1583"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		os.Exit(1)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

var lenRead = 128

func handleClient(conn net.Conn) {

	defer conn.Close()

	//conn.SetReadDeadline(time.Now().Add(2 * time.Minute))
	request := make([]byte, lenRead)
	readLen, err := conn.Read(request)
	if err != nil || readLen == 0 {
		conn.Close()
		conn = nil

	} else {
		email := strings.TrimSpace(string(request[:readLen]))
		addConnToMap(email, conn)

		for {
			readLen, err := conn.Read(request)
			if err != nil {
				conn.Close()
				conn = nil

				connMapLock.Lock()
				delete(connMap, email)

				connMapLock.Unlock()

				break
			}
			if readLen == 0 {
				fmt.Println("长度为0")
				break
			} else {
				msg := strings.TrimSpace(string(request[:readLen]))
				fmt.Println("收到 ： " + msg)
				addMsgToQueue(msg)
			}
			request = make([]byte, lenRead)
		}

	}

}

var msgQueueLock sync.Mutex
var msgQueue = list.New()

func addMsgToQueue(msg string) {
	msgQueueLock.Lock()
	msgQueue.PushBack(msg)
	msgQueueLock.Unlock()
}

var connMapLock sync.Mutex
var connMap map[string]net.Conn

func addConnToMap(email string, conn net.Conn) {
	connMapLock.Lock()
	if connMap == nil {
		connMap = make(map[string]net.Conn)
	}
	connMap[email] = conn
	connMapLock.Unlock()
}

//开启线程 将队列里的广播消息 一次发送出去
func broadCastMsg() {
	for {
		if msgQueue != nil && msgQueue.Len() > 0 {
			msgQueueLock.Lock()
			msg := msgQueue.Front()
			msgQueue.Remove(msg)
			msgQueueLock.Unlock()

			connMapLock.Lock()
			for email, conn := range connMap {
				if conn == nil {

					connMapLock.Lock()
					delete(connMap, email)
					connMapLock.Unlock()

				} else {
					fmt.Println("广播 ： " + msg.Value.(string))

					broadCastMsgDo(msg.Value.(string), conn, email)
				}
			}

			connMapLock.Unlock()
		}
		time.Sleep(1 * time.Second)
	}

}

//
func broadCastMsgDo(msg string, conn net.Conn, email string) {

	dataType := gjson.Get(msg, "type")

	if dataType.String() == "0" {
		//位置数据
		conn.Write([]byte(msg))
	} else if dataType.String() == "1" {
		//广播数据
		conn.Write([]byte(msg))
	} else if dataType.String() == "2" {
		//聊天数据
		to := gjson.Get(msg, "to")
		if to.String() == email {
			conn.Write([]byte(msg))
		}
	} else if dataType.String() == "3" {
		//用户下线
		conn.Write([]byte(msg))
		//
		email := gjson.Get(msg, "email")
		deleteUserFromDB(email.String())
	}

}

//从数据库中 把下线的用户删掉
func deleteUserFromDB(email string) error {
	if !hasDbInit {
		initDb()
	}
	insForm, err := myDb.Prepare("delete from user where email=?")
	if err != nil {
		return err
	}
	insForm.Exec(email)
	return nil
}
