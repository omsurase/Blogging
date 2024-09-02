package models

type User struct {
	ID       string `json:"id" bson:"_id"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenValidationRequest struct {
	Token string `json:"token"`
}
