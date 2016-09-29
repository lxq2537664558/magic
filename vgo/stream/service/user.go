package service

// User group user
type User struct {
	Name        string
	Sms         string
	Mail        string
	MessagePush string
}

func NewUser() *User {
	return &User{}
}

// Type       uint8
// Phone      string
// Operatorid string
