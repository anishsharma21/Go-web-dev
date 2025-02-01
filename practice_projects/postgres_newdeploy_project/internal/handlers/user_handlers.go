package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func AddUser(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		insertUserQuery := "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id;"
		randomName := randString(rand.Intn(5) + 5)
		randomEmail := fmt.Sprintf("%s@gmail.com", randomName)

		var id int
		err := db.QueryRow(insertUserQuery, randomName, randomEmail).Scan(&id)
		if err != nil {
			log.Printf("Error inserting user into table: %v\n", err)
			http.Error(w, "error adding user", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "User added to database:\nID: %d\nName: %s\nEmail: %s\n", id, randomName, randomEmail)
	})
}
