package models

type User struct {
	UserID string `json:"userID" dynamodbav:"userID"`
	Email  string `json:"email" dynamodbav:"email"`
}
