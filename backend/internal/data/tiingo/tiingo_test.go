package tiingo

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

// mockClient returns a fake HTTP client with canned response
func mockClient(jsonResponse string) *http.Client {
	return &http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(jsonResponse)),
				Header:     make(http.Header),
			}, nil
		}),
	}
}

// roundTripperFunc lets us use a function to implement http.RoundTripper
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

	service := &TiingoService{
		APIKey: "mock",
		Client: mockClient(mockJSON),
	}

	returns, err := service.GetMonthlyReturns("MOCK")
	require.NoError(t, err)
	require.Len(t, returns, 2)
	require.InDelta(t, 0.10, returns[0], 0.0001) // Jan -> Feb
	require.InDelta(t, 0.10, returns[1], 0.0001) // Feb -> Mar
}

func TestGetMonthlyReturns_InvalidJSON(t *testing.T) {
	badJSON := `{"not": "an array"}`
	service := &TiingoService{
		APIKey: "mock",
		Client: mockClient(badJSON),
	}

	_, err := service.GetMonthlyReturns("MOCK")
	require.Error(t, err)
}

func TestGetMonthlyReturns_HTTPError(t *testing.T) {
	client := &http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(bytes.NewBufferString("Internal Server Error")),
			}, nil
		}),
	}

	service := &TiingoService{
		APIKey: "mock",
		Client: client,
	}

	_, err := service.GetMonthlyReturns("MOCK")
	require.Error(t, err)
}
