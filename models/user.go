package models

type User struct {
	ID       string `bson:"_id,omitempty" json:"id"`
	Name     string `bson:"name" json:"name"`
	Email    string `bson:"email" json:"email"`
	Phone    string `bson:"phone" json:"phone"`
	Password string `bson:"password" json:"-"`
	Role     string `bson:"role" json:"role"`
}
