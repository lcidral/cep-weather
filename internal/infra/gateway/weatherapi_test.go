package gateway

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// customTransport is a http.RoundTripper that redirects all requests to a test server
type customTransport struct {
	serverURL string
}

func (t *customTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Create a new request to the test server
	testReq, err := http.NewRequest(req.Method, t.serverURL, req.Body)
	if err != nil {
		return nil, err
	}


	// Add all query parameters to the test request
	testQuery := testReq.URL.Query()
	for key, values := range req.URL.Query() {
		for _, value := range values {
			testQuery.Add(key, value)
		}
	}
	testReq.URL.RawQuery = testQuery.Encode()

	// Copy the headers
	testReq.Header = req.Header

	// Send the request to the test server
	return http.DefaultTransport.RoundTrip(testReq)
}

func TestWeatherAPIGateway_GetTemperature_URLEncoding(t *testing.T) {
	// Create a test server that verifies the URL is properly encoded
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that the query parameter is properly decoded
		query := r.URL.Query().Get("q")
		if query != "São Francisco do Sul" {
			t.Errorf("Expected decoded city name 'São Francisco do Sul', got: %s", query)
		}

		// Return a valid response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"location": {
				"name": "São Francisco do Sul",
				"region": "Santa Catarina",
				"country": "Brazil"
			},
			"current": {
				"temp_c": 25.0,
				"temp_f": 77.0
			}
		}`))
	}))
	defer server.Close()

	// Create a custom HTTP client that redirects all requests to our test server
	client := &http.Client{
		Transport: &customTransport{
			serverURL: server.URL,
		},
	}

	// Create the gateway with our custom client
	gateway := NewWeatherAPIGateway(client, "test-api-key")

	// Test with a city name containing special characters
	temp, err := gateway.GetTemperature("São Francisco do Sul")

	// Verify the result
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if temp != 25.0 {
		t.Errorf("Expected temperature 25.0, got: %f", temp)
	}
}
