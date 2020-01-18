package music

import (
	"database/sql"
	"errors"
)

type Subscription struct {
	SubscriptionID int    `json:"subscription_id"`
	Name           string `json:"name"`
	Price          int    `json:"price"`
	Duration       int    `json:"duration"`
}


func GetSubscriptions(db *sql.DB) ([]Subscription, error) {
	return nil, errors.New("not implemented")
}

func (s *Subscription) GetSubscriptionByID(db *sql.DB) error {
	return errors.New("not implemented")
}

func (s *Subscription) CreateSubscription(db *sql.DB) error {
	return errors.New("not implemented")
}
