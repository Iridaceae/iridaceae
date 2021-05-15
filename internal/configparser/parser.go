package configparser

import (
	"io"

	helpers2 "github.com/Iridaceae/iridaceae/internal/helpers"
)

// Parser allows us to handle read and write to and from file.
type Parser interface {
	// Marshal serialize data stream into io.Writer from our configparser
	// instance and return errors during serialization.
	Marshal(w io.Writer, c *helpers2.Settings) error

	// Unmarshal deserializes a configparser instance from io.Reader and
	// return our configparser instance and error occurred during serialization.
	Unmarshal(r io.Reader) (*helpers2.Settings, error)
}
