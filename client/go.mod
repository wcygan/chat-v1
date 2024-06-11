module github.com/wcygan/chat-v1/client

go 1.22.0

require (
	github.com/wcygan/chat-v1/generated/go v0.0.0
	google.golang.org/grpc v1.64.0
	google.golang.org/protobuf v1.34.2
)

require (
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240610135401-a8a62080eff3 // indirect
)

replace github.com/wcygan/chat-v1/generated/go => ../generated/go
