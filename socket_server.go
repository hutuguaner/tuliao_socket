package main

import (
	"container/list"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
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
	readLen,err:=conn.Read(request)
	if err!=nil||readLen==0 {
		conn.Close()
		conn = nil
		
	}else{
		email:=strings.TrimSpace(string(request[:readLen]))
		addConnToMap(email,conn)

		for {
			readLen, err := conn.Read(request)
			if err != nil{
				conn.Close()
				conn = nil
				connMapLock.Lock()
				delete(connMap,email)
				connMapLock.Unlock()
				break
			}
			if readLen == 0 {
				fmt.Println("长度为0")
				break
			} else {
				msg := strings.TrimSpace(string(request[:readLen]))
				fmt.Println("收到 ： "+msg)
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

func addConnToMap(email string,conn net.Conn) {
	connMapLock.Lock()
	if connMap==nil {
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
			for email,conn:=range connMap{
				if conn==nil {
					delete(connMap,email)
				}else{
					fmt.Println("广播 ： "+msg.Value.(string))
					conn.Write([]byte(msg.Value.(string)))
				}
			}
			
			connMapLock.Unlock()
		}
		time.Sleep(1 * time.Second)
	}

}
