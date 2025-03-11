tailwindcss:
	tailwindcss -i tailwindcss/styles.css -o static/styles.css --watch --optimize --minify

live:
	make -j2 templ tailwindcss

build:
	go build -o bin/echotest main.go
