package mongo

import (
	"context"
	"fmt"
	"time"

	mongodb "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoService struct {
	name    string
	host    string
	port    int
	timeout int
	every   int
	client  *mongodb.Client
	ok      bool
}

func NewService(name, host string, port int, SSL bool, timeout int, every int) *MongoService {
	return &MongoService{name: name, host: host, port: port, timeout: timeout, every: every}
}

func (s *MongoService) getURI() string {
	return fmt.Sprintf("mongodb://%s:%d", s.host, s.port)
}

func (s *MongoService) connect(ctx context.Context) error {
	var err error
	s.client, err = mongodb.NewClient(options.Client().ApplyURI(s.getURI()))
	if err != nil {
		return err
	}
	if err := s.client.Connect(ctx); err != nil {
		return err
	}

	return nil
}

func (s *MongoService) disconnect(ctx context.Context) error {
	return s.client.Disconnect(ctx)
}

func (s *MongoService) Check() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.timeout)*time.Second)
	defer cancel()
	if err := s.connect(ctx); err != nil {
		return err
	}
	defer s.disconnect(ctx)
	return s.client.Ping(ctx, readpref.Primary())
}

func (s *MongoService) IsUp() bool {
	return s.ok
}

func (s *MongoService) Down() {
	s.ok = false
}

func (s *MongoService) Up() {
	s.ok = true
}

func (s *MongoService) Every() time.Duration {
	return time.Duration(s.every) * time.Second
}