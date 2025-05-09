.PHONY: build docker-build run docker-run clean

build:
	go build -o build/backend main.go

docker-build:
	docker build -t my-backend .

run: build
	./build/backend

docker-run:
	docker-compose up --build

clean:
	rm -f build/backend
	docker-compose down

