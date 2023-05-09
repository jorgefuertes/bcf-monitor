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

type RedisService struct {
	name    string
	host    string
	port    int
	pass    string
	timeout int
	every   int
	client  *r.Client
	ok      bool
}

func NewService(name, host string, port int, pass string, timeout, every int) *RedisService {
	s := &RedisService{name: name, host: host, port: port, pass: pass, timeout: timeout, every: every}
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.timeout)*time.Second)
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

func (s *RedisService) Every() time.Duration {
	return time.Duration(s.every) * time.Second
}