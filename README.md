# CUE Schema Doc Gen

Generate CUE schema documentation in HTML format.

## Usage

```bash
> ./doc-gen -h
Usage of ./doc-gen:
  -i string
     input path (default ".")
  -o string
     output path (default "dist")
```

## Development

- build docker image

```bash
docker-compose build
```

- start a container

```bash
docker-compose run --rm mod
```

- convert schema to HTML

```bash
go run main.go ./test/schema/ ./test/dist/
```
