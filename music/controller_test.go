package music

import (
	"bytes"
	"encoding/json"
	"fmt"
	generateFake "github.com/bxcodec/faker"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const createSubscriptionTable = `
CREATE TABLE IF NOT EXISTS subscription
(
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    price INT NOT NULL,
	duration INT NOT NULL
)`

var a App

func TestMain(m *testing.M) {
	a = App{}
	a.Initialize("root", "", "bcc_music")

	ensureSubscriptionTableExists()
	ensureCustomerTableExists()
	ensureTransactionTableExists()

	code := m.Run()

	clearSubscriptionTable()
	clearCustomerTable()
	clearTransactionTable()

	os.Exit(code)
}

func ensureSubscriptionTableExists() {
	if _, err := a.DB.Exec(createSubscriptionTable); err != nil {
		log.Fatal(err)
	}
}

func clearSubscriptionTable() {
	a.DB.Exec("DELETE FROM subscription")
	a.DB.Exec("ALTER TABLE subscription AUTO_INCREMENT = 1")
}

func TestEmptySubscriptionTable(t *testing.T) {
	clearSubscriptionTable()

	req, _ := http.NewRequest("GET", "/api/subscriptions", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistentSubscription(t *testing.T) {
	clearSubscriptionTable()
	req, _ := http.NewRequest("GET", "/api/subscription/45", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Subscription not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Subscription not found'. Got '%s'", m["error"])
	}
}

func TestCreateSubscription(t *testing.T) {
	clearSubscriptionTable()

	payload := []byte(`{"name":"test subscription","price":30,"duration":10}`)

	req, _ := http.NewRequest("POST", "/api/subscription", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "test subscription" {
		t.Errorf("Expected subscription's name to be 'test subscription'. Got '%v'", m["name"])
	}

	if m["price"] != 30.0 {
		t.Errorf("Expected subscription's price to be '30'. Got '%v'", m["price"])
	}

	if m["duration"] != 10.0 {
		t.Errorf("Expected subscription's duration to be '10'. Got '%v'", m["duration"])
	}

	if m["subscription_id"] != 1.0 {
		t.Errorf("Expected subscription ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetSubscription(t *testing.T) {
	clearSubscriptionTable()
	addSubscription(1)
	req, _ := http.NewRequest("GET", "/api/subscription/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func addSubscription(count int) {
	if count < 1 {
		count = 1
	}
	for i := 0; i < count; i++ {
		statement := fmt.Sprintf("INSERT INTO subscription(name, price, duration) VALUES('%s', %d, %d)", generateFake.FirstName(), (i+1)*10, (i+1)*10)
		a.DB.Exec(statement)
	}
}

const createCustomerTable = `
CREATE TABLE IF NOT EXISTS customer
(
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    email VARCHAR(50) NOT NULL,
	phone VARCHAR(50) NOT NULL,
	balance INT NOT NULL
)`

func ensureCustomerTableExists() {
	if _, err := a.DB.Exec(createCustomerTable); err != nil {
		log.Fatal(err)
	}
}

func clearCustomerTable() {
	a.DB.Exec("DELETE FROM customer")
	a.DB.Exec("ALTER TABLE customer AUTO_INCREMENT = 1")
}

func TestEmptyCustomerTable(t *testing.T) {
	clearCustomerTable()

	req, _ := http.NewRequest("GET", "/api/customers", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistentCustomer(t *testing.T) {
	clearCustomerTable()

	req, _ := http.NewRequest("GET", "/api/customer/45", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Customer not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Customer not found'. Got '%s'", m["error"])
	}
}

func TestRegisterCustomer(t *testing.T) {
	clearCustomerTable()

	payload := []byte(`{"name":"test customer","email":"testcustomer@gmail.com","phone":"0813"}`)

	req, _ := http.NewRequest("POST", "/api/customer/register", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "test customer" {
		t.Errorf("Expected customer's name to be 'test customer'. Got '%v'", m["name"])
	}

	if m["email"] != "testcustomer@gmail.com" {
		t.Errorf("Expected customer's email to be 'testcustomer@gmail.com'. Got '%v'", m["price"])
	}

	if m["phone"] != "0813" {
		t.Errorf("Expected customer's phone to be '0813'. Got '%v'", m["duration"])
	}

	if m["customer_id"] != 1.0 {
		t.Errorf("Expected customer's ID to be '1'. Got '%v'", m["id"])
	}
}

func TestTopUpBalanceCustomer(t *testing.T) {
	clearCustomerTable()

	payload := []byte(`{"customer_id":1,"amount":10}`)

	addCustomer(1)

	req, _ := http.NewRequest("PUT", "/api/customer/topup", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["customer_id"] != 1.0 {
		t.Errorf("Expected customer's ID to be '1'. Got '%v'", m["customer_id"])
	}

	if m["amount"] != 10.0 {
		t.Errorf("Expected customer's amount to be '10'. Got '%v'", m["amount"])
	}
}

func TestGetCustomer(t *testing.T) {
	clearCustomerTable()
	addCustomer(1)
	req, _ := http.NewRequest("GET", "/api/customer/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func addCustomer(count int) {
	if count < 1 {
		count = 1
	}
	for i := 0; i < count; i++ {
		statement := fmt.Sprintf("INSERT INTO customer(name, email, phone, balance) VALUES('%s', '%s', '%s', %d)", generateFake.FirstName(), generateFake.Email(), generateFake.Phonenumber(), 0)
		a.DB.Exec(statement)
	}
}

const createTransactionTable = `
CREATE TABLE IF NOT EXISTS transaction
(
   id INT AUTO_INCREMENT PRIMARY KEY,
   customer_id INT NOT NULL,
   subscription_id INT NOT NULL,
	total INT NOT NULL,
	FOREIGN KEY (customer_id) REFERENCES customer(id),
	FOREIGN KEY (subscription_id) REFERENCES subscription(id)
)`

func ensureTransactionTableExists() {
	if _, err := a.DB.Exec(createTransactionTable); err != nil {
		log.Fatal(err)
	}
}

func clearTransactionTable() {
	a.DB.Exec("DELETE FROM transaction")
	a.DB.Exec("ALTER TABLE transaction AUTO_INCREMENT = 1")

	a.DB.Exec("DELETE FROM customer")
	a.DB.Exec("ALTER TABLE customer AUTO_INCREMENT = 1")

	a.DB.Exec("DELETE FROM subscription")
	a.DB.Exec("ALTER TABLE subscription AUTO_INCREMENT = 1")
}

func TestEmptyTransactionTable(t *testing.T) {
	clearTransactionTable()

	req, _ := http.NewRequest("GET", "/api/transactions", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistentTransaction(t *testing.T) {
	clearTransactionTable()

	req, _ := http.NewRequest("GET", "/api/transaction/45", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Transaction not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Transaction not found'. Got '%s'", m["error"])
	}
}

func TestCreateTransaction(t *testing.T) {
	clearTransactionTable()

	payload := []byte(`{"customer_id":1,"subscription_id":1,"total":10}`)

	addCustomer(1)
	addSubscription(1)
	req, _ := http.NewRequest("POST", "/api/transaction/payment", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["customer_id"] != 1.0 {
		t.Errorf("Expected transaction's customer_id to be '1'. Got '%v'", m["name"])
	}

	if m["subscription_id"] != 1.0 {
		t.Errorf("Expected transaction's subscription_id to be '1'. Got '%v'", m["price"])
	}

	if m["total"] != 10.0 {
		t.Errorf("Expected transaction's duration to be '10'. Got '%v'", m["duration"])
	}

	if m["transaction_id"] != 1.0 {
		t.Errorf("Expected transaction's ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetTransaction(t *testing.T) {
	clearTransactionTable()

	addCustomer(1)
	addSubscription(1)
	addTransaction(1)

	req, _ := http.NewRequest("GET", "/api/transaction/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func addTransaction(count int) {
	if count < 1 {
		count = 1
	}
	for i := 1; i <= count; i++ {
		statement := fmt.Sprintf("INSERT INTO transaction(customer_id, subscription_id, total) VALUES(%d, %d, %d)", i, i, (i+1)*10)
		a.DB.Exec(statement)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
