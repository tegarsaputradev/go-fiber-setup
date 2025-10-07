package models

type File struct {
	BaseModel
	Nama         string `json:"nama"`
	OriginalName string `json:"original_name"`
	MimeType     string `json:"mime_type"`
	Size         int64  `json:"size"`
	URL          string `json:"url"`
}
