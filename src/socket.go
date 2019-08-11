package main

// 开端口

import (
	"bufio"
	"net"
	"strconv"
)

type Socket struct {
	WS *Client
	listen *net.Listener
	port uint16
	sends *map[string] *chan []byte
}
func writePump(socket *Socket, conn net.Conn){
	send := make(chan []byte)
	addr := conn.RemoteAddr().String()
	(*socket.sends)[addr] = &send

	for {
		select {
			case data, ok:= <- send:
				if !ok {			// chan 被关闭 (客户端断开)
					msg := ParseCommand("close", addr, socket.port, nil, "", false)
					socket.WS.send <- msg

					(*socket.sends)[addr] = nil
					conn.Close()
					return
				}
				_, err := conn.Write(data)
				if err != nil {		// 存在错误 关闭客户端
					close(send)
					return
				}
				msg := ParseCommand("send", addr, socket.port,  &data, "", false)
				socket.WS.send <- msg
		}
	}
}

func handleConnection(socket *Socket, conn net.Conn) {
	addr := conn.RemoteAddr().String()
	go writePump(socket, conn)
	reader := bufio.NewReader(conn)
	for {
		data, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}
		msg := ParseCommand("recv", addr, socket.port, &data, "", false)
		socket.WS.send <- msg
	}
	send := (*socket.sends)[addr]
	if send != nil{
		close(*send)
	}
}

func listen(client *Client, port uint16) {
	ln, err := net.Listen("tcp", ":" + strconv.Itoa(int(port)))
	if err != nil {
		msg := ParseCommand("create", "", port, nil, "Listen Failure!", true)
		client.send <- msg
		return
	} else {
		msg := ParseCommand("create", "", port, nil, "Listen Success!", false)
		client.send <- msg
	}
	sends := map[string] *chan []byte {}
	socket := &Socket{
		WS:      client,
		listen:  &ln,
		sends:	 &sends,
	}
	(*client.sockets)[port] = socket

	for {
		conn, err := ln.Accept()
		if err != nil {
			break
		}
		msg := ParseCommand("connect", conn.RemoteAddr().String(), socket.port, nil, "", false)
		client.send <- msg
		go handleConnection(socket, conn)
	}
}
