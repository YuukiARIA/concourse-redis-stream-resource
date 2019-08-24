package main

import (
	"errors"
	"testing"

	"github.com/YuukiARIA/concourse-redis-stream-resource/models"
	"github.com/gomodule/redigo/redis"
)

type mockRedisConn struct {
	redis.Conn
	xrangeResponse map[string]interface{}
}

func (m mockRedisConn) Do(command string, args ...interface{}) (interface{}, error) {
	if command == "XRANGE" {
		key := args[0].(string)
		return m.xrangeResponse[key], nil
	} else {
		return nil, errors.New("unsupported")
	}
}

var conn = mockRedisConn{
	xrangeResponse: map[string]interface{}{
		"sample": []interface{}{
			[]interface{}{
				[]byte("1111111111111-0"),
				[][]byte{
					[]byte("name"),
					[]byte("test"),
					[]byte("value"),
					[]byte("123"),
				},
			},
		},
		"empty": []interface{}{},
	},
}

func Test_performCheck(t *testing.T) {
	request := models.CheckRequest{
		Source: models.Source{
			Host:     "localhost:6379",
			Password: nil,
			Key:      "sample",
			Fields:   []string{},
		},
		Version: nil,
	}

	response, err := performCheck(request, conn)
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
	request := models.CheckRequest{
		Source: models.Source{
			Host:     "localhost:6379",
			Password: nil,
			Key:      "empty",
			Fields:   []string{},
		},
		Version: nil,
	}

	response, err := performCheck(request, conn)
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
