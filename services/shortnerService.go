package services

import (
    "fmt"

    "github.com/mfvitale/curto/repository"
    "github.com/mfvitale/curto/services/core"
    log "github.com/sirupsen/logrus"
)

type ShortenOperationError struct {
    message string
    detail string
}

func (e ShortenOperationError) Error() string {
    return fmt.Sprintf("%s: %s", e.message, e.detail)
}

type ShortnerService struct {
    redisRepo repository.UrlRepository
    identifierService core.IdentifierService
}

func NewShortnerService(urlRepository repository.UrlRepository,
    identifierService core.IdentifierService) ShortnerService {
    return ShortnerService{urlRepository, identifierService}
}

func (s *ShortnerService) Shorten(url string) (string, error) {

    id, err := s.identifierService.NextID()
    if err != nil {
        log.Error(fmt.Sprintf("Error while getting identifier: %s", err))
        return "", ShortenOperationError{"Error shortening url", err.Error()}
    }

    hashValue := core.Base62hash(id)
    
    err = s.redisRepo.Store(hashValue, url)
    if err != nil {
        log.Error(fmt.Sprintf("Error while storing on Redis: %s", err))
        return "", ShortenOperationError{"Error shortening url", err.Error()}
    }

    return hashValue, nil
}

func (s *ShortnerService) Resolve(hashValue string) (string, error) {

    return  s.redisRepo.Get(hashValue)
}
