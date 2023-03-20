package types

// Folder is a simple type used for folder listings.
type Folder struct {
	Class   string   `json:"_class"`
	Name    string   `json:"name"`
	URL     string   `json:"url"`
	Folders []Folder `json:"jobs"`
}
