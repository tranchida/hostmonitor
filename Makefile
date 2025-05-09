run:
	go tool templ generate --watch --proxy="http://localhost:8080" --cmd="go run ."

build:
	go build -o bin/hostmonitor main.go
