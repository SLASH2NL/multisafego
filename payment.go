package multisafego

import (
	"errors"
	"strconv"
)

// Order is the container for holding all the order parametersa that should be posted
type Order map[string]interface{}

// Payment is the response from PlaceOrder
type Payment struct {
	URL     string `json:"payment_url"`
	OrderID int    `json:"order_id"`
}

// OrderInfo is returned from the GetOrder function and contains the TransactionId which can be
// used to grab the payment status
type OrderInfo struct {
	Amount         int `json:"amount"`
	AmountRefunded int `json:"amount_refunded"`
	//Created        time.Time `json:"created"`
	Currency string `json:"currency"`
	Customer struct {
		Email  string `json:"email"`
		Locale string `json:"locale"`
	} `json:"customer"`
	Description string `json:"description"`
	//Modified       time.Time `json:"modified"`
	OrderID        int `json:"order_id"`
	PaymentDetails struct {
		AccountHolderName     interface{} `json:"account_holder_name"`
		AccountID             interface{} `json:"account_id"`
		ExternalTransactionID interface{} `json:"external_transaction_id"`
		RecurringID           interface{} `json:"recurring_id"`
		Type                  string      `json:"type"`
	} `json:"payment_details"`
	Status        string `json:"status"`
	TransactionID int    `json:"transaction_id"`
}

// Transaction contains a transaction for an order
type Transaction struct {
	Amount   int    `json:"amount"`
	Created  string `json:"created"`
	Currency string `json:"currency"`
	Customer struct {
		Email string `json:"email"`
	} `json:"customer"`
	Description    string `json:"description"`
	OrderID        int    `json:"order_id"`
	OrderStatus    string `json:"order_status"`
	PaymentDetails struct {
		AccountHolderName     interface{} `json:"account_holder_name"`
		AccountID             interface{} `json:"account_id"`
		ExternalTransactionID string      `json:"external_transaction_id"`
		RecurringID           interface{} `json:"recurring_id"`
		Type                  string      `json:"type"`
	} `json:"payment_details"`
	Status        string  `json:"status"`
	TransactionID float64 `json:"transaction_id"`
	Type          string  `json:"type"`
}

// IsCompleted returns true if a payment is succesfull
func (t *Transaction) IsCompleted() bool {
	return t.Status == "completed"
}

// NewOrder is a shortcut method for creating the order map
func NewOrder() Order {
	o := make(Order)
	return o
}

// SetIssuer for the order to enable direct payments
func (o Order) SetIssuer(gateway string, issuer string) {
	o["gateway"] = gateway
	o["gateway_info"] = map[string]interface{}{
		"issuer_id": issuer,
	}
}

// SetPaymentOptions sets the urls for the payment
func (o Order) SetPaymentOptions(notifyURL, redirectURL, cancelURL string, closeWindow bool) {
	options := make(map[string]interface{})

	options["notification_url"] = notifyURL
	options["redirect_url"] = redirectURL
	options["cancel_url"] = cancelURL
	options["close_window"] = closeWindow

	o["payment_options"] = options
}

// PlaceOrder will post an order and return the response urls
func (m *MultiSafePay) PlaceOrder(o Order) (*Payment, *APIError) {
	m.baseURL.Path = Path("/orders/")

	mandatoryParameters := []string{
		"type",
		"order_id",
		"currency",
		"amount",
		"payment_options",
		"description",
	}

	for _, param := range mandatoryParameters {
		if _, ok := o[param]; !ok {
			return nil, errorToAPIError(errors.New(param + " is a required parameter"))
		}
	}

	var x Payment
	err := m.Execute(m.baseURL, "POST", o, &x)
	return &x, err
}

// GetOrder returns info about an order placed
func (m *MultiSafePay) GetOrder(id int) (*OrderInfo, *APIError) {
	m.baseURL.Path = Path("/orders/" + strconv.Itoa(id))

	var x OrderInfo
	err := m.Execute(m.baseURL, "GET", nil, &x)
	return &x, err
}

// GetTransaction returns infomation about the payment
func (m *MultiSafePay) GetTransaction(transactionID int) (*Transaction, *APIError) {
	m.baseURL.Path = Path("/transactions/" + strconv.Itoa(transactionID))

	var x Transaction
	err := m.Execute(m.baseURL, "GET", nil, &x)
	return &x, err
}
