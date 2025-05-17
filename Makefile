tailwind:
	curl -sL -C - https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 -o tailwindcss \
	&& chmod +x tailwindcss

live: tailwind
	go tool air

build:
	go build -o bin/hostmonitor main.go
