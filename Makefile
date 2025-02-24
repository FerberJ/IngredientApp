dev: build
	@make -j3 css templ air

run: build
	@./bin/app

build:
	@go mod tidy
	@go build -o bin/app .

css:
	@npm install
	@npx @tailwindcss/cli -i view/css/app.css -o public/styles.css --watch

templ:
	templ generate --watch --proxy=http://localhost:3000

air:
	@air