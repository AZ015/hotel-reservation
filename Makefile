build:
	go build -o bin/api

run: build
	./bin/api --listenAddr :8080

seed:
	go run scripts/seed.go

docker:
	echo "building docker file"
	docker build -t api .
	echo "running API inside Docker container"
	docker run -p 3000:3000 api

test:
	go test -v ./... -count=1