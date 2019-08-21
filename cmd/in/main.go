package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/YuukiARIA/concourse-redis-stream-resource/models"
	"github.com/gomodule/redigo/redis"
)

func writeFileString(filePath, content string) error {
	return ioutil.WriteFile(filePath, []byte(content), 0644)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "error: missing argument")
		os.Exit(1)
	}
	basePath := os.Args[1]

	var request models.GetRequest
	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		panic(err)
	}

	conn, err := redis.Dial("tcp", request.Source.Host)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		panic(err)
	}

	ret, err := conn.Do("XRANGE", request.Source.Key, request.Version.ID, request.Version.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		panic(err)
	}

	record := ret.([]interface{})[0].([]interface{})
	id := string(record[0].([]byte))

	fields, err := redis.StringMap(record[1], nil)
	if err != nil {
		panic(err)
	}

	metadata := make([]models.MetadataEntry, 0)
	for _, fieldName := range request.Source.Fields {
		value := fields[fieldName]
		if err := writeFileString(filepath.Join(basePath, fieldName), value); err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			panic(err)
		}
		metadata = append(metadata, models.MetadataEntry{Name: fieldName, Value: value})
	}

	response := models.GetResponse{
		Version:  models.Version{ID: id},
		Metadata: metadata,
	}

	if err := json.NewEncoder(os.Stdout).Encode(&response); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		panic(err)
	}
}
