package store

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// Errors

type sessionStore struct {
	redis *redis.Client
	c     *gin.Context
}

func NewSessionStore(redis *redis.Client, c *gin.Context) SessionStore {
	return &sessionStore{
		redis: redis,
		c:     c,
	}
}

func (s *sessionStore) Create(username string) (Session, error) {
	newSession := Session{
		Id:       uuid.NewString(),
		Username: username,
	}

	err := s.redis.Set(s.c, newSession.Id, newSession.Username, time.Hour).Err()

	if err != nil {
		return Session{}, errors.New("Failed to create session")
	}

	return newSession, nil
}

func (s *sessionStore) GetById(sessionId string) (Session, error) {
	username, err := s.redis.Get(s.c, sessionId).Result()

	if err != nil || username == "" {
		return Session{}, errors.New("Session not found")
	}

	return Session{
		Id:       sessionId,
		Username: username,
	}, nil
}

func (s *sessionStore) Delete(sessionId string) error {
	err := s.redis.Del(s.c, sessionId).Err()
	if err != nil {
		return err
	}
	return nil
}
