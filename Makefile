build:
	go build -o bin/api

run: build
	./bin/api --listenAddr :8080

seed:
	go run scripts/seed.go

test:
	go test -v ./... -count=1