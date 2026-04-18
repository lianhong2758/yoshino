NAME=yoshino
EXE_NAME=${NAME}exe
VERSION=0.0.1

build:
	@echo "build!"
	@go version
	@go env -w GOPROXY=https://goproxy.cn,direct
	@go mod tidy
	@go build -ldflags="-checklinkname=0 -H windowsgui"
	@echo "Done!"
	
build_proto:
	@echo "build_proto!"
	@cd ${PROTOPATH} && protoc --go_out=. --go_opt=paths=source_relative *.proto