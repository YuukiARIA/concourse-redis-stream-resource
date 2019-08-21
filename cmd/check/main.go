package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/YuukiARIA/concourse-redis-stream-resource/models"
	"github.com/gomodule/redigo/redis"
)

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

	firstID := "-"
	if request.Version != nil {
		firstID = request.Version.ID
	}
	fmt.Fprintf(os.Stderr, "first id = %s\n", firstID)
	ret, err := conn.Do("XRANGE", request.Source.Key, firstID, "+")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		panic(err)
	}

	var response models.CheckResponse
	for _, entry := range ret.([]interface{}) {
		entry := entry.([]interface{})
		version := models.Version{ID: string(entry[0].([]byte))}
		response = append(response, version)
	}

	if err := json.NewEncoder(os.Stdout).Encode(&response); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		panic(err)
	}
}
