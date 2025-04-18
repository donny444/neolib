package types

type Books struct {
	UUID            string
	Title           string
	Publisher       *string
	Category        *string
	Author          *string
	Page            *int
	Language        *string
	PublicationYear *int
	ISBN            string
}
