package models

// File represents a file stored in the database.
type File struct {
	// Name is the name of the file.
	Name string `json:"name"`
	// ContentType is the MIME type of the file.
	ContentType string `json:"content_type"`
	// Data is the file data.
	Data []byte `json:"data"`
}
