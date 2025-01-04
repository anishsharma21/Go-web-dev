.PHONY: dbup dbdown

dbup:
	~/go/bin/goose -dir ./db/migrations sqlite3 app.db up

dbdown:
	~/go/bin/goose -dir ./db/migrations sqlite3 app.db down