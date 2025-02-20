dev:
	go tool air

build:
	go build -o bin/echotest main.go

watch-css:
	npx tailwindcss -i tailwindcss/styles.css -o static/styles.css --watch