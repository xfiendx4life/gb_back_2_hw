run-rest:
	STORAGE="./storage" RESTPORT=:8080 go run ./cmd/main.go -mode=n -server=rest
run-rpc:
	STORAGE="./storage" RPCPORT=:8000 go run ./cmd/main.go -mode=n -server=rpc
run-both:
	STORAGE="./storage" RPCPORT=:8000 RESTPORT=:8080 go run ./cmd/main.go -mode=n
test-all:
	go test ./...