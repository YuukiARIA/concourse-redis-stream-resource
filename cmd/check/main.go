package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/YuukiARIA/concourse-redis-stream-resource/models"
	"github.com/gomodule/redigo/redis"
)

func performCheck(request models.CheckRequest, redisConn redis.Conn) (*models.CheckResponse, error) {
	firstID := "-"
	if request.Version != nil {
		firstID = request.Version.ID
	}

	ret, err := redisConn.Do("XRANGE", request.Source.Key, firstID, "+")
	if err != nil {
		return nil, err
	}

	var response models.CheckResponse
	for _, entry := range ret.([]interface{}) {
		entry := entry.([]interface{})
		version := models.Version{ID: string(entry[0].([]byte))}
		response = append(response, version)
	}

	return &response, nil
}

func main() {
	var request models.CheckRequest
	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		panic(err)
	}

	conn, err := redis.Dial("tcp", request.Source.Host)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		panic(err)
	}
	defer conn.Close()

	response, err := performCheck(request, conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		panic(err)
	}

	if err := json.NewEncoder(os.Stdout).Encode(&response); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		panic(err)
	}
}
