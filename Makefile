dev:
	go tool templ generate --watch --proxy="http://localhost:8080" --cmd="go run ."

build:
	go build -o bin/echotest main.go

watch-css:
	npx tailwindcss -i tailwindcss/styles.css -o static/styles.css --watch --optimize --minify