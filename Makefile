.PHONY: templ css build dev

# Download tailwindcss binary only if not already present
tailwindcss:
	@if [ ! -f ./tailwindcss ]; then \
		echo "Downloading tailwindcss..."; \
		curl -sL https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 -o tailwindcss; \
		chmod +x tailwindcss; \
	fi

templ:
	go tool templ generate

css: tailwindcss
	./tailwindcss -i input.css -o static/style.css

# Production build: minified CSS, generated templ, compiled binary
build: tailwindcss
	go tool templ generate
	./tailwindcss -m -i input.css -o static/style.css
	go build -o bin/hostmonitor .

# Development: run tailwindcss watch + templ watch + air concurrently
dev: tailwindcss css templ
	./tailwindcss --watch -i input.css -o static/style.css & \
	go tool templ generate --watch & \
	go tool air; \
	wait
