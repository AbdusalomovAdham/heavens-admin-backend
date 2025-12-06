package entity

type File struct {
	Id           int    `json:"id"`
	Path         string `json:"path"`
	Size         string `json:"size"`
	SavedName    string `json:"saved_name"`
	OriginalName string `json:"original_name"`
}
