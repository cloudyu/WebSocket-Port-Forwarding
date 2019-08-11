package main

type Data struct {
	Addr   	string			`json:"addr,omitempty"`				// 	ip:port 	格式
	Port   	uint16			`json:"port,omitempty"`				//	port 		操作的端口
	Data 	[]byte			`json:"data,omitempty"`				//	数据
	Message string			`json:"message,omitempty"`			//	消息什么的
	Err		bool			`json:"err"`						//	是否存在错误
}

type Command struct {
	Command string		`json:"command"`
	Data Data		`json:"data"`
}

