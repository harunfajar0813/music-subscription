package music

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const createSubscriptionTable = `
CREATE TABLE IF NOT EXISTS subscriptions
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

	code := m.Run()

	clearSubscriptionTable()

	os.Exit(code)
}

func ensureSubscriptionTableExists() {
	if _, err := a.DB.Exec(createSubscriptionTable); err != nil {
		log.Fatal(err)
	}
}

func clearSubscriptionTable() {
	a.DB.Exec("DELETE FROM subscriptions")
	a.DB.Exec("ALTER TABLE subscriptions AUTO_INCREMENT = 1")
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
