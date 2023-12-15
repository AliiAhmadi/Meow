package data

import (
	"fmt"
	"strconv"
)

// Declare a custom Runtime type which has the int32.
type Runtime int32

// Implement a MarshalJSON() method on the Runtime type so that it satisfies the
// json.Marshaler interface. This should return the JSON-encoded value for the movie
// runtime (in our case, it will return a string in the format "<runtime> mins").
func (rt Runtime) MarshalJSON() ([]byte, error) {

	// Create required format.
	jsonValue := fmt.Sprintf("%d mins", rt)

	// Use the strconv.Quote() function on the string to wrap it in double quotes. It
	// needs to be surrounded by double quotes in order to be a valid *JSON string*.
	quotedJSONValue := strconv.Quote(jsonValue)

	// Convert from string to an slice of bytes.
	return []byte(quotedJSONValue), nil
}
