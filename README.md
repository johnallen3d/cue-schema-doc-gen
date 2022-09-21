# CUE Schema Doc Gen

Generate CUE schema documentation in HTML format.

## Development

- build docker image

```bash
docker-compose build
```

- start a container

```bash
docker-compose run --rm mod
```

- build module, execute and view output

```bash
go build && ./doc-gen ./test/schema/ ./test/dist/
```
