package redis

import (
	"bcfmonitor/pkg/monitor/common"
	"context"
	"errors"
	"fmt"
	"time"

	r "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const TEST_KEY = "bcf:monitor:test"

type RedisService struct {
	common.MonitorizableBase
	name    string
	host    string
	port    int
	pass    string
	client  *r.Client
}

func NewService(name, host string, port int, pass string, timeout, every int) *RedisService {
	s := &RedisService{name: name, host: host, port: port, pass: pass}
	s.TimeoutSeconds = timeout
	s.EverySeconds = every
	s.client = r.NewClient(&r.Options{
		Addr:                  s.Address(),
		Password:              s.pass,
		DB:                    0,
		ContextTimeoutEnabled: true,
	})
	return s
}

func (s *RedisService) Address() string {
	return fmt.Sprintf("%s:%d", s.host, s.port)
}

func (s *RedisService) Check() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.Timeout())
	defer cancel()
	value := primitive.NewObjectID().Hex()
	if err := s.client.Set(ctx, TEST_KEY, value, time.Duration(2*time.Second)).Err(); err != nil {
		s.AddFail()
		return err
	}
	getValue, err := s.client.Get(ctx, TEST_KEY).Result()
	if err != nil {
		s.AddFail()
		return err
	}
	if getValue != value {
		s.AddFail()
		return errors.New("read value not equal to set value")
	}

	s.Reset()
	return nil
}

func (s *RedisService) Type() string {
	return "cache"
}

func (s *RedisService) Name() string {
	return s.name
}
