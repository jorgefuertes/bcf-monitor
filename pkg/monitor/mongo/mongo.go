package mongo

import (
	"bcfmonitor/pkg/monitor/common"
	"context"
	"fmt"
	"time"

	mongodb "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoService struct {
	common.MonitorizableBase
	name   string
	host   string
	port   int
	client *mongodb.Client
}

func NewService(name, host string, port int, SSL bool, timeout int, every int) *MongoService {
	s := &MongoService{name: name, host: host, port: port}
	s.TimeoutSeconds = timeout
	s.EverySeconds = every

	return s
}

func (s *MongoService) connect(ctx context.Context) error {
	var err error
	s.client, err = mongodb.Connect(ctx, options.Client().ApplyURI(s.Address()))
	if err != nil {
		return err
	}

	return nil
}

func (s *MongoService) disconnect(ctx context.Context) error {
	return s.client.Disconnect(ctx)
}


func (s *MongoService) Address() string {
	return fmt.Sprintf("mongodb://%s:%d", s.host, s.port)
}

func (s *MongoService) Check() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.Timeout())
	defer cancel()
	if err := s.connect(ctx); err != nil {
		s.AddFail()
		return err
	}
	defer s.disconnect(ctx)
	if err := s.client.Ping(ctx, readpref.Primary()); err != nil {
		s.AddFail()
		return err
	}

	status, err := s.client.Database("admin").RunCommand(ctx, map[string]interface{}{"serverStatus": 1}).DecodeBytes()
	if err != nil {
		s.AddFail()
		return fmt.Errorf("reading uptime: %s", err)
	}
	if time.Since(time.Unix(status.Lookup("uptime").AsInt64(), 0)) < s.Every() {
		s.AddFail()
		return fmt.Errorf("uptime is lower than %d seconds", s.EverySeconds)
	}

	s.Reset()
	return nil
}

func (s *MongoService) Type() string {
	return "database"
}

func (s *MongoService) Name() string {
	return s.name
}
