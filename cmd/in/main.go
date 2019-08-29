package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/YuukiARIA/concourse-redis-stream-resource/models"
	"github.com/YuukiARIA/concourse-redis-stream-resource/resource"
	"github.com/gomodule/redigo/redis"
)

func processXRangeEntry(entry interface{}) (string, map[string]string, error) {
	pair, ok := entry.([]interface{})
	if !ok {
		return "", nil, errors.New("type mismatch (pair)")
	}

	idBytes, ok := pair[0].([]byte)
	if !ok {
		return "", nil, errors.New("type mismatch (id)")
	}
	id := string(idBytes)

	fields, err := redis.StringMap(pair[1], nil)
	if err != nil {
		return "", nil, err
	}

	return id, fields, nil
}

func performGet(request models.GetRequest, fileRepository resource.FileRepository, redisConn redis.Conn) (*models.GetResponse, error) {
	entries, err := redis.Values(redisConn.Do("XRANGE", request.Source.Key, request.Version.ID, request.Version.ID))
	if err != nil {
		return nil, err
	}

	if len(entries) < 1 {
		return nil, errors.New("no entry")
	}

	id, fields, err := processXRangeEntry(entries[0])
	if err != nil {
		return nil, err
	}

	metadata := make([]models.MetadataEntry, 0)
	for _, fieldName := range request.Source.Fields {
		value := fields[fieldName]
		if err := fileRepository.Write(fieldName, value); err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			return nil, err
		}
		metadata = append(metadata, models.MetadataEntry{Name: fieldName, Value: value})
	}

	return &models.GetResponse{
		Version:  models.Version{ID: id},
		Metadata: metadata,
	}, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "error: missing argument")
		os.Exit(1)
	}
	fileRepository := resource.NewFileSystemRepository(os.Args[1])

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
	defer conn.Close()

	response, err := performGet(request, fileRepository, conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		panic(err)
	}

	if err := json.NewEncoder(os.Stdout).Encode(&response); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		panic(err)
	}
}
