package client

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
)

type MockDoer struct {
	Resp *http.Response
	Err  error
}

func (m *MockDoer) Do(req *http.Request) (*http.Response, error) {
	return m.Resp, m.Err
}

func TestListSensorsWithResponse(t *testing.T) {
	// Mock successful response
	mockResponse := `{"ok": true, "sensors": [{
		"name": "SensorXYZ",
		"lastseen": "2022-10-03T17:22:10.000001Z",
		"battery_charge": 100,
		"connected_using": "GSM",
		"wifi_password": "abcdefg",
		"online": true,
		"location": {
		  "latitude": 3.01,
		  "longitude": 2.12
		},
		"measuring_point": {
		  "name": "TheMeasuringPoint",
		  "id": 1,
		  "user_location": {
			"longitude": "8.0",
			"latitude": "12.1"
		  },
		  "active": true,
		  "swarm_type": "vibration",
		  "disable_led": false,
		  "log_flush_interval": 5,
		  "timezone": "Europe/Amsterdam",
		  "vtop_enabled": "On",
		  "atop_enabled": "Off",
		  "vector_enabled": "Off",
		  "guide_line": "DIN4150_3_80Hz",
		  "building_level": "dinFoundation",
		  "category": "CAT3",
		  "measurement_duration": 2,
		  "data_save_level": 0.2,
		  "noise_saving_enabled": "Off",
		  "vdv_enabled": "On",
		  "vdv_x": "BS6841_Wd",
		  "vdv_y": "BS6841_Wd",
		  "vdv_z": "BS6841_Wb",
		  "vdv_period": 30,
		  "trace_save_level": 20.0,
		  "trace_pre_trigger": 3.0,
		  "trace_post_trigger": 3.0,
		  "schedule_enable_1": "00:00:00",
		  "schedule_disable_1": "24:00:00",
		  "schedule_enable_2": "00:00:00",
		  "schedule_disable_2": "24:00:00",
		  "schedule_enable_3": "00:00:00",
		  "schedule_disable_3": "24:00:00",
		  "schedule_enable_4": "00:00:00",
		  "schedule_disable_4": "24:00:00",
		  "schedule_enable_5": "00:00:00",
		  "schedule_disable_5": "24:00:00",
		  "schedule_enable_6": "00:00:00",
		  "schedule_disable_6": "24:00:00",
		  "schedule_enable_0": "00:00:00",
		  "schedule_disable_0": "24:00:00",
		  "alarm_value": 50.0
		}
	  }]}`
	r := io.NopCloser(bytes.NewReader([]byte(mockResponse)))
	mockDoer := &MockDoer{
		Resp: &http.Response{
			StatusCode: http.StatusOK,
			Body:       r,
		},
	}

	mockDoer.Resp.Header = make(http.Header)
	mockDoer.Resp.Header.Set("Content-Type", "json")

	// Create client with mock doer
	client, err := NewClientWithResponses("https://api.example.com", "", WithHTTPClient(mockDoer))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Call ListSensorsWithResponse
	response, err := client.ListSensorsWithResponse(context.Background())
	if err != nil {
		t.Fatalf("ListSensorsWithResponse returned an error: %v", err)
	}

	if response.JSON200 == nil || len(*response.JSON200.Sensors) == 0 {
		t.Fatalf("Expected non-empty sensors list, got nil or empty")
	}

	if (*response.JSON200.Sensors)[0].Name == nil || *(*response.JSON200.Sensors)[0].Name != "SensorXYZ" {
		t.Fatalf("Expected sensor name to be 'SensorXYZ', got %v", *(*response.JSON200.Sensors)[0].Name)
	}
}
