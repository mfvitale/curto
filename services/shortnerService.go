package services

import (
	log "github.com/sirupsen/logrus"
	"github.com/go-redis/redis/v8"
	"github.com/mfvitale/curto/repository"
	"github.com/mfvitale/curto/services/core"
)

type shortnerService struct {
	redisClient *redis.Client
}

func NewShortnerService(redisClient *redis.Client) shortnerService {
	return shortnerService{redisClient}
}

func (s *shortnerService) Encode(url string) string {

	snowflakeGenerator := core.NewSnowflakeGenerator(int64(os.Getpid()), 10)

	id, err := snowflakeGenerator.NextID()
	if err != nil {
		return ""
	}

	hashValue, err := core.Base62hash(id)
	if err != nil {
		return ""
	}

	redisRepo := repository.NewRedisUrlRepository(s.redisClient)

	err = redisRepo.Store(hashValue, url)
	if err != nil {
		log.Error(err)
	}

	return hashValue
}
