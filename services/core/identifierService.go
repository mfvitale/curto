package core

type IdentifierService interface {
	nextID() (uint64, error)
}
