package apis

import (
	"encoding/json"
	"testing"
)

func TestTracking(t *testing.T) {
	res, err := NewLogisticsApi().GetTracking("3006507760217")
	if err != nil {
		t.Error(err)
	}

	jsonData, _ := json.Marshal(res)
	print(string(jsonData))
}
