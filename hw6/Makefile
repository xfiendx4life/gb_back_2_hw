run:
	docker run -p 6379:6379 -d --rm redis
	PORT=8080 HOST=localhost REDISHOST=localhost REDISPORT=6379 TTL=1 go run ./cmd/confrmr/main.go
