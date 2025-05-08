tailwindcss:
	npx @tailwindcss/cli -i style.css -o ./static/style.css --watch

run:
	go tool templ generate --watch --proxy="http://localhost:8080" --cmd="go run ."

live:
	make -j2 tailwindcss run

build:
	go build -o bin/hostmonitor main.go
