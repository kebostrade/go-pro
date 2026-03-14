module examples

go 1.23

require (
	github.com/dima/go-pro/advanced-topics/08-grpc-distributed v0.0.0
	google.golang.org/grpc v1.67.1
	google.golang.org/protobuf v1.35.1
)

replace github.com/dima/go-pro/advanced-topics/08-grpc-distributed => ../

require (
	golang.org/x/net v0.28.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241015192408-796eee8c2d53 // indirect
)
