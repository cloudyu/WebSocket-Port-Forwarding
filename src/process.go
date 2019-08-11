package main

import (
	"encoding/json"
)
func ParseCommand(command , addr string, port uint16, data *[]byte, message string, err bool) []byte {
	if data == nil {
		data = &[]byte{}
	}
	com := Command{
		Command: command,
		Data:    Data{
			Addr:	 addr,
			Port:    port,
			Data:    *data,
			Message: message,
			Err:     err,
		},
	}
	msg, _ := json.Marshal(com)
	return msg
}

func process(client *Client , message []byte){
	command := &Command{}
	err := json.Unmarshal(message, command)
	if err != nil {	// 错误格式不处理
		return
	}
	switch command.Command {
	case "create":
		port := command.Data.Port
		if port >= 10000 {
			go listen(client, port)
		} else {
			msg := ParseCommand("create", "", 0, nil, "Port out of range!", true)
			client.send <- msg
		}
	case "send":
		addr := command.Data.Addr
		port := command.Data.Port
		data := command.Data.Data
		if socket, ok := (*client.sockets)[port]; ok {
			if send, ok := (*socket.sends)[addr]; ok {
				if send != nil {
					*send <- data
				} else {
					msg := ParseCommand("send", "", port, nil, "Client not close!", true)
					client.send <- msg
				}
			} else {
				msg := ParseCommand("send", "", port, nil, "Client not found!", true)
				client.send <- msg
			}
		} else {
			msg := ParseCommand("send", "", port, nil, "Port not found!", true)
			client.send <- msg
		}
	case "close":
		addr := command.Data.Addr
		port := command.Data.Port
		if socket, ok := (*client.sockets)[port]; ok {
			if addr == "" { // 空addr 表示关闭整个端口
				for _, send := range *socket.sends{
					close(*send)
				}
				(*socket.listen).Close()
				delete(*client.sockets, port)
				msg := ParseCommand("close", addr, port, nil, "", false)
				client.send <- msg
			} else if send, ok := (*socket.sends)[addr]; ok {
				close(*send)	// 断线会有另一个包
			} else {
				msg := ParseCommand("close", addr, port, nil, "Client not found!", true)
				client.send <- msg
			}
		} else {
			msg := ParseCommand("close", addr, port, nil, "Port not found!", true)
			client.send <- msg
		}
	}
}