dev:
	npm run dev

build:
	npm run build

preview:
	npm run preview

install:
	npm install

docker-build:
	docker build -t hostmonitor .

.PHONY: dev build preview install docker-build
