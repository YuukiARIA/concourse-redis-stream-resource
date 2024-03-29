package main

import (
	"testing"

	"github.com/YuukiARIA/concourse-redis-stream-resource/mock"
	"github.com/YuukiARIA/concourse-redis-stream-resource/models"
)

func Test_processXRangeEntry(t *testing.T) {
	expectedID := "1234567890123-0"
	expectedKey := "key-1"
	expectedValue := "value-1"

	entry := []interface{}{
		[]byte(expectedID),
		[]interface{}{
			[]byte(expectedKey),
			[]byte(expectedValue),
		},
	}

	id, fields, err := processXRangeEntry(entry)
	if err != nil {
		t.Fatal(err)
	}
	if id != expectedID {
		t.Fatalf("unexpected id: expected=%s, actual=%s\n", expectedID, id)
	}
	if len(fields) != 1 {
		t.Fatalf("unexpected map count: expected=%d, actual=%d\n", 1, len(fields))
	}
	if fields[expectedKey] != expectedValue {
		t.Fatalf("unexpected value: expected=%s, actual=%s\n", expectedValue, fields[expectedKey])
	}
}

func Test_performGet(t *testing.T) {
	redisConn := mock.NewRedisConn()
	redisConn.AddXRangeReply(
		"sample",
		"1111111111111-0",
		map[string]string{
			"name":  "test",
			"value": "123",
		},
	)

	fileRepository := mock.NewMemoryFileRepository()

	request := models.GetRequest{
		Source: models.Source{
			Host:     "localhost:6379",
			Password: nil,
			Key:      "sample",
			Fields: []string{
				"name",
				"value",
			},
		},
		Version: models.Version{
			ID: "1111111111111-0",
		},
	}

	response, err := performGet(request, fileRepository, redisConn)
	if err != nil {
		t.Fatal(err)
	}
	if response == nil {
		t.Fatal("response is nil")
	}

	if response.Version.ID != "1111111111111-0" {
		t.Fatal("ID mismatch")
	}

	if len(response.Metadata) != 2 {
		t.Fatal("metadata count mismatch")
	}

	for _, metadataEntry := range response.Metadata {
		switch metadataEntry.Name {
		case "name":
			if metadataEntry.Value != "test" {
				t.Fatal("invalid metadata (name)")
			}
		case "value":
			if metadataEntry.Value != "123" {
				t.Fatal("invalid metadata (value)")
			}
		default:
			t.Fatal("unexpected metadata: " + metadataEntry.Name)
		}
	}

	contentCount := fileRepository.Count()
	if contentCount != 2 {
		t.Fatalf("unexpected count of files to be saved: expected=2, actual=%d\n", contentCount)
	}

	contentOfName := fileRepository.GetContent("name")
	if contentOfName != "test" {
		t.Fatalf("unexpected file content (name): expected=test, actual=%s\n", contentOfName)
	}

	contentOfValue := fileRepository.GetContent("value")
	if contentOfValue != "123" {
		t.Fatalf("unexpected file content (value): expected=123, actual=%s\n", contentOfValue)
	}
}
