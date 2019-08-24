# concourse-redis-stream-resource

[![Go Report Card](https://goreportcard.com/badge/github.com/YuukiARIA/concourse-redis-stream-resource)](https://goreportcard.com/report/github.com/YuukiARIA/concourse-redis-stream-resource)

## Source Configuration

- `host` (string): __*Required.*__ Redis host.
- `password` (string): *Optional.* Redis password.
- `key` (string): __*Required.*__ Redis stream key.
- `fields` (string array): Field names of stream entries to save.
