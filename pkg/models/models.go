package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Username string             `bson:"username" json:"username"`
	Password string             `bson:"password" json:"password"`
	Roles    []string           `bson:"roles" json:"roles"`
}

type UserClaims struct {
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}
