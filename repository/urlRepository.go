package repository

type UrlRepository interface {
	
	Store(string, string) error
	Get(string) (string, error)
}