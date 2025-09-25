package models

type Wallet struct {
	ID      string  `bson:"_id,omitempty" json:"id"`
	UserID  string  `bson:"user_id" json:"user_id"`
	Balance float64 `bson:"balance" json:"balance"`
}
