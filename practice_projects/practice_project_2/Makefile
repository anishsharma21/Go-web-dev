.PHONY: dbup dbdown up down gen

dbup up:
	goose -dir ./db/migrations sqlite3 app.db up

dbdown down:
	goose -dir ./db/migrations sqlite3 app.db down

gen:
	templ generate