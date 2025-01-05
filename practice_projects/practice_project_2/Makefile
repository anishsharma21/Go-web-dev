.PHONY: dbup dbdown up down templ g tailwind t

dbup up:
	goose -dir ./db/migrations sqlite3 app.db up

dbdown down:
	goose -dir ./db/migrations sqlite3 app.db down

templ g:
	templ generate

tailwind t:
	./tailwindcss -i public/css/input.css -o public/css/output.css