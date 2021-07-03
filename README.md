# gothmock

Gothmock is a CLI tool to mock APIs from API specification files. It only supports [OpenAPI 3](https://swagger.io/specification/) for now and is a WIP project in a very early stage.

## Installing:
```bash
go install github.com/renanferr/gothmock
```

## Usage:
The only required arg is the path to the spec file. This can be an OS filepath or a URI.

Flags may also be specified to determine in which port the server will listen and which response status and Content-type from the specification file should be used once a request is made to a valid path.

- Default Port: `:6666`
- Default Status Code: `200` (`"OK"`)
- Default Content-type: `"application/json"`

### Usage Examples: 
```bash
$ gothmock openapi3 ./example/openapi3/example.yml
```

```bash
$ gothmock openapi3 https://raw.githubusercontent.com/renanferr/gothmock/master/example/openapi3/example.yml --port 8080 --status 500 --content application/json
```
