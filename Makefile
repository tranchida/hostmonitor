tailwindcss:
	npx @tailwindcss/cli -i style.css -o ./static/style.css --watch

run:
	air

live:
	make -j2 tailwindcss run

build:
	go build -o bin/hostmonitor main.go
