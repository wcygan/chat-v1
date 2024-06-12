# Chat V1

This application is slightly smarter than [chat-v0](https://github.com/wcygan/chat-v0).

This is accomplished by leveraging an embedded [nats](https://github.com/nats-io/nats.go) server to help distribute messages to clients.

Another key difference is the API that clients use to connect & chat with; [chat-v0/proto/chat/v1/chat.proto](https://github.com/wcygan/chat-v0/blob/main/proto/chat/v1/chat.proto) maintains a simple schema where everything is a message that is sent to everyone else. 
Meanwhile, [chat-v1/proto/chat/v1/chat.proto](https://github.com/wcygan/chat-v1/blob/main/proto/chat/v1/chat.proto) takes things a step further by distinguishing "joining" a chat room from "sending" a chat message.
This will ultimately help [yap](https://github.com/wcygan/yap) move away from [ChatPacket](https://github.com/wcygan/yap/blob/caeba7ca64312c661eb518d07cec0f038a62ebd3/proto/chat/v1/chat.proto#L26) which requires [a very complicated setup](https://github.com/wcygan/yap/blob/caeba7ca64312c661eb518d07cec0f038a62ebd3/yap-api/internal/chat/chat.go#L24) to handle sending & receives messages

# Quickstart

Generate proto files

```
buf generate proto
```

Run server

```
cd server
go run cmd/main.go
```

Run client

```
cd client
go run cmd/main.go
```

# Chat Application Examples

- [chat-v0](https://github.com/wcygan/chat-v0) - server that iterates over internal state to distribute messages to clients
- [chat-v1](https://github.com/wcygan/chat-v1) - server that uses nats (pub/sub) to distribute messages to clients