package service

import (
	"encoding/json"
	"os"
	"sync"
	"testing"

	"github.com/fir1/port/internal/port/model"

	"github.com/stretchr/testify/require"
)

func TestNewJSONStream(t *testing.T) {
	stream := NewJSONStream()
	require.NotNil(t, stream.stream, "stream channel should not be nil")
}

func TestJSONStream_Watch(t *testing.T) {
	stream := NewJSONStream()
	ch := stream.Watch()

	require.NotNil(t, ch, "watch channel should not be nil")
}

func TestJSONStream_Start(t *testing.T) {
	// Create a temporary test file with JSON data
	data := `
{
  "AEAJM": {
    "name": "Ajman",
    "city": "Ajman",
    "country": "United Arab Emirates",
    "alias": [],
    "regions": [],
    "coordinates": [
      55.5136433,
      25.4052165
    ],
    "province": "Ajman",
    "timezone": "Asia/Dubai",
    "unlocs": [
      "AEAJM"
    ],
    "code": "52000"
  },
  "AEAUH": {
    "name": "Abu Dhabi",
    "coordinates": [
      54.37,
      24.47
    ],
    "city": "Abu Dhabi",
    "province": "Abu Z¸aby [Abu Dhabi]",
    "country": "United Arab Emirates",
    "alias": [],
    "regions": [],
    "timezone": "Asia/Dubai",
    "unlocs": [
      "AEAUH"
    ],
    "code": "52001"
  }
}
`
	tmpfile, err := os.CreateTemp("", "testfile.json")
	require.NoError(t, err, "Failed to create temporary file")
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.WriteString(data)
	require.NoError(t, err, "Failed to write data to temporary file")
	require.NoError(t, tmpfile.Close(), "Failed to close temporary file")

	expectedPorts := map[string]model.Port{
		"AEAJM": {
			Name:        "Ajman",
			City:        "Ajman",
			Country:     "United Arab Emirates",
			Alias:       []interface{}{},
			Regions:     []interface{}{},
			Coordinates: []float64{55.5136433, 25.4052165},
			Province:    "Ajman",
			Timezone:    "Asia/Dubai",
			Unlocs:      []string{"AEAJM"},
			Code:        "52000",
		},
		"AEAUH": {
			Name:        "Abu Dhabi",
			Coordinates: []float64{54.37, 24.47},
			City:        "Abu Dhabi",
			Province:    "Abu Z¸aby [Abu Dhabi]",
			Country:     "United Arab Emirates",
			Alias:       []interface{}{},
			Regions:     []interface{}{},
			Timezone:    "Asia/Dubai",
			Unlocs:      []string{"AEAUH"},
			Code:        "52001",
		},
	}

	jsonStream := NewJSONStream()

	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		for entry := range jsonStream.Watch() {
			if entry.Error != nil {
				require.Fail(t, "Unexpected error:", entry.Error.Error())
			} else {
				port, exists := expectedPorts[entry.PortCode]
				require.Truef(t, exists, "Unexpected port code: %s", entry.PortCode)

				// Compare the ports
				require.Truef(t, portsEqual(entry.Port, port), "Port mismatch for code %s", entry.PortCode)
			}
		}
	}(&wg)

	file, err := os.Open("../../../data/ports-test.json")
	require.NoError(t, err, "Failed to open file `/data/ports-test.json`")

	jsonStream.Start(file)
	wg.Wait()

	// Now close the file after it has been fully processed
	// require.NoError(t, tmpfile.Close(), "Failed to close temporary file")
}

// Helper function to compare two Port structs for equality.
func portsEqual(p1, p2 model.Port) bool {
	p1JSON, _ := json.Marshal(p1)
	p2JSON, _ := json.Marshal(p2)
	return bytes.Equal(p1JSON, p2JSON)
}
