package music

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("%s:%s@/%s", user, password, dbname)

	var err error
	a.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(address string) {
	log.Fatal(http.ListenAndServe(address, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/api/subscriptions", a.getSubscriptions).Methods("GET")
	a.Router.HandleFunc("/api/subscription/{id:[0-9]+}", a.getSubscriptionByID).Methods("GET")
	a.Router.HandleFunc("/api/subscription", a.createSubscription).Methods("POST")

	a.Router.HandleFunc("/api/customers", a.getCustomers).Methods("GET")
	a.Router.HandleFunc("/api/customer/{id:[0-9]+}", a.getCustomerByID).Methods("GET")
	a.Router.HandleFunc("/api/customer/register", a.registerCustomer).Methods("POST")
	a.Router.HandleFunc("/api/customer/topup", a.topUpBalanceCustomer).Methods("PUT")

	a.Router.HandleFunc("/api/transactions", a.getTransactions).Methods("GET")
	a.Router.HandleFunc("/api/transaction/{id:[0-9]+}", a.getTransactionByID).Methods("GET")
	a.Router.HandleFunc("/api/transaction/payment", a.createTransaction).Methods("POST")
}

func (a *App) getSubscriptions(w http.ResponseWriter, r *http.Request) {
	subscriptions, err := GetSubscriptions(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, subscriptions)
}

func (a *App) getSubscriptionByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid subscription ID")
		return
	}

	s := Subscription{SubscriptionID: id}
	if err := s.GetSubscriptionByID(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Subscription not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, s)
}

func (a *App) createSubscription(w http.ResponseWriter, r *http.Request) {
	var s Subscription
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&s); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	if err := s.CreateSubscription(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, s)
}

func (a *App) getCustomers(w http.ResponseWriter, r *http.Request) {
	customers, err := GetCustomers(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, customers)
}

func (a *App) getCustomerByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	c := Customer{CustomerID: id}
	if err := c.GetCustomerByID(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Customer not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, c)
}

func (a *App) registerCustomer(w http.ResponseWriter, r *http.Request) {
	var c Customer
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&c); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	if err := c.RegisterCustomer(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, c)
}

func (a *App) topUpBalanceCustomer(w http.ResponseWriter, r *http.Request) {
	var topUpBalance TopUpBalanceCustomer
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&topUpBalance); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	if err := topUpBalance.TopUpBalanceCustomer(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, topUpBalance)
}

func (a *App) getTransactions(w http.ResponseWriter, r *http.Request) {
	transactions, err := GetTransactions(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, transactions)
}

func (a *App) getTransactionByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid transaction ID")
		return
	}

	t := Transaction{TransactionID: id}
	if err := t.GetTransactionByID(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Transaction not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, t)
}

func (a *App) createTransaction(w http.ResponseWriter, r *http.Request) {
	var t Transaction
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	oldBalance, _ := t.GetBalanceCustomerByID(a.DB)
	subsPrice, _ := t.GetPriceSubscriptionByID(a.DB)
	errDecreased := t.DecreasedCustomerBalance(a.DB, oldBalance, subsPrice)
	if errDecreased != nil {
		respondWithError(w, http.StatusInternalServerError, errDecreased.Error())
	} else {
		if err := t.CreateTransaction(a.DB); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusCreated, t)
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
