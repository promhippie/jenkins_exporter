package types

// Build defines the response from specific builds.
type Build struct {
	Timestamp int64 `json:"timestamp"`
	Duration  int64 `json:"duration"`
}
