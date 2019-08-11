# WebSocket-Port-Forwarding
1. 通过WebSocket动态监听端口
2. 把端口接收到的数据通过WebSocket返回给客户端, 并支持通信
3. 当WebSocket断开后, 自动清空创建的端口

## 客户端消息

### 打开端口
```json
{"command": "create", "data": {"port": 10000}}
```

### 发送数据
```json
{"Command":"send","Data":{"addr":"127.0.0.1:12345","port":10000,"data":"Q2xvdWRZdQo="}}
```

### 关闭某个客户端
```json
{"command": "close", "data": {"addr":"127.0.0.1:12345", "port": 10000}}
```

### 关闭端口
```json
{"command": "close", "data": {"port": 10000}}
```

## 服务端消息

### 创建端口
```json
{"command":"create","data":{"port":10000,"message":"Listen Success!","err":false}}
{"command":"create","data":{"port":10000,"message":"Listen Failure!","err":true}}
```
### 用户连接
```json
{"command":"connect","data":{"addr":"127.0.0.1:123456","err":false}}
```
### 发送数据 (客户端发的, 再原样返回)
```json
{"Command":"send","Data":{"addr":"127.0.0.1:12345","port":10000,"data":"Q2xvdWRZdQo=","err":false}}
```
```json
{"command":"send","data":{"port":10001,"message":"Port not found!","err":true}}
{"command":"send","data":{"port":10000,"message":"Client not found!","err":true}}
```

### 接收数据
```json
{"Command":"recv","Data":{"addr":"127.0.0.1:12345","port":10000,"data":"Q2xvdWRZdQo=","err":false}}
```
### 用户断开
```json
{"command":"close","data":{"addr":"127.0.0.1:12345","err":false}}
```

### 关闭端口
```json
{"command":"close","data":{"port":10000,"err":false
{"command":"close","data":{"port":10000,"message":"Port not found!","err":true}}
```
