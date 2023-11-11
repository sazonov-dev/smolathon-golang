package models

// Collection Masters
type Master struct {
	ID           string `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName    string `bson:"first_name" json:"first_name"`
	LastName     string `bson:"last_name" json:"last_name"`
	Email        string `bson:"email" json:"email"`
	PasswordHash string `bson:"password_hash" json:"password_hash"`
}

// Collection Events
type Event struct {
	Date string `bson:"date" json:"date"`
	Time string `bson:"time" json:"time"`
}

type Card struct {
	ID             string  `bson:"_id,omitempty" json:"id"`
	Title          string  `bson:"title" json:"title"`
	Description    string  `bson:"description" json:"description"`
	ImageName      string  `bson:"image_name" json:"image_name"`
	Location       string  `bson:"location" json:"location"`
	UpcomingEvents []Event `bson:"events,omitempty" json:"events"`
	PhoneNumber    string  `bson:"phone_number" json:"phone_number"`
}
