package main

import (
	"testing"

	"github.com/YuukiARIA/concourse-redis-stream-resource/mock"
	"github.com/YuukiARIA/concourse-redis-stream-resource/models"
)

func Test_performCheck(t *testing.T) {
	redisConn := mock.NewRedisConn()
	redisConn.AddXRangeReply(
		"sample",
		"1111111111111-0",
		map[string]string{
			"name":  "test",
			"value": "123",
		},
	)

	request := models.CheckRequest{
		Source: models.Source{
			Host:     "localhost:6379",
			Password: nil,
			Key:      "sample",
			Fields:   []string{},
		},
		Version: nil,
	}

	response, err := performCheck(request, redisConn)
	if err != nil {
		t.Fatal(err)
	}
	if response == nil {
		t.Fatal("response is nil")
	}
	if len(*response) != 1 {
		t.Fatal("length")
	}
	if (*response)[0].ID != "1111111111111-0" {
		t.Fatal("ID mismatch")
	}
}

func Test_performCheck_NoEntry(t *testing.T) {
	redisConn := mock.NewRedisConn()
	redisConn.AddXRangeReply(
		"sample",
		"1111111111111-0",
		map[string]string{
			"name":  "test",
			"value": "123",
		},
	)

	request := models.CheckRequest{
		Source: models.Source{
			Host:     "localhost:6379",
			Password: nil,
			Key:      "empty",
			Fields:   []string{},
		},
		Version: nil,
	}

	response, err := performCheck(request, redisConn)
	if err != nil {
		t.Fatal(err)
	}
	if response == nil {
		t.Fatal("response is nil")
	}
	if len(*response) != 0 {
		t.Fatal("length")
	}
}
