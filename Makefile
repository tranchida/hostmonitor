templ:
	go tool templ generate --watch --proxy="http://localhost:8080" --cmd="go run ."

tailwindcss:
	npx tailwindcss -i tailwindcss/styles.css -o static/styles.css --watch --optimize --minify

live:
	make -j2 templ tailwindcss

build:
	go build -o bin/echotest main.go
