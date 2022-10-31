package services

import (
    log "github.com/sirupsen/logrus"
    "github.com/mfvitale/curto/repository"
    "github.com/mfvitale/curto/services/core"
)

type ShortnerService struct {
    redisRepo repository.UrlRepository
	identifierService core.IdentifierService
}

func NewShortnerService(urlRepository repository.UrlRepository,
	identifierService core.IdentifierService) ShortnerService {
    return ShortnerService{urlRepository, identifierService}
}

func (s *ShortnerService) Shorten(url string) string {

    id, err := s.identifierService.NextID()
    if err != nil {
        return ""
    }

    hashValue, err := core.Base62hash(id)
    if err != nil {
        return ""
    }

    err = s.redisRepo.Store(hashValue, url)
    if err != nil {
        log.Error(err)
    }

    return hashValue
}

func (s *ShortnerService) Resolve(hashValue string) string {

    originalUrl, err := s.redisRepo.Get(hashValue)
    if err != nil {
        log.Error(err)
    }

    return originalUrl
}
