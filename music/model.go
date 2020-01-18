package music

import (
	"database/sql"
	"errors"
	"fmt"
)

type Subscription struct {
	SubscriptionID int    `json:"subscription_id"`
	Name           string `json:"name"`
	Price          int    `json:"price"`
	Duration       int    `json:"duration"`
}

func GetSubscriptions(db *sql.DB) ([]Subscription, error) {
	statement := fmt.Sprintf("SELECT * FROM subscriptions")
	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var subscriptions []Subscription

	for rows.Next() {
		var s Subscription
		if err := rows.Scan(&s.SubscriptionID, &s.Name, &s.Price, &s.Duration); err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, s)
	}
	if subscriptions == nil {
		return []Subscription{}, nil
	} else {
		return subscriptions, nil
	}
}

func (s *Subscription) GetSubscriptionByID(db *sql.DB) error {
	return errors.New("not implemented")
}

func (s *Subscription) CreateSubscription(db *sql.DB) error {
	return errors.New("not implemented")
}
