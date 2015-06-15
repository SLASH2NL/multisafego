package multisafego

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reflect"
)

var (
	// TESTURL is the url to use in test mode
	TESTURL *url.URL
	// PRODURL is the url to use in production mode
	PRODURL *url.URL
)

func init() {
	TESTURL, _ = url.Parse("https://testapi.multisafepay.com")
	PRODURL, _ = url.Parse("https://api.multisafepay.com")
}

// MultiSafePay is the wrapper for calls to the Api.
// Use the New function to construct a basic variant
type MultiSafePay struct {
	apiKey  string
	baseURL *url.URL
	debug   bool
	Logs    chan errorLog
}

// A APIError occurs when the response has a non 200 status code
// default error values like a timeout will be wrapped in an APIError with code -1
type APIError struct {
	Code    int         `json:"error_code"`
	Message string      `json:"error_info"`
	Data    interface{} `json:"data"`
}

func (a APIError) String() string {
	return fmt.Sprintf("[%d] %s", a.Code, a.Message)
}

type errorLog struct {
	URL      string
	Request  []byte
	Response []byte
}

// New returns a basic multisafepay object that can be set in debug mode to dump requests/responses
func New(apiKey string, baseURL *url.URL, debug bool) *MultiSafePay {
	return &MultiSafePay{apiKey, baseURL, debug, make(chan errorLog, 100)}
}

func errorToAPIError(err error) *APIError {
	return &APIError{-1, err.Error(), nil}
}

// Path is a shortcut for prepending /v1/json to the path
func Path(path string) string {
	return "/v1/json" + path
}

func (m *MultiSafePay) log(url string, request, response []byte) {
	if len(m.Logs) == cap(m.Logs) {
		<-m.Logs
	}

	m.Logs <- errorLog{url, request, response}
}

// Execute will do a json call to multisafepay
// payload represents the request body that will be encoded as json
// the returnVal will be transformed the same way json.Unmarshal does.
// If an error occurs an APIError will be returned
func (m *MultiSafePay) Execute(url *url.URL, method string, payload interface{}, returnVal interface{}) *APIError {
	requestBody := &bytes.Buffer{}
	if payload != nil {
		encoder := json.NewEncoder(requestBody)
		err := encoder.Encode(payload)

		if err != nil {
			return errorToAPIError(err)
		}
	}

	res, err := m.Call(url, method, requestBody)
	if err != nil {
		return errorToAPIError(err)
	}

	// Copy the response body in buffer to re-read the contents
	buff := &bytes.Buffer{}
	io.Copy(buff, res.Body)
	res.Body.Close()

	// If we get a non 200 status code try to parse the response into APIError
	if res.StatusCode != http.StatusOK {
		var apiErr *APIError

		// unmarshall the response and check if we have an error
		if err = json.Unmarshal(buff.Bytes(), &apiErr); err == nil && apiErr.Code != 0 {
			return apiErr
		}

		return errorToAPIError(fmt.Errorf("Expected status code %d got %d but could not parse an ApiErr for the response body: %s", http.StatusOK, res.StatusCode, string(buff.Bytes())))
	}

	var checkSuccess struct {
		Success bool            `json:"success"`
		Data    json.RawMessage `json:"data"`
	}

	err = json.Unmarshal(buff.Bytes(), &checkSuccess)
	if err != nil {
		return errorToAPIError(fmt.Errorf("Could not parse response: %s into:%s got error:", string(buff.Bytes()), reflect.TypeOf(returnVal).Name(), err))
	}

	if checkSuccess.Success == false {
		return errorToAPIError(fmt.Errorf("The response returned a 200 OK but contained a success=false for response: %s", string(buff.Bytes())))
	}

	if returnVal != nil {
		err = json.Unmarshal(checkSuccess.Data, returnVal)
		if err != nil {
			return errorToAPIError(fmt.Errorf("Could not parse data field: %s into:%s got error:", string(buff.Bytes()), reflect.TypeOf(returnVal).Name(), err))
		}
	}

	return nil
}

// Call can be used to get the raw unparsed http.Response.
// Remember to close the response body when using this function.
func (m *MultiSafePay) Call(url *url.URL, method string, payload io.Reader) (*http.Response, error) {

	req, err := http.NewRequest(method, url.String(), payload)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header["api_key"] = []string{m.apiKey}

	var requestDump []byte
	if m.debug {
		requestDump, _ = httputil.DumpRequest(req, true)
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	req.Body.Close()

	if m.debug {
		b, _ := httputil.DumpResponse(res, true)
		m.log(url.String(), requestDump, b)
	}

	return res, err
}
