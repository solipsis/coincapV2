package coincap

import (
	"strconv"
	"time"
)

// Timestamp is wrapper around time.Time with custom unmarshaling behaviour
// specific to the format returned by the CoinCap API
type Timestamp struct {
	time.Time
}

// UnmarshalJSON implements json.Unmarshaler
// Custom unmarshaller to handle that the timestamp is not in a standard format
func (t *Timestamp) UnmarshalJSON(b []byte) error {

	// CoinCap timestamp is unix milliseconds
	m, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}
	// Convert from milleseconds to nanoseconds
	t.Time = time.Unix(0, m*1e6)
	return nil
}
