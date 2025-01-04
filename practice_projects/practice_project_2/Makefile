.PHONY: dbup dbdown up down

dbup up:
	~/go/bin/goose -dir ./db/migrations sqlite3 app.db up

dbdown down:
	~/go/bin/goose -dir ./db/migrations sqlite3 app.db down