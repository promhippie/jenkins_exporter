package types

// Hudson defines the root type returned by the API.
type Hudson struct {
	Mode         string   `json:"mode"`
	NumExecutors int      `json:"numExecutors"`
	Folders      []Folder `json:"jobs"`
}
