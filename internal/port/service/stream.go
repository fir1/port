package service

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/fir1/port/internal/port/model"
)

// Entry represents each stream. If the stream fails, an error will be present.
type Entry struct {
	Error    error
	PortCode string
	Port     model.Port
}

// Stream helps transmit each streams within a channel.
type Stream struct {
	stream chan Entry
}

func NewJSONStream() Stream {
	return Stream{
		stream: make(chan Entry),
	}
}

// Watch watches JSON streams. Each stream entry will either have an error or a
// Post object. Client code does not need to explicitly exit after catching an
// error as the `Start` method will close the channel automatically.
func (s Stream) Watch() <-chan Entry {
	return s.stream
}

// Start starts streaming JSON file line by line. If an error occurs, the channel
// will be closed.
// To handle large JSON files and limited resources efficiently, we can use a streaming-based
// approach to read the file in chunks rather than reading the entire file at once.
// It allows us to handle large JSON files without loading the entire file into memory.
// This way, we can process the JSON data in smaller portions and reduce memory usage.
func (s Stream) Start(file io.Reader) {
	// Stop streaming channel as soon as nothing left to read in the file.
	defer close(s.stream)

	decoder := json.NewDecoder(file)

	// Check for the opening curly braces to start the JSON data
	tok, err := decoder.Token()
	if err != nil {
		s.stream <- Entry{Error: fmt.Errorf("decode opening delimiter: %w", err)}
		return
	}
	if delim, ok := tok.(json.Delim); !ok || delim != '{' {
		s.stream <- Entry{Error: fmt.Errorf("expected '{' at the beginning of JSON data")}
		return
	}

	// Read file content as long as there is something.
	i := 1
	for decoder.More() {
		portCode, err := decoder.Token()
		if err != nil {
			s.stream <- Entry{Error: fmt.Errorf("decode line %d: %w", i, err)}
			return
		}

		var port model.Port
		err = decoder.Decode(&port)
		if err != nil {
			s.stream <- Entry{Error: fmt.Errorf("decode line %d: %w", i, err)}
			return
		}

		s.stream <- Entry{Port: port, PortCode: portCode.(string)}
		i++
	}

	// Check for the closing curly braces to end the JSON data
	tok, err = decoder.Token()
	if err != nil {
		s.stream <- Entry{Error: fmt.Errorf("decode closing delimiter: %w", err)}
		return
	}

	if delim, ok := tok.(json.Delim); !ok || delim != '}' {
		s.stream <- Entry{Error: fmt.Errorf("expected '}' at the end of JSON data")}
		return
	}
}
