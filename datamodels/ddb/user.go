package ddb

type User struct {
	ID             string `json:"UserID"`
	Name           string `json:"Name"`
	Email          string `json:"Email"`
	PhoneNo        string `json:"PhoneNo"`
	ProfilePicLink string `json:"ProfilePicLink"`
	IsActive       bool   `json:"IsActive"`
}
