generate:
	protoc --proto_path=proto proto/*.proto --go_out=. --go-grpc_out=.

build:
	go build -o bin/server cmd/server/main.go

test:
	go test -v ./...

run:
	go run server/main.go

docker-build:
	docker build -t chromadb-sizing-estimator .

docker-run:
	docker run -p 8080:8080 -it chromadb-sizing-estimator

lint:
	golangci-lint run