package core

type IdentifierService interface {
    NextID() (uint64, error)
}
