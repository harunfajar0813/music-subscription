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

type CustomerResponse struct {
	CustomerID int `json:"customer_id"`
}

type TopUpBalanceCustomer struct {
	CustomerID int `json:"customer_id"`
	Amount     int `json:"amount"`
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

func (c *Customer) RegisterCustomer(db *sql.DB) error {
	statement := fmt.Sprintf("INSERT INTO customer(name, email, phone, balance) VALUES('%s', '%s', '%s', %d)", c.Name, c.Email, c.Phone, 0)
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}
	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&c.CustomerID)
	if err != nil {
		return err
	}
	return nil
}

func (topUpBalanceCustomer *TopUpBalanceCustomer) TopUpBalanceCustomer(db *sql.DB) error {
	statement1 := fmt.Sprintf("SELECT balance FROM customer WHERE id = %d", topUpBalanceCustomer.CustomerID)
	row, err1 := db.Query(statement1)

	if err1 != nil {
		return err1
	}

	defer row.Close()

	var oldBalanceCustomer int
	for row.Next() {
		if err2 := row.Scan(&oldBalanceCustomer); err2 != nil {
			return err2
		}
	}

	statement2 := fmt.Sprintf("UPDATE customer SET balance = %d WHERE id = %d", oldBalanceCustomer+topUpBalanceCustomer.Amount, topUpBalanceCustomer.CustomerID)
	_, err3 := db.Exec(statement2)
	return err3
}

type Transaction struct {
	TransactionID  int `json:"transaction_id"`
	CustomerID     int `json:"customer_id"`
	SubscriptionID int `json:"subscription_id"`
	Total          int `json:"total"`
}

type RenewTransaction struct {
	TransactionID int `json:"transaction_id"`
}

type TransactionResponse struct {
	TransactionID int `json:"transaction_id"`
}

func GetTransactions(db *sql.DB) ([]Transaction, error) {
	statement := fmt.Sprintf("SELECT * FROM transaction")
	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var transactions []Transaction

	for rows.Next() {
		var t Transaction
		if err := rows.Scan(&t.TransactionID, &t.CustomerID, &t.SubscriptionID, &t.Total); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	if transactions == nil {
		return []Transaction{}, nil
	} else {
		return transactions, nil
	}
}

func (t *Transaction) GetTransactionByID(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT * FROM transaction WHERE id=%d", t.TransactionID)
	return db.QueryRow(statement).Scan(&t.TransactionID, &t.CustomerID, &t.SubscriptionID, &t.Total)
}

func (t *Transaction) CreateTransaction(db *sql.DB) error {
	statement := fmt.Sprintf("INSERT INTO transaction(customer_id, subscription_id, total) VALUES(%d, %d, %d)", t.CustomerID, t.SubscriptionID, t.Total)
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}

	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&t.TransactionID)
	if err != nil {
		return err
	}
	return nil
}

func GetBalanceCustomerByID(db *sql.DB, customerID int) (int, error) {
	statement := fmt.Sprintf("SELECT balance FROM customer WHERE id=%d", customerID)
	row, err := db.Query(statement)

	if err != nil {
		return 0 - 1, err
	}

	defer row.Close()

	var oldBalance int

	for row.Next() {
		if err2 := row.Scan(&oldBalance); err2 != nil {
			return oldBalance, nil
		}
	}
	return oldBalance, errors.New("got value successfully")
}

func GetPriceSubscriptionByID(db *sql.DB, subscriptionID int) (int, error) {
	statement := fmt.Sprintf("SELECT price FROM subscription WHERE id=%d", subscriptionID)
	row, err := db.Query(statement)

	if err != nil {
		return 0 - 1, err
	}

	defer row.Close()

	var price int

	for row.Next() {
		if err2 := row.Scan(&price); err2 != nil {
			return price, nil
		}
	}
	return price, errors.New("got value successfully")
}

func (t *Transaction) DebitCustomerBalance(db *sql.DB, oldBalance int, subsPrice int) error {
	if oldBalance-subsPrice < 0 {
		return errors.New("balance wasn't enough")
	} else {
		statement := fmt.Sprintf("UPDATE customer SET balance = %d WHERE id = %d", oldBalance-subsPrice, t.CustomerID)
		_, err := db.Exec(statement)

		if err != nil {
			return err
		}
		return nil
	}
}
