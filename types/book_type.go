package types

type Books struct {
	ISBN            string
	Title           string
	Publisher       *string
	Category        *string
	Author          *string
	Pages           *int
	Language        *string
	PublicationYear *int
}
