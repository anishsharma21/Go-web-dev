package types

import "time"

// TODO set json tags so that values can be sent as JSON instead
type User struct {
	Id int
	Name string
	Email string
	CreatedAt time.Time
}