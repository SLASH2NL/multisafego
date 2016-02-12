package multisafego

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var (
	msp *MultiSafePay
)

func TestMain(t *testing.M) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/json/orders/20", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
    "success": true,
    "data": {
        "transaction_id": "77962", 
        "order_id": "1291", 
        "created": "2015-02-12T14:19:38",
        "currency": "EUR",
        "amount": "235", 
        "description": "order description",
        "amount_refunded": "0",
        "status": "initialized"
    }
}`))
	})

	mux.HandleFunc("/v1/json/orders/21", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
    "success": true,
    "data": {
        "transaction_id": 77962, 
        "order_id": 1291, 
        "created": "2015-02-12T14:19:38",
        "currency": "EUR",
        "amount": 135, 
        "description": "order description",
        "amount_refunded": 0,
        "status": "initialized"
    }
}`))
	})

	mux.HandleFunc("/v1/json/orders/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
  "success" : true,
  "data" : {
    "payment_url" : "http://test.nl",
    "order_id" : "20"
  }
}`))
	})

	mux.HandleFunc("/v1/json/transactions/44", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
    "success": true,
    "data": {
        "transaction_id": 77962, 
        "order_id": "1291", 
        "created": "2015-02-12T14:19:38",
        "currency": "EUR",
        "amount": "135", 
        "description": "order description",
        "amount_refunded": 0,
        "status": "initialized"
    }
}`))
	})

	ts := httptest.NewServer(mux)

	defer ts.Close()

	url, err := url.Parse(ts.URL)
	if err == nil {
		msp = New("test", url, false)
		t.Run()
	}
}

func TestGetOrderWithStrings(t *testing.T) {
	info, err := msp.GetOrder(20)
	if err != nil {
		t.Fatal(err)
	}

	if info.Amount != 235 {
		t.Fatalf("expected 235 got:%v", info.Amount)
	}

	if info.OrderID != 1291 {
		t.Fatalf("expected 135 got:%v", info.Amount)
	}
}

func TestGetOrderWithInts(t *testing.T) {
	info, err := msp.GetOrder(21)
	if err != nil {
		t.Fatal(err)
	}

	if info.Amount != 135 {
		t.Fatalf("expected 135 got:%v", info.Amount)
	}

	if info.OrderID != 1291 {
		t.Fatalf("expected 1291 got:%v", info.OrderID)
	}
}

func TestGetPaymentWithString(t *testing.T) {
	o := NewOrder()
	o["type"] = "test"
	o["order_id"] = 20
	o["currency"] = "EUR"
	o["amount"] = 80
	o["payment_options"] = "test"
	o["description"] = "merp"

	payment, err := msp.PlaceOrder(o)
	if err != nil {
		t.Fatal(err)
	}

	if payment.OrderID != 20 {
		t.Fatalf("expected orderID 20 got: %+v", payment.OrderID)
	}
}

func TestGetTransactionWithString(t *testing.T) {
	transaction, err := msp.GetTransaction(44)
	if err != nil {
		t.Fatal(err)
	}

	if transaction.Amount != 135 {
		t.Fatalf("expected amount 135 got:%v", transaction.Amount)
	}
}
