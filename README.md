# Kafka playground
A repository created to quickly test different kafka concepts and see how they actually work.

## Getting started
Spin up containers
`docker-compose up`

Download dependencies
`go mod download`

Run a consumer
`go run cmd/consumer/consumer.go`

Run a producer 
`go run cmd/producer/producer.go`

## Useful commands
Register a schema
```sh
curl -X POST -H "Content-Type: application/vnd.schemaregistry.v1+json" \                                                                                                              ✹
  --data '{"schema": "{\"type\": \"record\", \"namespace\": \"test\", \"name\": \"Test\", \"fields\": [{\"name\": \"id1\",\"type\": \"int\"}, {\"name\":\"name\", \"type\": \"string\"}]}"}' \
  http://localhost:8081/subjects/test/versions
```

Register a schema with specific version and id
```sh
curl -X POST -H "Content-Type: application/vnd.schemaregistry.v1+json" \                                                                                                            ⏎ ✹
  --data '{"version": 5, "id":100, "schema": "{\"type\": \"record\", \"namespace\": \"test\", \"name\": \"Test\", \"fields\": [{\"name\":\"new\",\"type\":\"string\", \"default\":\"undefined\"},{\"name\": \"id\",\"type\": \"int\"}, {\"name\":\"newname\", \"type\": \"string\"}]}"}' \
  http://localhost:8081/subjects/test/versions
```

Change import mode of a subject
```sh
curl -X PUT -H "Content-Type: application/vnd.schemaregistry.v1+json" \\n  --data '{"mode":"IMPORT"}' \\n  http://localhost:8081/mode/test\?force\=true
```
