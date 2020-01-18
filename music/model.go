package music

import (
	"database/sql"
	"fmt"
)

type Subscription struct {
	SubscriptionID int    `json:"subscription_id"`
	Name           string `json:"name"`
	Price          int    `json:"price"`
	Duration       int    `json:"duration"`
}

func GetSubscriptions(db *sql.DB) ([]Subscription, error) {
	statement := fmt.Sprintf("SELECT * FROM subscription")
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
	statement := fmt.Sprintf("SELECT * FROM subscription WHERE id=%d", s.SubscriptionID)
	return db.QueryRow(statement).Scan(&s.SubscriptionID, &s.Name, &s.Price, &s.Duration)
}

func (s *Subscription) CreateSubscription(db *sql.DB) error {
	statement := fmt.Sprintf("INSERT INTO subscription(name, price, duration) VALUES('%s', %d, %d)", s.Name, s.Price, s.Duration)
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}
	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&s.SubscriptionID)
	if err != nil {
		return err
	}
	return nil
}
