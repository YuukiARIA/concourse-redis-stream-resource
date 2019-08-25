package mock

import (
	"errors"
	"strings"

	"github.com/gomodule/redigo/redis"
)

type redisConn struct {
	redis.Conn
	xrangeReply map[string][]interface{}
}

func NewRedisConn() redisConn {
	return redisConn{xrangeReply: map[string][]interface{}{}}
}

func (r redisConn) SetXRangeReplyEmpty(key string) {
	r.xrangeReply[key] = make([]interface{}, 0)
}

func (r redisConn) AddXRangeReply(key, id string, fieldMap map[string]string) {
	fields := make([]interface{}, 0)
	for k, v := range fieldMap {
		fields = append(fields, []byte(k), []byte(v))
	}

	entry := []interface{}{[]byte(id), fields}

	entries := r.xrangeReply[key]
	if entries == nil {
		entries = make([]interface{}, 0)
	}
	entries = append(entries, entry)
	r.xrangeReply[key] = entries
}

func (r redisConn) Do(command string, args ...interface{}) (interface{}, error) {
	switch strings.ToUpper(command) {
	case "XRANGE":
		key := args[0].(string)
		return r.xrangeReply[key], nil
	default:
		return nil, errors.New("unsupported")
	}
}
