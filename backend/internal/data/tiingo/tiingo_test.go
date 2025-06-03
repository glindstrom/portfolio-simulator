package tiingo

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// mockClient returns a fake HTTP client with canned response.
func mockClient(statusCode int, jsonResponse string) *http.Client {
	return &http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: statusCode,
				Body:       io.NopCloser(bytes.NewBufferString(jsonResponse)),
				Header:     make(http.Header),
			}, nil
		}),
	}
}

// roundTripperFunc lets us use a function to implement http.RoundTripper.
type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestGetMonthlyReturns_Mocked(t *testing.T) {
	mockJSON := `[{
		"date": "2023-01-31T00:00:00.000Z",
		"adjClose": 100.0
	}, {
		"date": "2023-02-28T00:00:00.000Z",
		"adjClose": 110.0
	}, {
		"date": "2023-03-31T00:00:00.000Z",
		"adjClose": 121.0
	}]`

	service := &Service{
		APIKey: "mock_api_key",
		Client: mockClient(http.StatusOK, mockJSON),
	}

	returns, err := service.GetMonthlyReturns("MOCK")
	require.NoError(t, err)
	require.Len(t, returns, 2)
	if len(returns) == 2 {
		require.InDelta(t, 0.10, returns[0], 0.0001)
		require.InDelta(t, 0.10, returns[1], 0.0001)
	}
}

func TestGetMonthlyReturns_InvalidJSON(t *testing.T) {
	// Removed unused 'badJSON' variable.
	// Using malformedJSON to ensure a parsing error within the data processing chain.
	malformedJSON := `[{"date": "2023-01-31T00:00:00.000Z", "adjClose": "not-a-float"}]`

	service := &Service{
		APIKey: "mock_api_key",
		Client: mockClient(http.StatusOK, malformedJSON),
	}

	_, err := service.GetMonthlyReturns("MOCK_INVALID_JSON_DATA")
	require.Error(t, err)
}

func TestGetMonthlyReturns_HTTPError(t *testing.T) {
	service := &Service{
		APIKey: "mock_api_key",
		Client: mockClient(http.StatusInternalServerError, `{"error": "Internal Server Error"}`),
	}

	_, err := service.GetMonthlyReturns("MOCK_HTTP_ERROR")
	require.Error(t, err)
}

func TestGetMonthlyReturns_NoAPIKey(t *testing.T) {
	service := &Service{
		APIKey: "", // No API Key
		Client: mockClient(http.StatusOK, "[]"),
	}

	_, err := service.GetMonthlyReturns("MOCK_NO_KEY")
	require.Error(t, err)
	require.Contains(t, err.Error(), "API key is not configured")
}

func TestGetMonthlyReturns_NoDataReturned(t *testing.T) {
	mockJSON := `[]` // Empty array of prices

	service := &Service{
		APIKey: "mock_api_key",
		Client: mockClient(http.StatusOK, mockJSON),
	}

	_, err := service.GetMonthlyReturns("MOCK_NO_DATA")
	require.Error(t, err)
	// Test asserts that GetMonthlyPrices returns an error if Tiingo provides no price data.
	require.Contains(t, err.Error(), "no price data returned")
}

func TestGetMonthlyReturns_InsufficientDataForReturns(t *testing.T) {
	mockJSON := `[{
		"date": "2023-01-31T00:00:00.000Z",
		"adjClose": 100.0
	}]`

	service := &Service{
		APIKey: "mock_api_key",
		Client: mockClient(http.StatusOK, mockJSON),
	}

	returns, err := service.GetMonthlyReturns("MOCK_ONE_POINT")
	require.NoError(t, err)
	require.Nil(t, returns) // common.ToMonthlyReturns returns nil for < 2 prices.
}
