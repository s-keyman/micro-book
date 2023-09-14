package domain

import "time"

// User 领域对象，是 DDD 中的 entity
// BO(business object)
type User struct {
	Id       uint64
	Email    string
	Password string
	Nickname string
	AboutMe  string
	Birthday time.Time
	Ctime    time.Time
}
