# gothmock

Gothmock is a CLI tool to mock APIs from API specification files. It only supports [OpenAPI 3](https://swagger.io/specification/) for now and is a WIP project in a very early stage.

## Installing:
```bash
go install github.com/renanferr/gothmock
```

## Usage example:
```bash
gothmock openapi3 ./example/openapi3/example.yml --port 8080 --status 200 --content application/json
```

or

```bash
gothmock openapi3 https://raw.githubusercontent.com/renanferr/gothmock/master/example/openapi3/example.yml --port 8080 --status 200 --content application/json
```
