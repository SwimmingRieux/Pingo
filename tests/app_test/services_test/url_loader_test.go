package services_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"pingo/configs"
	"pingo/internal/app/services"
	"testing"
)

var configForUrlLoaderTest, _ = configs.NewConfig()

type urlLoaderOkResponseTest struct {
	name         string
	testFunction func(t *testing.T)
}

var urlLoaderOkResponseTests = []urlLoaderOkResponseTest{
	{
		name:         "should return valid text when url is valid",
		testFunction: urlLoaderValidTest,
	},
	{
		name:         "should return err when can't read response",
		testFunction: urlLoaderCantReadTest,
	},
}

type urlLoaderBadResponseTest struct {
	name         string
	responseCode int
	testFunction func(t *testing.T, responseCode int)
}

var urlLoaderBadResponseTests = []urlLoaderBadResponseTest{
	// Client Errors (4xx)
	{"should return 'Bad Request' when server responds with 400", http.StatusBadRequest, urlLoaderReturnsBadResponseTest},
	{"should return 'Unauthorized' when server responds with 401", http.StatusUnauthorized, urlLoaderReturnsBadResponseTest},
	{"should return 'Forbidden' when server responds with 403", http.StatusForbidden, urlLoaderReturnsBadResponseTest},
	{"should return 'Not Found' when server responds with 404", http.StatusNotFound, urlLoaderReturnsBadResponseTest},
	{"should return 'Method Not Allowed' when server responds with 405", http.StatusMethodNotAllowed, urlLoaderReturnsBadResponseTest},
	{"should return 'Request Timeout' when server responds with 408", http.StatusRequestTimeout, urlLoaderReturnsBadResponseTest},
	{"should return 'Conflict' when server responds with 409", http.StatusConflict, urlLoaderReturnsBadResponseTest},
	{"should return 'Gone' when server responds with 410", http.StatusGone, urlLoaderReturnsBadResponseTest},
	{"should return 'Too Many Requests' when server responds with 429", http.StatusTooManyRequests, urlLoaderReturnsBadResponseTest},

	// Server Errors (5xx)
	{"should return 'Internal Server Error' when server responds with 500", http.StatusInternalServerError, urlLoaderReturnsBadResponseTest},
	{"should return 'Not Implemented' when server responds with 501", http.StatusNotImplemented, urlLoaderReturnsBadResponseTest},
	{"should return 'Bad Gateway' when server responds with 502", http.StatusBadGateway, urlLoaderReturnsBadResponseTest},
	{"should return 'Service Unavailable' when server responds with 503", http.StatusServiceUnavailable, urlLoaderReturnsBadResponseTest},
	{"should return 'Gateway Timeout' when server responds with 504", http.StatusGatewayTimeout, urlLoaderReturnsBadResponseTest},
}

func TestLoad(t *testing.T) {
	t.Parallel()
	for _, testCase := range urlLoaderOkResponseTests {
		t.Run(testCase.name, testCase.testFunction)
	}
	for _, testCase := range urlLoaderBadResponseTests {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.testFunction(t, testCase.responseCode)
		})
	}
}

func urlLoaderValidTest(t *testing.T) {
	t.Parallel()
	// Arrange
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("mock response"))
	}))
	defer mockServer.Close()

	loader := services.NewUrlLoader(configForUrlLoaderTest)

	// Act
	result, err := loader.Load(mockServer.URL)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	assert.Equal(t, "mock response", result)
}

func urlLoaderCantReadTest(t *testing.T) {
	t.Parallel()
	// Arrange
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		conn, _, _ := w.(http.Hijacker).Hijack()
		conn.Close() // Force an early close to simulate a read error
	}))
	defer mockServer.Close()

	loader := services.NewUrlLoader(configForUrlLoaderTest)

	// Act
	_, err := loader.Load(mockServer.URL)

	// Assert
	assert.Error(t, err)
}

func urlLoaderReturnsBadResponseTest(t *testing.T, responseCode int) {
	t.Parallel()
	// Arrange
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(responseCode)
	}))
	defer mockServer.Close()

	loader := services.NewUrlLoader(configForUrlLoaderTest)

	// Act
	_, err := loader.Load(mockServer.URL)

	// Assert
	if err == nil {
		t.Fatalf("expected an error but got none")
	}
	expectedError := fmt.Errorf("%v %v", configForUrlLoaderTest.Errors.HttpStatus, responseCode)
	assert.EqualError(t, expectedError, err.Error())
}
