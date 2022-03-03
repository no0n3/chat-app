package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

type TcpData struct {
	UserId  string
	Message []byte
}

var TCP_PORT = ":1235"

func tcpSendMessage(message TcpData, ip string) error {
	json, err := json.Marshal(message)
	if err != nil {
		return err
	}
	tcpCon, err := net.Dial("tcp", ip+TCP_PORT)
	if err != nil {
		return err
	}
	defer tcpCon.Close()

	fmt.Fprintf(tcpCon, string(json)+"\n")

	return nil
}

func tcpListen() error {
	l, err := net.Listen("tcp", TCP_PORT)
	if err != nil {
		return err
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleTcpCon(c)
	}
}

func handleTcpCon(c net.Conn) {
	defer c.Close()

	netData, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		return
	}

	receivedMsg := &TcpData{}
	json.Unmarshal([]byte(netData), receivedMsg)

	var msg map[string]interface{}
	json.Unmarshal(receivedMsg.Message, &msg)

	WS_HUB.sendMessageChan <- WsMessageData{
		userId:  receivedMsg.UserId,
		message: receivedMsg.Message,
	}
}
