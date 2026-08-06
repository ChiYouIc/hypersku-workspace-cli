package apis

import "testing"

func TestWarehouse_GetWarehouseTracking(t *testing.T) {
	tests := []struct {
		name           string
		trackingNumber string
		wantErr        bool
		wantLen        int
	}{
		{name: "Test Case 1", trackingNumber: "79022253801788", wantErr: false, wantLen: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := NewWarehouseApi()
			got, gotErr := api.GetWarehouseTracking(tt.trackingNumber)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetWarehouseTracking() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetWarehouseTracking() succeeded unexpectedly")
			}

			if len(got) != tt.wantLen {
				t.Errorf("GetWarehouseTracking() = %v, want %v", len(got), tt.wantLen)
			}
		})
	}
}
