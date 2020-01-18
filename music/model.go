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

type Customer struct {
	CustomerID int    `json:"customer_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Balance    int    `json:"balance"`
}

func GetCustomers(db *sql.DB) ([]Customer, error) {
	statement := fmt.Sprintf("SELECT * FROM customer")
	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var customers []Customer

	for rows.Next() {
		var c Customer
		if err := rows.Scan(&c.CustomerID, &c.Name, &c.Email, &c.Phone, &c.Balance); err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}
	if customers == nil {
		return []Customer{}, nil
	} else {
		return customers, nil
	}
}

func (c *Customer) GetCustomerByID(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT * FROM customer WHERE id=%d", c.CustomerID)
	return db.QueryRow(statement).Scan(&c.CustomerID, &c.Name, &c.Email, &c.Phone, &c.Balance)
}