package types

type Books struct {
	ISBN            string
	Title           string  `json:"title"`
	Publisher       *string `json:"publisher"`
	Category        *string `json:"category"`
	Author          *string `json:"author"`
	Pages           *int    `json:"pages"`
	Language        *string `json:"language"`
	PublicationYear *int    `json:"publication_year"`
	IsRead          *bool   `json:"is_read"`
	Path            string
	FileExtension   *string `json:"file_extension"`
}
