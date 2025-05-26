package api

import "testing"

func TestSimulationRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     SimulationRequest
		wantErr bool
	}{
		{
			name: "valid input",
			req: SimulationRequest{
				InitialVal:  10000,
				Periods:     240,
				Simulations: 500,
				Withdrawal:  0.04,
				Method:      "normal",
				Portfolio: []AssetRequest{
					{Ticker: "AAPL", Weight: 0.6},
					{Ticker: "MSFT", Weight: 0.4},
				},
			},
			wantErr: false,
		},
		{
			name: "weight sum too low",
			req: SimulationRequest{
				InitialVal:  10000,
				Periods:     240,
				Simulations: 500,
				Withdrawal:  0.04,
				Method:      "bootstrap",
				Portfolio: []AssetRequest{
					{Ticker: "AAPL", Weight: 0.4},
					{Ticker: "MSFT", Weight: 0.4},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid method",
			req: SimulationRequest{
				InitialVal:  10000,
				Periods:     240,
				Simulations: 500,
				Withdrawal:  0.04,
				Method:      "invalid",
				Portfolio: []AssetRequest{
					{Ticker: "AAPL", Weight: 1.0},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}
