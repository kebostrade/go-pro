module server

go 1.23

require (
	github.com/dima/go-pro/advanced-topics/08-grpc-distributed v0.0.0
	github.com/fatih/color v1.17.0
	google.golang.org/grpc v1.67.1
	google.golang.org/protobuf v1.35.1
)

replace github.com/dima/go-pro/advanced-topics/08-grpc-distributed => ../

require (
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/sys v0.0.0-20211116161374-3aa7ad689f93 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241015192408-796eee8c2d53 // indirect
)
