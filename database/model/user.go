package model

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	MongoID  primitive.ObjectID `json:"_id"      bson:"_id,omitempty"      redis:"-"`
	ID       string             `json:"-"        bson:"id,omitempty"       redis:"id"`
	Name     string             `json:"name"     bson:"name,omitempty"     redis:"name"`
	Username string             `json:"username" bson:"username,omitempty" redis:"username"`
	Email    string             `json:"email"    bson:"email,omitempty"    redis:"-"`
	Password string             `json:"-"        bson:"password,omitempty" redis:"-"`
}

func NewUser(name, username, email, password string) *User {
	return &User{
		ID:       uuid.NewString(),
		MongoID:  primitive.NewObjectID(),
		Username: username,
		Name:     name,
		Email:    email,
		Password: password,
	}
}

// HSet("myhash", []string{"key1", "value1", "key2", "value2"})
func (u *User) ToHSET() []string {
	uH := []string{
		"id", u.ID,
		"mID", u.MongoID.Hex(),
		"name", u.Name,
		"username", u.Username,
	}

	return uH
}
