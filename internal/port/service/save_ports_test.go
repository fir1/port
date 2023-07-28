package service

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/fir1/port/config"
	"github.com/fir1/port/internal/port/model"
	"github.com/fir1/port/internal/port/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSavePortsFromFile(t *testing.T) {
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
	fmt.Printf("%+v", expectedPorts)

	cnf, err := config.NewParsedConfig()
	assert.NoError(t, err, "Unexpected error")

	portService := NewPortService(repository.NewPostRepositoryMemoryDB(), cnf)

	// Use a timeout context to limit the test duration
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testCases := []struct {
		name          string
		fileName      string
		file          multipart.File
		expectedError error
		expectedPorts map[string]model.Port
	}{
		{
			name:          "Success",
			fileName:      tmpfile.Name(),
			expectedError: nil,
			expectedPorts: map[string]model.Port{
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
			},
		},
		{
			name:          "FailureBothFilePathAndFileExist",
			fileName:      tmpfile.Name(),
			file:          tmpfile,
			expectedError: errors.New("both filePath and file cannot be provided at the same time"),
		},
		{
			name:          "FailureBothFilePathAndFileEmpty",
			fileName:      "",
			file:          nil,
			expectedError: errors.New("either filePath or file must be provided"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err = portService.SavePortsFromFile(ctx, tc.fileName, tc.file)
			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
			} else {
				assert.NoError(t, err, "Unexpected error")
			}

			if tc.expectedPorts != nil {
				actualPorts, err := portService.ListPorts(ctx)
				assert.NoError(t, err, "Unexpected error")
				assert.NotNil(t, actualPorts, "actualPorts is nil")
				assert.Len(t, actualPorts, len(tc.expectedPorts), "Unexpected number of actualPorts")

				// To test the SavePortsFromFile function with concurrent processing, we need to consider that the order of ports
				// in the actualPorts map may not be the same as in the expectedPorts map due to concurrent goroutines.
				// We can use the reflect.DeepEqual function from the reflect package to compare the maps without
				// worrying about the order of elements.
				assert.True(t, reflect.DeepEqual(expectedPorts, actualPorts), "Expected and actual ports do not match")
			}
		})
	}
}
