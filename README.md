# gothmock

## Installing:
```bash
go install github.com/renanferr/gothmock
```

## Usage example:
```bash
gothmock openapi3 ./example/openapi3/example.yml --port 8080 --status 200 --content application/json
```

of

```bash
gothmock openapi3 https://raw.githubusercontent.com/renanferr/gothmock/master/example/openapi3/example.yml --port 8080 --status 200 --content application/json
```
