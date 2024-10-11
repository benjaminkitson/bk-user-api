package models

type User struct {
	// TODO: ID
	Email string `json:"email" dynamodbav:"email"`
}
