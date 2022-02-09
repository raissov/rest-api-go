package user

type User struct {
	ID           string `json:"ID" bson:"_id,omitempty"`
	Username     string
	PasswordHash string
	Email        string
}
