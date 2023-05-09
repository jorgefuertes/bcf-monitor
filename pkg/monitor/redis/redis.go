package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	r "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const TEST_KEY = "bcf:monitor:test"
const TIMEOUT_SEC = 5

type RedisService struct {
	name   string
	host   string
	port   int
	pass   string
	client *r.Client
	ok     bool
}

func NewService(name, host string, port int, pass string) *RedisService {
	s := &RedisService{name: name, host: host, port: port, pass: pass}
	s.client = r.NewClient(&r.Options{
		Addr:                  s.getAddress(),
		Password:              s.pass,
		DB:                    0,
		ContextTimeoutEnabled: true,
	})
	return s
}

func (s *RedisService) getAddress() string {
	return fmt.Sprintf("%s:%d", s.host, s.port)
}

func (s *RedisService) Check() error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_SEC*time.Second)
	defer cancel()
	value := primitive.NewObjectID().Hex()
	if err := s.client.Set(ctx, TEST_KEY, value, time.Duration(2*time.Second)).Err(); err != nil {
		return err
	}
	getValue, err := s.client.Get(ctx, TEST_KEY).Result()
	if err != nil {
		return err
	}
	if getValue != value {
		return errors.New("read value not equal to set value")
	}
	return nil
}

func (s *RedisService) IsUp() bool {
	return s.ok
}

func (s *RedisService) Down() {
	s.ok = false
}

func (s *RedisService) Up() {
	s.ok = true
}
