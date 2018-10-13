package types

// Event contains new data from event source or an error.
type Event struct {
	Data []byte
	Err  error
}
