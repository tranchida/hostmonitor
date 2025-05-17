tailwind:
	curl -sL -C - https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 -o tailwindcss \
	&& chmod +x tailwindcss

templ:
	go tool templ generate

live: tailwind
	go tool air

build: tailwind templ
	./tailwindcss  -m -i input.css -o static/style.css
	go tool templ generate
	go build -o bin/hostmonitor main.go
